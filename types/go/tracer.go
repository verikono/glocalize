package types

import (
	"context"
	"time"
)

type LogLevel string   // loglevel as a string "info", "warn" etc
type LogMessage string // a message for the logged event

type ReadableTracerLog interface {
	GetLogs(options ...OptionerFn) ([]TracerLogInterface, error)
	SetLogFilter(options ...OptionerFn) ([]TracerLogInterface, error)
}

type TracerLogFilter interface {
	SetStart(epoch EpochSeconds) error
	SetEnd(epoch EpochSeconds) error
	SetLevel(level LogLevel) error
	SetFilterString(needle string) error
}

type TracerLogInterface interface {
	GetTime() EpochSeconds
	GetTimestamp() time.Time
	GetLevel() LogLevel
	GetMessage() LogMessage
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
