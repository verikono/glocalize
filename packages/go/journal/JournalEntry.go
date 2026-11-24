package journal

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

type JournalEntry struct {
	timestamp  EpochNanoSeconds
	level      SeverityLevel
	message    string
	caller     string // `${packagename}/${function}:${line}`
	stacktrace string // `not implemented for now.`
}

func NewJournalEntry(level SeverityLevel, message string) *JournalEntry {
	return &JournalEntry{
		timestamp:  EpochNanoSeconds(time.Now().UnixNano()),
		level:      level,
		message:    message,
		caller:     captureCaller(2), // skip NewJournalEntry + captureCaller
		stacktrace: "not implemented for now.",
	}
}

// -- JournalEntryReader --

func (e *JournalEntry) GetLevel() SeverityLevel { return e.level }
func (e *JournalEntry) GetMessage() string      { return e.message }
func (e *JournalEntry) GetTimestamp() int64     { return int64(e.timestamp) }
func (e *JournalEntry) GetCaller() string       { return e.caller }
func (e *JournalEntry) GetStackTrace() string   { return e.stacktrace }

func captureCaller(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown:0"
	}
	name := fn.Name()
	// extract `${package}/${function}` from full path like `kaltra.io/path/to/package.FuncName`
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return name + ":" + strconv.Itoa(line)
}
