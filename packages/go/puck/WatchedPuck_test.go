package puck_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	NTFY "github.com/fsnotify/fsnotify"

	"github.com/verikono/glocalize/packages/go/puck"
)

// --- helpers --- //

func fromPuck(t *testing.T, pattern string) *puck.Puck {
	t.Helper()
	p, ok := puck.From(pattern).(*puck.Puck)
	if !ok {
		t.Fatalf("From(%q) returned %T, want *puck.Puck", pattern, p)
	}
	return p
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}

func awaitEvent(t *testing.T, ch <-chan puck.FSChangeEvent) puck.FSChangeEvent {
	t.Helper()
	select {
	case event := <-ch:
		return event
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for fs event")
		return puck.FSChangeEvent{}
	}
}

func assertNoEvent(t *testing.T, ch <-chan puck.FSChangeEvent, d time.Duration) {
	t.Helper()
	select {
	case event := <-ch:
		t.Fatalf("unexpected event: %+v", event)
	case <-time.After(d):
	}
}

// --- watcher behaviour --- //

func TestOnChange_SingleFile_DeliversEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watch.prisma")
	writeFile(t, path, "v1")

	p := fromPuck(t, path)
	events := make(chan puck.FSChangeEvent, 16)
	if err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error {
		events <- event
		return nil
	})); err != nil {
		t.Fatalf("OnChange: %v", err)
	}
	t.Cleanup(func() { _ = p.StopWatching() })

	writeFile(t, path, "v2")

	event := awaitEvent(t, events)
	if event.Path != path {
		t.Fatalf("event.Path = %q, want %q", event.Path, path)
	}
	if event.Op&NTFY.Write == 0 {
		t.Fatalf("event.Op = %v, want a write op", event.Op)
	}
}

func TestOnChange_Glob_DeliversEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.prisma")
	writeFile(t, path, "v1")
	glob := filepath.Join(dir, "*.prisma")

	p := fromPuck(t, glob)
	events := make(chan puck.FSChangeEvent, 16)
	if err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error {
		events <- event
		return nil
	})); err != nil {
		t.Fatalf("OnChange: %v", err)
	}
	t.Cleanup(func() { _ = p.StopWatching() })

	writeFile(t, path, "v2")

	event := awaitEvent(t, events)
	if event.Path != path {
		t.Fatalf("event.Path = %q, want %q", event.Path, path)
	}
}

func TestOnChange_RequiresWatcher(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "x.prisma"))
	err := p.OnChange()
	if err == nil || !strings.Contains(err.Error(), "no function acquired") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnChange_PropagatesOptionError(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "x.prisma"))
	boom := errors.New("boom")
	if err := p.OnChange(func(args ...any) error { return boom }); !errors.Is(err, boom) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnChange_NonExistentPathFails(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "nope.prisma"))
	err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error { return nil }))
	if err == nil || !strings.Contains(err.Error(), "failed watching") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnChange_GlobNoMatchesFails(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "*.prisma"))
	err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error { return nil }))
	if err == nil || !strings.Contains(err.Error(), "matched no files") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnChange_InvalidGlobFails(t *testing.T) {
	p := fromPuck(t, "/tmp/[bad*.prisma")
	err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error { return nil }))
	if err == nil || !strings.Contains(err.Error(), "invalid glob") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStopWatching_StopsDelivery(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watch.prisma")
	writeFile(t, path, "v1")

	p := fromPuck(t, path)
	events := make(chan puck.FSChangeEvent, 16)
	fn := func(event puck.FSChangeEvent) error {
		events <- event
		return nil
	}
	if err := p.OnChange(puck.WithWatcher(fn)); err != nil {
		t.Fatalf("OnChange: %v", err)
	}
	if err := p.StopWatching(puck.WithWatcher(fn)); err != nil {
		t.Fatalf("StopWatching: %v", err)
	}

	writeFile(t, path, "v2")
	assertNoEvent(t, events, 300*time.Millisecond)
}

func TestStopWatching_StopsAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watch.prisma")
	writeFile(t, path, "v1")

	p := fromPuck(t, path)
	events := make(chan puck.FSChangeEvent, 16)
	if err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error {
		events <- event
		return nil
	})); err != nil {
		t.Fatalf("OnChange: %v", err)
	}
	if err := p.OnChange(puck.WithWatcher(func(event puck.FSChangeEvent) error {
		return nil
	})); err != nil {
		t.Fatalf("OnChange: %v", err)
	}

	if err := p.StopWatching(); err != nil {
		t.Fatalf("StopWatching: %v", err)
	}

	writeFile(t, path, "v2")
	assertNoEvent(t, events, 300*time.Millisecond)
}

