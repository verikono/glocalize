package tracer_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/verikono/glocalize/packages/go/tracer"
)

// TestParseLogLineOrigin proves a leading bracketed token that is not a level
// is read as the origin, in either order relative to the level.
func TestParseLogLineOrigin(t *testing.T) {
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.Local)

	tests := []struct {
		name    string
		line    string
		origin  string
		level   tracer.LogLevel
		message string
	}{
		{
			name:    "origin after level",
			line:    "20:23:33.199 info [orm] sql query",
			origin:  "orm",
			level:   tracer.InfoLogLevel,
			message: "sql query",
		},
		{
			name:    "origin before level",
			line:    "20:23:33.199 [orm] warn slow query",
			origin:  "orm",
			level:   tracer.WarnLogLevel,
			message: "slow query",
		},
		{
			name:    "origin only",
			line:    "[mcp] tool registered",
			origin:  "mcp",
			level:   tracer.InfoLogLevel,
			message: "tool registered",
		},
		{
			name:    "bracketed level is not an origin",
			line:    "[ERROR] boom",
			origin:  "",
			level:   tracer.ErrorLogLevel,
			message: "boom",
		},
		{
			name:    "operation token is not an origin",
			line:    "[26/10/20 20:23:33] [WRITE]:/tmp/x.go",
			origin:  "",
			level:   tracer.InfoLogLevel,
			message: "/tmp/x.go",
		},
		{
			name:    "no origin",
			line:    "20:23:33.199 warn this is a warning",
			origin:  "",
			level:   tracer.WarnLogLevel,
			message: "this is a warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := tracer.ParseLogLine(tt.line, now)
			if log.GetOrigin() != tt.origin {
				t.Fatalf("origin = %q, want %q", log.GetOrigin(), tt.origin)
			}
			if log.GetLevel() != tt.level {
				t.Fatalf("level = %q, want %q", log.GetLevel(), tt.level)
			}
			if log.GetMessage() != tracer.LogMessage(tt.message) {
				t.Fatalf("message = %q, want %q", log.GetMessage(), tt.message)
			}
		})
	}
}

// TestParseLogLinesAttachesDetails proves an indented line becomes a detail of
// the log above it rather than a new entry.
func TestParseLogLinesAttachesDetails(t *testing.T) {
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.Local)
	logs := tracer.ParseLogLines([]string{
		"20:23:33.199 info [orm] sql query",
		"\tSELECT * FROM Agent",
		"\targs: []",
		"20:23:34.000 info [orm] sql exec",
	}, now)

	if len(logs) != 2 {
		t.Fatalf("logs = %d, want 2", len(logs))
	}
	details := logs[0].GetDetails()
	if len(details) != 2 {
		t.Fatalf("details = %d, want 2: %v", len(details), details)
	}
	if details[0] != "SELECT * FROM Agent" || details[1] != "args: []" {
		t.Fatalf("details = %v, want the two indented lines", details)
	}
	if len(logs[1].GetDetails()) != 0 {
		t.Fatalf("second log details = %v, want none", logs[1].GetDetails())
	}
}

// TestLoggerEmitsOriginAndDetails proves a Logger writes a line a tracer parses
// back into the same origin, level, message and details.
func TestLoggerEmitsOriginAndDetails(t *testing.T) {
	var buf bytes.Buffer
	logger := tracer.NewLogger("orm", &buf)

	logger.Warn("slow query", "SELECT 1", "took 900ms")
	logger.Info("multi\nline message", "first\nsecond")

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	if len(lines) != 6 {
		t.Fatalf("lines = %d, want 6:\n%s", len(lines), buf.String())
	}

	logs := tracer.ParseLogLines(lines, time.Now())
	if len(logs) != 2 {
		t.Fatalf("parsed logs = %d, want 2", len(logs))
	}
	if logs[0].GetOrigin() != "orm" {
		t.Fatalf("origin = %q, want orm", logs[0].GetOrigin())
	}
	if logs[0].GetLevel() != tracer.WarnLogLevel {
		t.Fatalf("level = %q, want warn", logs[0].GetLevel())
	}
	if logs[0].GetMessage() != "slow query" {
		t.Fatalf("message = %q, want %q", logs[0].GetMessage(), "slow query")
	}
	if got := logs[0].GetDetails(); len(got) != 2 || got[0] != "SELECT 1" || got[1] != "took 900ms" {
		t.Fatalf("details = %v, want [SELECT 1, took 900ms]", got)
	}
	// a message containing a newline is flattened onto one line.
	if logs[1].GetMessage() != "multi line message" {
		t.Fatalf("message = %q, want %q", logs[1].GetMessage(), "multi line message")
	}
	// a detail containing a newline is split so every line stays a detail.
	if got := logs[1].GetDetails(); len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Fatalf("details = %v, want [first, second]", got)
	}
}

