package tracer

import (
	"regexp"
	"strings"
	"time"
)

// TracerLog is one parsed log line: where it came from, when it happened, its
// level, its concise message and any verbose details.
//
// A log carries two message types. Message is the concise, succinct line a
// reader scans; Details is the optional list of extra lines a package may
// produce in any quantity, shown only when the log is focused.
type TracerLog struct {
	origin  string
	time    time.Time
	level   LogLevel
	message LogMessage
	details []string
}

// NewTracerLog builds a log from its parts, with no origin and no details.
func NewTracerLog(ts time.Time, level LogLevel, message LogMessage) *TracerLog {
	return &TracerLog{time: ts, level: level, message: message}
}

// NewTracerLogFrom builds a log with an origin and optional details. The origin
// identifies the package that produced the entry.
func NewTracerLogFrom(origin string, ts time.Time, level LogLevel, message LogMessage, details ...string) *TracerLog {
	return &TracerLog{origin: origin, time: ts, level: level, message: message, details: details}
}

// AddDetail appends one detail line to the log.
func (l *TracerLog) AddDetail(detail string) {
	l.details = append(l.details, detail)
}

// GetDetails returns the log's verbose detail lines, if any.
func (l *TracerLog) GetDetails() []string {
	return l.details
}

// GetLevel returns the log's severity.
func (l *TracerLog) GetLevel() LogLevel {
	return l.level
}

// GetMessage returns the log's concise message.
func (l *TracerLog) GetMessage() LogMessage {
	return l.message
}

// GetOrigin returns the identifier of the package that produced the log.
func (l *TracerLog) GetOrigin() string {
	return l.origin
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

// originRx matches a leading bracketed origin token, as a package emits it:
// "[orm] SQL Query: ...". The token must be followed by whitespace or the end
// of the line, so an operation token such as "[WRITE]:/path" is not mistaken
// for an origin.
var originRx = regexp.MustCompile(`^\[([A-Za-z][A-Za-z0-9_.\-/]*)\](?:\s+|$)`)

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

	// the level and the origin may appear in either order, so consume each at
	// most once, whichever comes first.
	var haveLevel, haveOrigin bool
	for {
		if !haveLevel {
			if match := levelRx.FindStringSubmatch(rest); match != nil {
				log.level = normaliseLevel(match[1])
				rest = strings.TrimSpace(rest[len(match[0]):])
				haveLevel = true
				continue
			}
		}
		if !haveOrigin {
			if match := originRx.FindStringSubmatch(rest); match != nil && !isLevelToken(match[1]) {
				log.origin = match[1]
				rest = strings.TrimSpace(rest[len(match[0]):])
				haveOrigin = true
				continue
			}
		}
		break
	}
	if match := opRx.FindString(rest); match != "" {
		rest = strings.TrimSpace(rest[len(match):])
	}

	log.message = LogMessage(rest)
	return log
}

// isLevelToken reports whether token names a level, so a bracketed level is
// never mistaken for an origin.
func isLevelToken(token string) bool {
	switch strings.ToUpper(token) {
	case "TRACE", "DEBUG", "INFO", "WARN", "WARNING", "ERROR", "FATAL", "PANIC":
		return true
	}
	return false
}

// isDetailLine reports whether a raw line continues the previous log rather
// than starting a new one. An indented line is a detail.
func isDetailLine(line string) bool {
	if line == "" {
		return false
	}
	return line[0] == ' ' || line[0] == '\t'
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