func TestStopWatching_OnlyStopsMatchingListener(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watch.prisma")
	writeFile(t, path, "v1")

	p := fromPuck(t, path)
	eventsA := make(chan puck.FSChangeEvent, 16)
	eventsB := make(chan puck.FSChangeEvent, 16)
	fnA := func(event puck.FSChangeEvent) error {
		eventsA <- event
		return nil
	}
	fnB := func(event puck.FSChangeEvent) error {
		eventsB <- event
		return nil
	}
	if err := p.OnChange(puck.WithWatcher(fnA)); err != nil {
		t.Fatalf("OnChange A: %v", err)
	}
	if err := p.OnChange(puck.WithWatcher(fnB)); err != nil {
		t.Fatalf("OnChange B: %v", err)
	}

	if err := p.StopWatching(puck.WithWatcher(fnA)); err != nil {
		t.Fatalf("StopWatching: %v", err)
	}

	writeFile(t, path, "v2")

	awaitEvent(t, eventsB)
	assertNoEvent(t, eventsA, 300*time.Millisecond)
}

func TestStopWatching_UnknownListenerIsNoOp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watch.prisma")
	writeFile(t, path, "v1")

	p := fromPuck(t, path)
	events := make(chan puck.FSChangeEvent, 16)
	fn := func(event puck.FSChangeEvent) error {
		events <- event
		return nil
	}
	if err := p.OnChange(puck.WithWatcher(fn)); err != nil {
		t.Fatalf("OnChange: %v", err)
	}
	t.Cleanup(func() { _ = p.StopWatching() })

	if err := p.StopWatching(puck.WithWatcher(func(event puck.FSChangeEvent) error {
		return nil
	})); err != nil {
		t.Fatalf("StopWatching: %v", err)
	}

	writeFile(t, path, "v2")
	awaitEvent(t, events)
}

func TestStopWatching_PropagatesOptionError(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "x.prisma"))
	boom := errors.New("boom")
	if err := p.StopWatching(func(args ...any) error { return boom }); !errors.Is(err, boom) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithWatcher_IgnoresNonWatcherArgs(t *testing.T) {
	opt := puck.WithWatcher(func(event puck.FSChangeEvent) error { return nil })
	if err := opt(1, "two", 3.0); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

// --- puck reading / metadata --- //

func TestPuck_FileReadsAndMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.prisma")
	writeFile(t, path, "datasource db {}")

	p := fromPuck(t, path)

	if got := p.Path(); got != path {
		t.Fatalf("Path = %q, want %q", got, path)
	}
	if got := p.Url(); got != "file://"+path {
		t.Fatalf("Url = %q", got)
	}
	if got := p.Protocol(); got != puck.ProtocolFile {
		t.Fatalf("Protocol = %q", got)
	}
	if got := p.Format(); got != puck.FormatPrisma {
		t.Fatalf("Format = %q", got)
	}
	if got := p.Extension(); got != "prisma" {
		t.Fatalf("Extension = %q", got)
	}
	if got := p.ResourceType(); got != puck.FileResource {
		t.Fatalf("ResourceType = %q", got)
	}
	if !p.Exists() {
		t.Fatal("Exists = false, want true")
	}
	if err := p.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	content, err := p.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(content) != "datasource db {}" {
		t.Fatalf("Read = %q", content)
	}
	if size := p.Size(); size != len("datasource db {}") {
		t.Fatalf("Size = %d", size)
	}
}

func TestPuck_Directory(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.prisma"), "aaa")
	writeFile(t, filepath.Join(dir, "b.prisma"), "bbbbb")

	p := fromPuck(t, dir)

	if got := p.ResourceType(); got != puck.DirectoryResource {
		t.Fatalf("ResourceType = %q", got)
	}
	if got := p.Extension(); got != "" {
		t.Fatalf("Extension = %q, want empty", got)
	}
	if _, err := p.Read(); err == nil {
		t.Fatal("Read of directory should error")
	}
	if size := p.Size(); size != 8 {
		t.Fatalf("Size = %d, want 8", size)
	}
}

