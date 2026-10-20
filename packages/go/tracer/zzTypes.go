package tracer

import (
	T "github.com/verikono/glocalize/types/go"
)

type (
	EpochSeconds = T.EpochSeconds
	LogMessage   = T.LogMessage
)

type (
	OptionerFn = T.OptionerFn
)

type (
	LogLevel = T.LogLevel
)

type (
	TracerLogInterface       = T.TracerLogInterface
	TracerLogFilterInterface = T.TracerLogFilter
	TracerLogSource          = T.TracerLogSource
	ReadableTracerLog        = T.ReadableTracerLog
)

const (
	TraceLogLevel = T.TraceLogLevel
	InfoLogLevel  = T.InfoLogLevel
	WarnLogLevel  = T.WarnLogLevel
	ErrorLogLevel = T.ErrorLogLevel
	FatalLogLevel = T.FatalLogLevel
)