// TestRenderFormatOrigin proves %o renders the origin and honours the U suffix.
func TestRenderFormatOrigin(t *testing.T) {
	ts := time.Date(2026, 10, 20, 20, 23, 33, 199_000_000, time.Local)
	log := tracer.NewTracerLogFrom("orm", ts, tracer.WarnLogLevel, "slow query")

	tests := []struct {
		format string
		want   string
	}{
		{"%t %l [%o] %m", "20:23:33.199 warn [orm] slow query"},
		{"%oU %m", "ORM slow query"},
		{"%m", "slow query"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got, err := tracer.RenderFormat(tt.format, log)
			if err != nil {
				t.Fatalf("RenderFormat(%q) error = %v", tt.format, err)
			}
			if got != tt.want {
				t.Fatalf("RenderFormat(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

// TestFilterMatchesOriginAndDetails proves the needle searches the origin and
// the detail lines, not just the message.
func TestFilterMatchesOriginAndDetails(t *testing.T) {
	log := tracer.NewTracerLogFrom("orm", time.Now(), tracer.InfoLogLevel, "sql query", "SELECT * FROM Agent")

	tests := []struct {
		needle string
		want   bool
	}{
		{"orm", true},
		{"SELECT", true},
		{"agent", true},
		{"mcp", false},
	}

	for _, tt := range tests {
		t.Run(tt.needle, func(t *testing.T) {
			filter, err := tracer.NewTracerLogFilter()
			if err != nil {
				t.Fatalf("NewTracerLogFilter() error = %v", err)
			}
			if err := filter.SetFilterString(tt.needle); err != nil {
				t.Fatalf("SetFilterString(%q) error = %v", tt.needle, err)
			}
			if got := filter.Match(log); got != tt.want {
				t.Fatalf("Match(%q) = %v, want %v", tt.needle, got, tt.want)
			}
		})
	}
}

// TestTracerAttachesDetailsAndFocusShowsThem proves the viewer folds an indented
// line into the log above it and renders the details only when focused.
func TestTracerAttachesDetailsAndFocusShowsThem(t *testing.T) {
	viewer := newTestTracer(t,
		"20:23:33.199 info [orm] sql query",
		"\tSELECT * FROM Agent",
		"20:23:34.000 info [orm] sql exec",
	)
	drain(t, viewer)

	logs := viewer.Logs()
	if len(logs) != 2 {
		t.Fatalf("logs = %d, want 2", len(logs))
	}
	if logs[0].GetOrigin() != "orm" {
		t.Fatalf("origin = %q, want orm", logs[0].GetOrigin())
	}
	if got := logs[0].GetDetails(); len(got) != 1 || got[0] != "SELECT * FROM Agent" {
		t.Fatalf("details = %v, want [SELECT * FROM Agent]", got)
	}

	// the detail is not in the log pane until the log is focused.
	if strings.Contains(viewer.View(), "SELECT * FROM Agent") {
		t.Fatal("View() shows the detail before the log is focused")
	}
	viewer.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !strings.Contains(viewer.View(), "SELECT * FROM Agent") {
		t.Fatalf("View() after enter does not show the detail:\n%s", viewer.View())
	}
}
