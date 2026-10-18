package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/verikono/glocalize/packages/go/puck"
)

const (
	riotDevDirEnv  = "RIOT_DEV_DIR"
	riotModuleDecl = "module github.com/verikono/glocalize/apps/riot"
	hotTimestamp   = "06/01/02 15:04:05"
	hotDebounce    = 300 * time.Millisecond

	// hotDispatchFlush bounds how long a refresh dispatch may take to reach
	// connected clients before the process is replaced by exec.
	hotDispatchFlush = 250 * time.Millisecond
)

// RiotDevDirEnv names the environment variable that points at the riot module
// directory for the hot service.
const RiotDevDirEnv = riotDevDirEnv

// hotLog writes hot-service lines without the standard log prefix: change
// lines carry their own [yy/mm/dd hh:mm:ss] timestamp per the spec.
var hotLog = log.New(os.Stdout, "", 0)

// ValidateDevDir verifies that RIOT_DEV_DIR, when set, points at the riot
// module. It is a no-op returning nil when RIOT_DEV_DIR is unset.
func ValidateDevDir() error {
	_, err := resolveDevDir()
	return err
}

// SetServiceHot arms the development hot-reloader. When RIOT_DEV_DIR points at
// the riot module, any change beneath it triggers a rebuild and restart of the
// running daemon. It is a no-op when RIOT_DEV_DIR is unset.
//
// dispatch is called with a refresh message before the daemon is replaced, so
// connected browser clients reload and reconnect to the new process.
func SetServiceHot(port int, dispatch func(string)) error {
	devDir, err := resolveDevDir()
	if err != nil {
		return err
	}
	if devDir == "" {
		return nil
	}

	reloader := newHotReloader(port, devDir, dispatch)
	p, ok := puck.From(devDir).(*puck.Puck)
	if !ok {
		return fmt.Errorf("riot: hot service: unexpected puck type %T", puck.From(devDir))
	}
	if err := p.OnChange(puck.WithWatcher(reloader.onChange)); err != nil {
		return fmt.Errorf("riot: hot service: %w", err)
	}
	go reloader.run()
	return nil
}

// resolveDevDir returns the absolute, validated RIOT_DEV_DIR or "" when unset.
func resolveDevDir() (string, error) {
	devDir := strings.TrimSpace(os.Getenv(riotDevDirEnv))
	if devDir == "" {
		return "", nil
	}
	return resolveRiotDir(devDir)
}

// ResolveHotDir returns the absolute riot-module directory the hot service
// should watch: RIOT_DEV_DIR when set, otherwise the current working directory.
func ResolveHotDir() (string, error) {
	devDir := strings.TrimSpace(os.Getenv(riotDevDirEnv))
	if devDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("riot: hot service: resolve working directory: %w", err)
		}
		return resolveRiotDir(cwd)
	}
	return resolveRiotDir(devDir)
}

func resolveRiotDir(devDir string) (string, error) {
	abs, err := filepath.Abs(devDir)
	if err != nil {
		return "", fmt.Errorf("riot: hot service: resolve %s: %w", riotDevDirEnv, err)
	}
	if err := assertRiotModule(abs); err != nil {
		return "", err
	}
	return abs, nil
}

func assertRiotModule(devDir string) error {
	data, err := os.ReadFile(filepath.Join(devDir, "go.mod"))
	if err != nil {
		return fmt.Errorf("riot: hot service: %s is not a Go module: %w", devDir, err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == riotModuleDecl {
			return nil
		}
	}
	return fmt.Errorf("riot: hot service: %s is the wrong directory for Glocalize::Riot (go.mod does not declare %q)", devDir, riotModuleDecl)
}

// hotReloader debounces file changes and rebuilds + restarts the daemon once
// the tree has settled.
type hotReloader struct {
	port     int
	devDir   string
	events   chan puck.FSChangeEvent
	dispatch func(string)
}

func newHotReloader(port int, devDir string, dispatch func(string)) *hotReloader {
	return &hotReloader{
		port:     port,
		devDir:   devDir,
		events:   make(chan puck.FSChangeEvent, 64),
		dispatch: dispatch,
	}
}

// onChange is the puck watcher callback. It logs the change and feeds the
// debounced rebuild loop; changes under .git/ or dist/ are ignored so a rebuild
// never re-triggers itself.
func (h *hotReloader) onChange(event puck.FSChangeEvent) error {
	if ignoredPath(h.devDir, event.Path) {
		return nil
	}
	hotLog.Printf("[%s] [%s]:%s", time.Now().Format(hotTimestamp), changeType(event.Op), event.Path)
	select {
	case h.events <- event:
	default:
	}
	return nil
}

// run debounces incoming changes and rebuilds once the burst settles.
func (h *hotReloader) run() {
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	for {
		select {
		case <-h.events:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(hotDebounce)
		case <-timer.C:
			h.rebuildAndRestart()
		}
	}
}

func (h *hotReloader) rebuildAndRestart() {
	if err := h.rebuild(); err != nil {
		hotLog.Printf("[%s] [ERROR] %v", time.Now().Format(hotTimestamp), err)
		return
	}
	// Tell connected browser clients to reload before the process is replaced.
	// Dispatch is queued, so give the writer a bounded moment to flush: the
	// exec below destroys this process and anything still in flight with it.
	if h.dispatch != nil {
		h.dispatch("riot::refresh")
		time.Sleep(hotDispatchFlush)
	}
	bin := filepath.Join(h.devDir, "dist", "riot")
	argv := []string{bin, "serve", "--daemon", "--port", strconv.Itoa(h.port)}
	if err := syscall.Exec(bin, argv, os.Environ()); err != nil {
		hotLog.Printf("[%s] [ERROR] %v", time.Now().Format(hotTimestamp), err)
	}
}

func (h *hotReloader) rebuild() error {
	cmd := exec.Command("go", "build", "-o", filepath.Join("dist", "riot"), ".")
	cmd.Dir = h.devDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rebuild failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// changeType maps an fsnotify op to the primary change token logged.
func changeType(op fsnotify.Op) string {
	switch {
	case op&fsnotify.Write != 0:
		return "WRITE"
	case op&fsnotify.Create != 0:
		return "CREATE"
	case op&fsnotify.Remove != 0:
		return "REMOVE"
	case op&fsnotify.Rename != 0:
		return "RENAME"
	case op&fsnotify.Chmod != 0:
		return "CHMOD"
	default:
		return "CHANGE"
	}
}

// ignoredPath reports whether a change under devDir should be ignored. Build
// output (dist/) and the git store are excluded.
func ignoredPath(devDir, path string) bool {
	rel, err := filepath.Rel(devDir, path)
	if err != nil {
		return true
	}
	if rel == "." {
		return true
	}
	top := strings.SplitN(rel, string(filepath.Separator), 2)[0]
	return top == ".git" || top == "dist"
}
