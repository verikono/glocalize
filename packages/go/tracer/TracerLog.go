package tracer

import (
	"regexp"
	"strings"
	"time"
)

// TracerLog is one parsed log line: when it happened, its level and its message.
type TracerLog struct {
	time    time.Time
	level   LogLevel
	message LogMessage
	TracerLogInterface
}

// NewTracerLog builds a log from its parts.
func NewTracerLog(ts time.Time, level LogLevel, message LogMessage) *TracerLog {
	return &TracerLog{time: ts, level: level, message: message}
}

// GetLevel returns the log's severity.
func (l *TracerLog) GetLevel() LogLevel {
	return l.level
}

// GetMessage returns the log's message text.
func (l *TracerLog) GetMessage() LogMessage {
	return l.message
}

// GetTime returns the log's timestamp truncated to whole epoch seconds.
func (l *TracerLog) GetTime() EpochSeconds {
	return EpochSeconds(l.time.Unix())
}

// GetTimestamp returns the log's full timestamp, including sub-second precision.
func (l *TracerLog) GetTimestamp() time.Time {
	return l.time
}

// timestampRx matches a leading timestamp in any of the layouts riot and the
// standard library emit, with or without surrounding brackets.
var timestampRx = regexp.MustCompile(
	`^\[?(\d{4}[/-]\d{2}[/-]\d{2}[ T]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?|\d{2}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}|\d{2}:\d{2}:\d{2}(?:\.\d+)?)\]?\s*`,
)

// levelRx matches a leading level token, bracketed or bare.
var levelRx = regexp.MustCompile(`(?i)^\[?(TRACE|DEBUG|INFO|WARN|WARNING|ERROR|FATAL|PANIC)\]?[:\s]+`)

// opRx matches a leading bracketed operation token, as the riot hot service
// emits: "[WRITE]:/path/to/file".
var opRx = regexp.MustCompile(`^\[[A-Z]+\]:\s*`)

// timestampLayouts are tried in order against a captured timestamp.
var timestampLayouts = []string{
	"2006/01/02 15:04:05.000",
	"2006/01/02 15:04:05",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02 15:04:05",
	"06/01/02 15:04:05",
	"15:04:05.000",
	"15:04:05",
}

// parseLogLine parses one raw line. A line with no recognisable timestamp is
// stamped with now; a line with no recognisable level is treated as info.
func parseLogLine(line string, now time.Time) *TracerLog {
	log := &TracerLog{time: now, level: InfoLogLevel}
	rest := strings.TrimSpace(line)

	if match := timestampRx.FindStringSubmatch(rest); match != nil {
		if ts, ok := parseTimestamp(match[1], now); ok {
			log.time = ts
			rest = strings.TrimSpace(rest[len(match[0]):])
		}
	}
	if match := levelRx.FindStringSubmatch(rest); match != nil {
		log.level = normaliseLevel(match[1])
		rest = strings.TrimSpace(rest[len(match[0]):])
	}
	if match := opRx.FindString(rest); match != "" {
		rest = strings.TrimSpace(rest[len(match):])
	}

	log.message = LogMessage(rest)
	return log
}

// parseTimestamp resolves a captured timestamp, anchoring a time-only value to
// the day of now.
func parseTimestamp(value string, now time.Time) (time.Time, bool) {
	for _, layout := range timestampLayouts {
		ts, err := time.ParseInLocation(layout, value, time.Local)
		if err != nil {
			continue
		}
		if !strings.Contains(layout, "2006") && !strings.Contains(layout, "06/") {
			return time.Date(now.Year(), now.Month(), now.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), time.Local), true
		}
		return ts, true
	}
	return time.Time{}, false
}

// normaliseLevel maps a raw level token onto the canonical LogLevel set.
func normaliseLevel(raw string) LogLevel {
	switch strings.ToUpper(raw) {
	case "TRACE":
		return TraceLogLevel
	case "DEBUG", "INFO":
		return InfoLogLevel
	case "WARN", "WARNING":
		return WarnLogLevel
	case "ERROR":
		return ErrorLogLevel
	case "FATAL", "PANIC":
		return FatalLogLevel
	default:
		return InfoLogLevel
	}
}
