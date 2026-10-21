package tracer

import (
	"strings"
	"time"
)

// ParseLogLine parses one raw log line into a TracerLog. A line with no
// recognisable timestamp is stamped with now; a line with no recognisable level
// is treated as info. A leading bracketed token that is not a level is read as
// the origin, e.g. "[orm] SQL Query: ...".
func ParseLogLine(line string, now time.Time) *TracerLog {
	return parseLogLine(line, now)
}

// ParseLogLines parses a block of raw lines into logs, attaching each indented
// line to the log above it as a detail. It is the block form of what a tracer
// does as lines arrive.
func ParseLogLines(lines []string, now time.Time) []*TracerLog {
	logs := make([]*TracerLog, 0, len(lines))
	var last *TracerLog
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if isDetailLine(line) && last != nil {
			last.AddDetail(strings.TrimSpace(line))
			continue
		}
		last = parseLogLine(line, now)
		logs = append(logs, last)
	}
	return logs
}

// RenderFormat applies a format string to one log, e.g. "%l | %t | %m". An
// empty format selects DefaultFormat. An unknown option is an error.
func RenderFormat(raw string, log TracerLogInterface) (string, error) {
	format, err := parseFormat(raw)
	if err != nil {
		return "", err
	}
	return format.render(log), nil
}
