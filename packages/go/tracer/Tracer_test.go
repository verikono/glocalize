package tracer_test

import (
	"testing"
	"time"

	"github.com/verikono/glocalize/packages/go/tracer"
)

func TestParseFormat(t *testing.T) {
	ts := time.Date(2026, 10, 20, 20, 23, 33, 199_000_000, time.Local)
	log := tracer.NewTracerLog(ts, tracer.WarnLogLevel, "this is a warning")

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"default", "%t %l %m", "20:23:33.199 warn this is a warning"},
		{"reordered", "%l | %t | %m", "warn | 20:23:33.199 | this is a warning"},
		{"upper message", "%t %l %mU", "20:23:33.199 warn THIS IS A WARNING"},
		{"upper level", "%lU %m", "WARN this is a warning"},
		{"escaped percent", "100%% %m", "100% this is a warning"},
		{"empty selects default", "", "20:23:33.199 warn this is a warning"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestParseFormatRejectsUnknownOption(t *testing.T) {
	log := tracer.NewTracerLog(time.Now(), tracer.InfoLogLevel, "hi")
	if _, err := tracer.RenderFormat("%x %m", log); err == nil {
		t.Fatal("RenderFormat with an unknown option: error = nil, want an error")
	}
	if _, err := tracer.RenderFormat("%m %", log); err == nil {
		t.Fatal("RenderFormat with a trailing percent: error = nil, want an error")
	}
}

func TestParseLogLine(t *testing.T) {
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.Local)

	tests := []struct {
		name    string
		line    string
		level   tracer.LogLevel
		message string
		hour    int
	}{
		{"riot hot line", "[26/10/20 20:23:33] [WRITE]:/tmp/x.go", tracer.InfoLogLevel, "/tmp/x.go", 20},
		{"stdlib prefix", "2026/10/20 20:23:33 Running migrations...", tracer.InfoLogLevel, "Running migrations...", 20},
		{"bracketed level", "[ERROR] boom", tracer.ErrorLogLevel, "boom", 12},
		{"bare level", "warn careful", tracer.WarnLogLevel, "careful", 12},
		{"time only", "20:23:33.199 warn this is a warning", tracer.WarnLogLevel, "this is a warning", 20},
		{"no timestamp", "just a message", tracer.InfoLogLevel, "just a message", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := tracer.ParseLogLine(tt.line, now)
			if log.GetLevel() != tt.level {
				t.Fatalf("level = %q, want %q", log.GetLevel(), tt.level)
			}
			if log.GetMessage() != tracer.LogMessage(tt.message) {
				t.Fatalf("message = %q, want %q", log.GetMessage(), tt.message)
			}
			if got := log.GetTimestamp().Hour(); got != tt.hour {
				t.Fatalf("hour = %d, want %d", got, tt.hour)
			}
		})
	}
}

// TestNeedleTable is the truth table from SPEC.md, verbatim.
func TestNeedleTable(t *testing.T) {
	tests := []struct {
		message string
		needle  string
		want    bool
	}{
		{"jimbo loves to thunks in bed", "unk", true},
		{"jimbo loves to thunks in bed", "^jimbo*", true},
		{"jimbo loves to thunks in bed", " bed", true},
		{"jimbo loves to thunks in bed", "jim", true},
		{"jimbo loves to junkies in bed", "asdf", false},
		{"jimbo loves to junkies in bed", "^junk", false},
		{"jimbo loves to junkies in bed", " bedri", false},
		{"jimbo loves to junkies in bed", "john", false},
	}

	for _, tt := range tests {
		t.Run(tt.message+"/"+tt.needle, func(t *testing.T) {
			filter, err := tracer.NewTracerLogFilter()
			if err != nil {
				t.Fatalf("NewTracerLogFilter() error = %v", err)
			}
			if err := filter.SetFilterString(tt.needle); err != nil {
				t.Fatalf("SetFilterString(%q) error = %v", tt.needle, err)
			}
			log := tracer.NewTracerLog(time.Now(), tracer.InfoLogLevel, tracer.LogMessage(tt.message))
			if got := filter.Match(log); got != tt.want {
				t.Fatalf("Match(%q, %q) = %v, want %v", tt.message, tt.needle, got, tt.want)
			}
		})
	}
}

func TestFilterByLevelAndWindow(t *testing.T) {
	base := time.Date(2026, 10, 20, 12, 0, 0, 0, time.Local)
	early := tracer.NewTracerLog(base, tracer.InfoLogLevel, "early")
	late := tracer.NewTracerLog(base.Add(time.Hour), tracer.ErrorLogLevel, "late")

	filter, err := tracer.NewTracerLogFilter()
	if err != nil {
		t.Fatalf("NewTracerLogFilter() error = %v", err)
	}
	if err := filter.SetLevel(tracer.ErrorLogLevel); err != nil {
		t.Fatalf("SetLevel() error = %v", err)
	}
	if filter.Match(early) {
		t.Fatal("Match(early) = true, want false for an error-only filter")
	}
	if !filter.Match(late) {
		t.Fatal("Match(late) = false, want true for an error-only filter")
	}

	if err := filter.SetLevel(""); err != nil {
		t.Fatalf("SetLevel(\"\") error = %v", err)
	}
	if err := filter.SetStart(tracer.EpochSeconds(base.Add(30 * time.Minute).Unix())); err != nil {
		t.Fatalf("SetStart() error = %v", err)
	}
	if filter.Match(early) {
		t.Fatal("Match(early) = true, want false before the window start")
	}
	if !filter.Match(late) {
		t.Fatal("Match(late) = false, want true inside the window")
	}
}

func TestFilterRejectsBadRegex(t *testing.T) {
	filter, err := tracer.NewTracerLogFilter()
	if err != nil {
		t.Fatalf("NewTracerLogFilter() error = %v", err)
	}
	// An uncompilable needle is not an error: it degrades to a literal match.
	if err := filter.SetFilterString("/[unclosed"); err != nil {
		t.Fatalf("SetFilterString(/[unclosed) error = %v, want nil", err)
	}
	log := tracer.NewTracerLog(time.Now(), tracer.InfoLogLevel, "literal /[unclosed text")
	if !filter.Match(log) {
		t.Fatal("Match() = false, want true for a literal fallback match")
	}
}
