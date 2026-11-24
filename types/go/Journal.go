package types

import "fmt"

type JournalPackage struct {
	NewJournal     func(name string, options ...OptionerFn) (Journalist, error)
	NewNoopJournal func() Journalist
}

type JournalDependencies struct {
	CoreDependencies
}

type JournalName string

type SeverityLevel int

const (
	SeverityLevelDebug    SeverityLevel = iota
	SeverityLevelTrace    SeverityLevel = 5
	SeverityLevelProlix   SeverityLevel = 10
	SeverityLevelVerbose  SeverityLevel = 20
	SeverityLevelInfo     SeverityLevel = 30
	SeverityLevelNotice   SeverityLevel = 40
	SeverityLevelWarn     SeverityLevel = 50
	SeverityLevelError    SeverityLevel = 60
	SeverityLevelCritical SeverityLevel = 70
	SeverityLevelFatal    SeverityLevel = 80
)

// String returns the canonical uppercase name of the severity level.
func (s SeverityLevel) String() string {
	switch s {
	case SeverityLevelDebug:
		return "DEBUG"
	case SeverityLevelTrace:
		return "TRACE"
	case SeverityLevelProlix:
		return "PROLIX"
	case SeverityLevelVerbose:
		return "VERBOSE"
	case SeverityLevelInfo:
		return "INFO"
	case SeverityLevelNotice:
		return "NOTICE"
	case SeverityLevelWarn:
		return "WARN"
	case SeverityLevelError:
		return "ERROR"
	case SeverityLevelCritical:
		return "CRITICAL"
	case SeverityLevelFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("SEVERITY(%d)", int(s))
	}
}

type JournalEntryReader interface {
	GetLevel() SeverityLevel
	GetMessage() string
	GetTimestamp() int64
	GetCaller() string
	GetStackTrace() string
}

type JournalEntryWriter interface {
}

type Journalist interface {
	GetKind() JournalHandlerKind
	GetName() JournalName
	SetHandler(handler JournalHandlerKind) error

	// Severity-ordered output methods. The first argument is a fmt-style
	// format string; remaining arguments are formatted into it, e.g.
	// Info("connected to %s (gateway: %s)", ssid, gateway).
	Critical(format string, args ...any)
	Error_(format string, args ...any)
	Warn(format string, args ...any)
	Notice(format string, args ...any)
	Info(format string, args ...any)
	Debug(format string, args ...any)
	Trace(format string, args ...any)
	Verbose(format string, args ...any)
	Prolix(format string, args ...any)
	Log(format string, args ...any)
}

type JournalEntryProcessor interface {
	JournalEntryReader
	JournalEntryWriter
}

type JournalHandlerKind string
type JournalHandlerProvider = func(options ...OptionerFn) (JournalHandler, error)

type JournalHandler interface {
	GetKind() JournalHandlerKind
	ProcessEntry(entry JournalEntryProcessor) error
}

func (j *JournalName) Validate() error {
	// journal name must be a-z, periods, and hyphens only
	if j == nil {
		return nil
	}
	name := string(*j)
	if len(name) == 0 {
		return nil
	}
	for _, c := range name {
		if (c < 'a' || c > 'z') && c != '.' && c != '-' {
			return fmt.Errorf("invalid journal name: %s", name)
		}
	}
	return nil
}

// NoopJournal is a Journalist that discards every message. Obtain one via
// NewNoopJournal; it is the safe fallback when no real journal is wired in.
type NoopJournal struct{}

func NewNoopJournal() Journalist { return NoopJournal{} }

func (NoopJournal) GetKind() JournalHandlerKind         { return "" }
func (NoopJournal) GetName() JournalName                { return "" }
func (NoopJournal) SetHandler(JournalHandlerKind) error { return nil }
func (NoopJournal) Critical(string, ...any)             {}
func (NoopJournal) Error_(string, ...any)               {}
func (NoopJournal) Warn(string, ...any)                 {}
func (NoopJournal) Notice(string, ...any)               {}
func (NoopJournal) Info(string, ...any)                 {}
func (NoopJournal) Debug(string, ...any)                {}
func (NoopJournal) Trace(string, ...any)                {}
func (NoopJournal) Verbose(string, ...any)              {}
func (NoopJournal) Prolix(string, ...any)               {}
func (NoopJournal) Log(string, ...any)                  {}
