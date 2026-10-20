package tracer_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/verikono/glocalize/packages/go/tracer"
)

// writeFile writes content to path, failing the test on error.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// appendFile appends content to path, failing the test on error.
func appendFile(t *testing.T, path, content string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("append %s: %v", path, err)
	}
}

// fakeSource replays a fixed set of lines, then closes.
type fakeSource struct {
	lines []string
}

func (s *fakeSource) Follow(ctx context.Context) (<-chan string, error) {
	lines := make(chan string, len(s.lines))
	for _, line := range s.lines {
		lines <- line
	}
	close(lines)
	return lines, nil
}

func newTestTracer(t *testing.T, lines ...string) *tracer.Tracer {
	t.Helper()
	viewer, err := tracer.NewTracer(tracer.WithApp("riot"), tracer.WithSource(&fakeSource{lines: lines}))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	return viewer
}

// drain feeds every line from the source into the model.
func drain(t *testing.T, viewer *tracer.Tracer) {
	t.Helper()
	cmd := viewer.Init()
	for cmd != nil {
		msg := cmd()
		if _, ok := msg.(tracer.SourceDoneMsg); ok {
			return
		}
		_, cmd = viewer.Update(msg)
	}
}

func TestTracerBuffersAndRenders(t *testing.T) {
	viewer := newTestTracer(t,
		"20:23:33.199 warn this is a warning",
		"20:23:34.000 info all good",
	)
	drain(t, viewer)

	if got := len(viewer.Logs()); got != 2 {
		t.Fatalf("len(Logs()) = %d, want 2", got)
	}

	view := viewer.View()
	for _, want := range []string{"riot", "20:23:33.199 warn this is a warning", "[Q]uit", "[F]ilter", "[C]lear"} {
		if !strings.Contains(view, want) {
			t.Fatalf("View() missing %q:\n%s", want, view)
		}
	}
}

func TestTracerFilterAndClear(t *testing.T) {
	viewer := newTestTracer(t,
		"20:23:33.199 warn this is a warning",
		"20:23:34.000 info all good",
	)
	drain(t, viewer)

	visible, err := viewer.SetLogFilter(tracer.WithNeedle("warn"))
	if err != nil {
		t.Fatalf("SetLogFilter() error = %v", err)
	}
	if len(visible) != 1 {
		t.Fatalf("len(visible) = %d, want 1", len(visible))
	}

	// "c" clears the filter.
	viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if got := len(viewer.Visible()); got != 2 {
		t.Fatalf("len(Visible()) after clear = %d, want 2", got)
	}
}

func TestTracerFocusAndUnfocus(t *testing.T) {
	viewer := newTestTracer(t, "20:23:33.199 warn this is a warning")
	drain(t, viewer)

	viewer.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !strings.Contains(viewer.View(), "this is a warning") {
		t.Fatal("View() after enter does not show the focused log")
	}

	viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if strings.Count(viewer.View(), "this is a warning") != 1 {
		t.Fatalf("View() after x should show the log once:\n%s", viewer.View())
	}
}

func TestTracerFilterUI(t *testing.T) {
	viewer := newTestTracer(t,
		"20:23:33.199 warn this is a warning",
		"20:23:34.000 info all good",
	)
	drain(t, viewer)

	// "f" opens the one-line needle editor.
	viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	if !strings.Contains(viewer.View(), "filter:") {
		t.Fatalf("View() after f does not show the filter editor:\n%s", viewer.View())
	}
	if !strings.Contains(viewer.View(), "[Ctrl-F]ilter") {
		t.Fatalf("View() after f does not show the active filter hint:\n%s", viewer.View())
	}

	// Typing narrows the visible logs live.
	for _, r := range "warn" {
		viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	if got := len(viewer.Visible()); got != 1 {
		t.Fatalf("len(Visible()) after typing warn = %d, want 1", got)
	}

	// Enter closes the editor and keeps the filter.
	viewer.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if strings.Contains(viewer.View(), "filter:") {
		t.Fatalf("View() after enter still shows the filter editor:\n%s", viewer.View())
	}
	if got := len(viewer.Visible()); got != 1 {
		t.Fatalf("len(Visible()) after enter = %d, want 1", got)
	}
}

func TestTracerQuitsOnQ(t *testing.T) {
	viewer := newTestTracer(t)
	_, cmd := viewer.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("Update(q) returned no command, want tea.Quit")
	}
	if msg := cmd(); msg == nil {
		t.Fatal("quit command returned nil message")
	}
}

func TestTracerRequiresSource(t *testing.T) {
	if _, err := tracer.NewTracer(); err == nil {
		t.Fatal("NewTracer() with no source: error = nil, want an error")
	}
}

func TestFileSourceFollowsAppendedLines(t *testing.T) {
	path := t.TempDir() + "/riot.log"
	writeFile(t, path, "first line\n")

	source := tracer.NewFileSource(path)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lines, err := source.Follow(ctx)
	if err != nil {
		t.Fatalf("Follow() error = %v", err)
	}

	appendFile(t, path, "second line\n")

	select {
	case line := <-lines:
		if line != "second line" {
			t.Fatalf("line = %q, want %q", line, "second line")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for the appended line")
	}
}
