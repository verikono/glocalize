package puck

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	NTFY "github.com/fsnotify/fsnotify"
)

// FSChangeEvent describes a single change emitted by a watched resource.
type FSChangeEvent struct {
	Op   NTFY.Op
	Path AbsoluteFsPath
}

// FileSystemWatcher is the callback invoked for every file-system change.
type FileSystemWatcher = func(FSChangeEvent) error

// Watcher is the option target populated by WithWatcher.
type Watcher struct {
	fn FileSystemWatcher
}

// watchTarget owns a live fsnotify watcher and routes events to one listener.
type watchTarget struct {
	watcher *NTFY.Watcher
	fn      FileSystemWatcher
	once    sync.Once
}

// run pumps fsnotify events into the listener until the watcher is closed.
// Newly created directories are added to the watcher so the whole tree stays watched.
func (t *watchTarget) run() {
	for {
		select {
		case event, ok := <-t.watcher.Events:
			if !ok {
				return
			}
			if event.Op&NTFY.Create != 0 {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					_ = addRecursive(t.watcher, AbsoluteFsPath(event.Name))
				}
			}
			_ = t.fn(FSChangeEvent{Op: event.Op, Path: AbsoluteFsPath(event.Name)})
		case err, ok := <-t.watcher.Errors:
			if !ok {
				return
			}
			_ = err
		}
	}
}

// close shuts the underlying fsnotify watcher down exactly once.
func (t *watchTarget) close() {
	t.once.Do(func() {
		_ = t.watcher.Close()
	})
}

// watchPaths resolves the absolute paths this puck should watch.
func (p *Puck) watchPaths() ([]AbsoluteFsPath, error) {
	if p.glob == "" {
		return []AbsoluteFsPath{p.Path()}, nil
	}

	matches, err := filepath.Glob(p.glob)
	if err != nil {
		return nil, fmt.Errorf("invalid glob %s: %w", p.glob, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("glob %s matched no files", p.glob)
	}
	return matches, nil
}

// addWatchTarget creates a live watcher for paths and dispatches changes to fn.
func (p *Puck) addWatchTarget(fn FileSystemWatcher, paths []AbsoluteFsPath) error {
	fsWatcher, err := NTFY.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed creating watcher: %w", err)
	}

	for _, path := range paths {
		if err := addRecursive(fsWatcher, path); err != nil {
			_ = fsWatcher.Close()
			return fmt.Errorf("failed watching %s: %w", path, err)
		}
	}

	target := &watchTarget{watcher: fsWatcher, fn: fn}

	p.watchersMu.Lock()
	p.watchers = append(p.watchers, target)
	p.watchersMu.Unlock()

	go target.run()
	return nil
}

// addRecursive adds path to the watcher. Directories are added together with
// every subdirectory beneath them, so changes anywhere in the tree are delivered.
func addRecursive(w *NTFY.Watcher, path AbsoluteFsPath) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if err := w.Add(path); err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}
	return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && p != path {
			return w.Add(p)
		}
		return nil
	})
}

// stopWatcher closes every watchTarget whose listener is fn.
func (p *Puck) stopWatcher(fn FileSystemWatcher) {
	targetFn := reflect.ValueOf(fn).Pointer()

	p.watchersMu.Lock()
	kept := p.watchers[:0]
	var closing []*watchTarget
	for _, target := range p.watchers {
		if reflect.ValueOf(target.fn).Pointer() == targetFn {
			closing = append(closing, target)
		} else {
			kept = append(kept, target)
		}
	}
	p.watchers = kept
	p.watchersMu.Unlock()

	for _, target := range closing {
		target.close()
	}
}

// stopAllWatchers closes every watchTarget held by this puck.
func (p *Puck) stopAllWatchers() {
	p.watchersMu.Lock()
	closing := p.watchers
	p.watchers = nil
	p.watchersMu.Unlock()

	for _, target := range closing {
		target.close()
	}
}
