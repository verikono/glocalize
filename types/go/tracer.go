package types

import (
	"context"
	"time"
)

type LogLevel string   // loglevel as a string "info", "warn" etc
type LogMessage string // a message for the logged event

// ReadableTracerLog is the read seam a tracer exposes: the whole buffer, the
// filtered view, and the ability to change the filter. Logs is unfiltered;
// Visible applies the current filter.
type ReadableTracerLog interface {
	Logs() []TracerLogInterface
	Visible() []TracerLogInterface
	SetLogFilter(options ...OptionerFn) ([]TracerLogInterface, error)
}

type TracerLogFilter interface {
	SetStart(epoch EpochSeconds) error
	SetEnd(epoch EpochSeconds) error
	SetLevel(level LogLevel) error
	SetFilterString(needle string) error
}

// TracerLogInterface is one log entry. A log carries two message types: a
// concise Message and an optional list of Details for the verbose story.
type TracerLogInterface interface {
	GetOrigin() string       // An identifier of origin, eg the package which created this.
	GetTime() EpochSeconds   // Time the log entry was created (in unix epoch time - seconds)
	GetTimestamp() time.Time // Timestamp instead of epoch time
	GetLevel() LogLevel      // warn, log, error, fatal - you know the drill
	GetMessage() LogMessage  // the concise message
	GetDetails() []string    // the details, the waffle, the extended story if any
}

// TracerLogSource yields log lines for a tracer to consume. Follow returns a
// channel that closes when ctx is cancelled or the source is exhausted.
type TracerLogSource interface {
	Follow(ctx context.Context) (<-chan string, error)
}

const (
	TraceLogLevel LogLevel = "trace"
	InfoLogLevel  LogLevel = "info"
	WarnLogLevel  LogLevel = "warn"
	ErrorLogLevel LogLevel = "error"
	FatalLogLevel LogLevel = "fatal"
)