func TestPuck_GlobReadsAndSizes(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.prisma"), "aaa")
	writeFile(t, filepath.Join(dir, "b.prisma"), "bbbbb")
	glob := filepath.Join(dir, "*.prisma")

	p := fromPuck(t, glob)

	if got := p.Path(); got != "" {
		t.Fatalf("Path = %q, want empty", got)
	}
	if got := p.Extension(); got != "prisma" {
		t.Fatalf("Extension = %q", got)
	}
	if got := p.ResourceType(); got != puck.GlobFileResource {
		t.Fatalf("ResourceType = %q", got)
	}

	content, err := p.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !strings.Contains(string(content), "aaa") || !strings.Contains(string(content), "bbbbb") {
		t.Fatalf("Read = %q", content)
	}
	if size := p.Size(); size != 8 {
		t.Fatalf("Size = %d, want 8", size)
	}
}

func TestPuck_GlobNoMatches(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "*.prisma"))
	if _, err := p.Read(); err == nil || !strings.Contains(err.Error(), "matched no files") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPuck_InvalidGlobRead(t *testing.T) {
	p := fromPuck(t, "/tmp/[bad*.prisma")
	if _, err := p.Read(); err == nil || !strings.Contains(err.Error(), "invalid glob") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPuck_GlobMatchesUnreadable(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "a.prisma"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFile(t, filepath.Join(dir, "b.prisma"), "bbb")

	p := fromPuck(t, filepath.Join(dir, "*.prisma"))
	if _, err := p.Read(); err == nil || !strings.Contains(err.Error(), "failed reading") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPuck_InvalidGlobSizePanics(t *testing.T) {
	p := fromPuck(t, "/tmp/[bad*.prisma")
	mustPanic(t, func() { _ = p.Size() })
}

func TestPuck_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gone.prisma")
	p := fromPuck(t, path)

	if p.Exists() {
		t.Fatal("Exists = true, want false")
	}
	if got := p.ResourceType(); got != puck.UnknownResource {
		t.Fatalf("ResourceType = %q", got)
	}
	if _, err := p.Read(); err == nil {
		t.Fatal("Read should error")
	}
	mustPanic(t, func() { _ = p.Size() })
}

func TestPuck_UnreadableFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; permission checks are bypassed")
	}
	path := filepath.Join(t.TempDir(), "secret.prisma")
	writeFile(t, path, "locked")
	if err := os.Chmod(path, 0); err != nil {
		t.Fatalf("chmod: %v", err)
	}

	p := fromPuck(t, path)
	if _, err := p.Read(); err == nil || !strings.Contains(err.Error(), "failed reading") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPuck_ReadGlobDeprecated(t *testing.T) {
	p := fromPuck(t, filepath.Join(t.TempDir(), "x.prisma"))
	if _, err := p.ReadGlob("*.prisma"); err == nil {
		t.Fatal("ReadGlob should return the deprecated error")
	}
}

func TestPuck_FormatOfDirectories(t *testing.T) {
	t.Run("no extension", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "noext")
		if err := os.Mkdir(sub, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if got := fromPuck(t, sub).Format(); got != puck.FormatUnknown {
			t.Fatalf("Format = %q, want unknown", got)
		}
	})
	t.Run("unknown extension", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "dir.txt")
		if err := os.Mkdir(sub, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if got := fromPuck(t, sub).Format(); got != puck.FormatUnknown {
			t.Fatalf("Format = %q, want unknown", got)
		}
	})
}

func TestPuck_StatErrorsPanic(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(sub, "x.prisma")
	writeFile(t, path, "v1")
	p := fromPuck(t, path)

	if err := os.RemoveAll(sub); err != nil {
		t.Fatalf("RemoveAll: %v", err)
	}
	writeFile(t, sub, "now a file")

	mustPanic(t, func() { _ = p.Exists() })
	mustPanic(t, func() { _ = p.ResourceType() })
}

func TestPuck_SizeDirectoryReadErrorPanics(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; permission checks are bypassed")
	}
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	p := fromPuck(t, sub)

	if err := os.Chmod(sub, 0); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	mustPanic(t, func() { _ = p.Size() })
}

func TestPuck_ConstructionPanics(t *testing.T) {
	t.Run("default NewPuck panics (unknown format)", func(t *testing.T) {
		mustPanic(t, func() { _ = puck.NewPuck() })
	})
	t.Run("relative path panics", func(t *testing.T) {
		mustPanic(t, func() { _ = puck.From("relative.txt") })
	})
	t.Run("glob without extension panics", func(t *testing.T) {
		mustPanic(t, func() { _ = puck.From("/tmp/noext*") })
	})
	t.Run("unknown format file panics", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "x.txt")
		writeFile(t, path, "hello")
		mustPanic(t, func() { _ = puck.From(path) })
	})
	t.Run("option error panics", func(t *testing.T) {
		mustPanic(t, func() { _ = puck.NewPuck(func(args ...any) error { return errors.New("boom") }) })
	})
}
