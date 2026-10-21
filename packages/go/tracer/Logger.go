package tracer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// logTimestampLayout is the timestamp a Logger stamps on each entry. It carries
// the date so a log file read the next day is unambiguous.
const logTimestampLayout = "2006/01/02 15:04:05.000"

// Logger writes tracer-native log lines for one package. Every line it emits
// carries the package name as its origin, so a tracer can attribute the entry
// and a reader can filter on it.
//
// A Logger is safe for concurrent use. The zero value is not usable: build one
// with NewLogger.
type Logger struct {
	origin string
	out    io.Writer
	mu     sync.Mutex
}

// NewLogger builds a logger that stamps every line with origin. A nil out
// writes to os.Stdout, which is where a daemon's output is redirected from.
func NewLogger(origin string, out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{origin: origin, out: out}
}

// Origin returns the package name stamped on every line.
func (l *Logger) Origin() string {
	return l.origin
}

// Error writes an error-level entry.
func (l *Logger) Error(message string, details ...string) {
	l.write(ErrorLogLevel, message, details)
}

// Fatal writes a fatal-level entry. It does not exit: the caller owns that.
func (l *Logger) Fatal(message string, details ...string) {
	l.write(FatalLogLevel, message, details)
}

// Info writes an info-level entry.
func (l *Logger) Info(message string, details ...string) {
	l.write(InfoLogLevel, message, details)
}

// Trace writes a trace-level entry.
func (l *Logger) Trace(message string, details ...string) {
	l.write(TraceLogLevel, message, details)
}

// Warn writes a warn-level entry.
func (l *Logger) Warn(message string, details ...string) {
	l.write(WarnLogLevel, message, details)
}

// write emits one entry: a header line carrying the timestamp, level, origin
// and concise message, then one indented line per detail. A detail containing
// newlines is split so every physical line stays indented and therefore parses
// back as a detail; a message containing newlines is flattened onto one line,
// because a message is by definition concise and a newline would break the
// header.
func (l *Logger) write(level LogLevel, message string, details []string) {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s [%s] %s\n", time.Now().Format(logTimestampLayout), level, l.origin, flatten(message))
	for _, detail := range details {
		for _, line := range strings.Split(detail, "\n") {
			fmt.Fprintf(&b, "\t%s\n", line)
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = io.WriteString(l.out, b.String()) // a failed write to a log sink is not actionable
}

// flatten collapses every run of whitespace in s to a single space, so a value
// that must occupy one line cannot break the line format.
func flatten(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
