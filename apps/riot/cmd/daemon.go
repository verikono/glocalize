package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

const (
	daemonArgv     = "serve --daemon"
	executableName = "riot"
	logFileName    = "riot.log"
)

// StartRiot launches the server as a detached background process and returns
// immediately. It is a no-op when a server is already listening on the port.
func StartRiot(port int) error {
	return startRiot(port, "")
}

// StartRiotHot launches the server detached with hot-reloading armed on the
// resolved riot dev directory (RIOT_DEV_DIR, else the current directory).
func StartRiotHot(port int) error {
	devDir, err := utils.ResolveHotDir()
	if err != nil {
		return err
	}
	return startRiot(port, devDir)
}

func startRiot(port int, hotDir string) error {
	if isRiotRunning(port) {
		fmt.Println("riot: already running")
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("riot: resolve executable: %w", err)
	}

	logPath := filepath.Join(os.TempDir(), logFileName)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("riot: open log: %w", err)
	}
	defer logFile.Close()

	child := exec.Command(exe, "serve", "--daemon", "--port", strconv.Itoa(port))
	if hotDir != "" {
		child.Env = withEnv(os.Environ(), utils.RiotDevDirEnv, hotDir)
	}
	child.Stdout = logFile
	child.Stderr = logFile
	child.Stdin = nil
	child.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := child.Start(); err != nil {
		return fmt.Errorf("riot: start: %w", err)
	}
	pid := child.Process.Pid
	if err := child.Process.Release(); err != nil {
		return fmt.Errorf("riot: release: %w", err)
	}

	hot := willRunHot(hotDir)

	if err := saveRuntimeFlags(runtimeFlags{Port: port, Hot: hot, DevDir: hotDir}); err != nil {
		return err
	}

	fmt.Print(startBanner(pid, port, logPath, hot))
	return nil
}

// willRunHot reports whether the daemon child will arm hot reload: either hot
// was requested directly or RIOT_DEV_DIR is inherited from the environment.
func willRunHot(hotDir string) bool {
	return hotDir != "" || os.Getenv(utils.RiotDevDirEnv) != ""
}

// withEnv returns env with key replaced by key=value, so the daemon child
// inherits a single, authoritative value even when key is already present.
func withEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			continue
		}
		out = append(out, e)
	}
	return append(out, prefix+value)
}

// startBanner formats the success banner printed after a new daemon starts.
func startBanner(pid, port int, logPath string, hot bool) string {
	base := fmt.Sprintf("http://localhost:%d", port)

	lines := [][2]string{
		{"pid", strconv.Itoa(pid)},
		{"gql", base + "/graph"},
		{"mcp", base + "/mcp"},
		{"database", utils.ResolveDataURL()},
	}
	if hot {
		lines = append(lines, [2]string{"hot", "yes (hot-reload armed)"})
	}
	lines = append(lines, [2]string{"logs", logPath})

	var b strings.Builder
	b.WriteString("riot started\n")
	for _, line := range lines {
		fmt.Fprintf(&b, "%-12s : %s\n", line[0], line[1])
	}
	return b.String()
}

// isRiotRunning reports whether anything is listening on 127.0.0.1:port.
func isRiotRunning(port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), 300*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// waitPortFree waits until nothing is listening on port, or the timeout elapses.
func waitPortFree(port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !isRiotRunning(port) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return !isRiotRunning(port)
}

// RunningPid returns the pid of the running riot daemon, or an error.
func RunningPid() (int, error) {
	return findRiotPid()
}

// findRiotPid scans /proc for a process whose argv contains the daemon marker
// and whose executable basename is "riot".
func findRiotPid() (int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if !matchesDaemonArgv(pid) {
			continue
		}
		if !isRiotExecutable(pid) {
			continue
		}
		return pid, nil
	}
	return 0, fmt.Errorf("riot: no running daemon found")
}

// findRiotPids returns every pid whose executable basename is "riot", except
// the calling process itself (so `riot down` never terminates itself).
func findRiotPids() ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	self := os.Getpid()
	var pids []int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if pid == self {
			continue
		}
		if isRiotExecutable(pid) {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

func matchesDaemonArgv(pid int) bool {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	if err != nil {
		return false
	}
	return strings.Contains(strings.ReplaceAll(string(data), "\x00", " "), daemonArgv)
}

func isRiotExecutable(pid int) bool {
	exe, err := os.Readlink(filepath.Join("/proc", strconv.Itoa(pid), "exe"))
	if err != nil {
		return false
	}
	exe = strings.TrimSuffix(exe, " (deleted)")
	return filepath.Base(exe) == executableName
}

func isProcessAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}

// procUptime returns how long pid has been running, read from /proc.
func procUptime(pid int) (time.Duration, error) {
	stat, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return 0, err
	}
	s := string(stat)
	// comm is in parens and may contain spaces — split after the last ')'.
	end := strings.LastIndexByte(s, ')')
	if end < 0 {
		return 0, fmt.Errorf("riot: malformed stat for pid %d", pid)
	}
	fields := strings.Fields(s[end+1:])
	// fields[0] is state (field 3); starttime is field 22 → index 19 here.
	if len(fields) < 20 {
		return 0, fmt.Errorf("riot: short stat for pid %d", pid)
	}
	startTicks, err := strconv.ParseUint(fields[19], 10, 64)
	if err != nil {
		return 0, err
	}
	const clkTck = 100
	start := time.Duration(startTicks) * time.Second / clkTck

	uptimeRaw, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	upSeconds, _, _ := strings.Cut(strings.TrimSpace(string(uptimeRaw)), " ")
	up, err := strconv.ParseFloat(upSeconds, 64)
	if err != nil {
		return 0, err
	}
	systemUp := time.Duration(up * float64(time.Second))
	return systemUp - start, nil
}
