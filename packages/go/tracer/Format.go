package tracer

import (
	"fmt"
	"strings"
)

// DefaultFormat renders a log as "20:23:33.199 warn this is a warning".
const DefaultFormat = "%t %l %m"

// DefaultFormatWithOrigin renders a log with the package that produced it:
// "20:23:33.199 warn [orm] this is a warning".
const DefaultFormatWithOrigin = "%t %l [%o] %m"

// timestampLayout is the %t rendering: wall-clock time to the millisecond.
const timestampLayout = "15:04:05.000"

// formatSegment is one %option of a format string.
type formatSegment struct {
	kind  byte
	upper bool
}

// logFormat is a parsed format string: literals interleaved with segments, so
// len(literals) == len(segments)+1.
type logFormat struct {
	raw      string
	literals []string
	segments []formatSegment
}

// parseFormat compiles a format string. An empty string selects DefaultFormat.
func parseFormat(raw string) (logFormat, error) {
	if raw == "" {
		raw = DefaultFormat
	}

	format := logFormat{raw: raw}
	var literal strings.Builder

	for i := 0; i < len(raw); i++ {
		if raw[i] != '%' {
			literal.WriteByte(raw[i])
			continue
		}
		if i+1 >= len(raw) {
			return logFormat{}, fmt.Errorf("format %q: trailing %%", raw)
		}
		i++
		kind := raw[i]
		if kind == '%' {
			literal.WriteByte('%')
			continue
		}
		if kind != 't' && kind != 'l' && kind != 'm' && kind != 'o' {
			return logFormat{}, fmt.Errorf("format %q: unknown option %%%c", raw, kind)
		}
		segment := formatSegment{kind: kind}
		if i+1 < len(raw) && raw[i+1] == 'U' {
			segment.upper = true
			i++
		}
		format.literals = append(format.literals, literal.String())
		literal.Reset()
		format.segments = append(format.segments, segment)
	}

	format.literals = append(format.literals, literal.String())
	return format, nil
}

// String returns the format string as written.
func (f logFormat) String() string {
	return f.raw
}

// render applies the format to one log.
func (f logFormat) render(log TracerLogInterface) string {
	var out strings.Builder
	for i, segment := range f.segments {
		out.WriteString(f.literals[i])
		out.WriteString(segment.value(log))
	}
	out.WriteString(f.literals[len(f.segments)])
	return out.String()
}

// value renders one segment of a log.
func (s formatSegment) value(log TracerLogInterface) string {
	var value string
	switch s.kind {
	case 't':
		value = log.GetTimestamp().Format(timestampLayout)
	case 'l':
		value = string(log.GetLevel())
	case 'm':
		value = string(log.GetMessage())
	case 'o':
		value = log.GetOrigin()
	}
	if s.upper {
		value = strings.ToUpper(value)
	}
	return value
}
