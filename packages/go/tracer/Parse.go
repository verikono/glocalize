package tracer

import (
	"time"
)

// ParseLogLine parses one raw log line into a TracerLog. A line with no
// recognisable timestamp is stamped with now; a line with no recognisable level
// is treated as info.
func ParseLogLine(line string, now time.Time) *TracerLog {
	return parseLogLine(line, now)
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
