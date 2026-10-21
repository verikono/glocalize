package tracer

import (
	"regexp"
	"strings"
)

// TracerLogFilter narrows a log stream by time window, level and needle.
//
// A needle is a plain substring match against every field of a log (time,
// level, origin, message and details) unless it begins with "/", in which case
// the remainder is compiled as a regular expression. Matching is
// case-insensitive either way.
type TracerLogFilter struct {
	start  EpochSeconds
	end    EpochSeconds
	level  LogLevel
	needle string
	rx     *regexp.Regexp
}

// NewTracerLogFilter builds a filter from options.
func NewTracerLogFilter(options ...OptionerFn) (*TracerLogFilter, error) {
	filter := &TracerLogFilter{}
	for _, option := range options {
		if err := option(filter); err != nil {
			return nil, err
		}
	}
	return filter, nil
}

// SetEnd bounds the window at epoch (inclusive). Zero means unbounded.
func (f *TracerLogFilter) SetEnd(epoch EpochSeconds) error {
	f.end = epoch
	return nil
}

// SetFilterString sets the needle.
//
// SPEC.md's prose says only a "/"-prefixed needle is a regular expression, but
// its truth table requires "^jimbo*" (unprefixed) to match "jimbo loves to
// thunks in bed" — which is only true under regex semantics. The table is the
// acceptance criterion, so every needle is compiled as a case-insensitive
// regex; a leading "/" is stripped first, and a needle that will not compile
// falls back to a literal substring match.
func (f *TracerLogFilter) SetFilterString(needle string) error {
	f.needle = needle
	f.rx = nil
	if needle == "" {
		return nil
	}
	pattern := strings.TrimPrefix(needle, "/")
	rx, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil // not a valid regex: matchNeedle falls back to a literal substring
	}
	f.rx = rx
	return nil
}

// SetLevel restricts the filter to one level. An empty level matches all.
func (f *TracerLogFilter) SetLevel(level LogLevel) error {
	f.level = level
	return nil
}

// SetStart bounds the window at epoch (inclusive). Zero means unbounded.
func (f *TracerLogFilter) SetStart(epoch EpochSeconds) error {
	f.start = epoch
	return nil
}

// Active reports whether the filter would exclude anything.
func (f *TracerLogFilter) Active() bool {
	return f.start != 0 || f.end != 0 || f.level != "" || f.needle != ""
}

// Needle returns the current needle.
func (f *TracerLogFilter) Needle() string {
	return f.needle
}

// Level returns the current level restriction.
func (f *TracerLogFilter) Level() LogLevel {
	return f.level
}

// Match reports whether log passes the filter.
func (f *TracerLogFilter) Match(log TracerLogInterface) bool {
	if f.start != 0 && log.GetTime() < f.start {
		return false
	}
	if f.end != 0 && log.GetTime() > f.end {
		return false
	}
	if f.level != "" && log.GetLevel() != f.level {
		return false
	}
	return f.matchNeedle(log)
}

// matchNeedle applies the needle to every field of the log: its timestamp,
// level, origin, message and each detail line.
func (f *TracerLogFilter) matchNeedle(log TracerLogInterface) bool {
	if f.needle == "" {
		return true
	}
	fields := []string{
		log.GetTimestamp().Format(timestampLayout),
		string(log.GetLevel()),
		log.GetOrigin(),
		string(log.GetMessage()),
	}
	fields = append(fields, log.GetDetails()...)
	if f.rx != nil {
		for _, field := range fields {
			if f.rx.MatchString(field) {
				return true
			}
		}
		return false
	}
	needle := strings.ToLower(f.needle)
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), needle) {
			return true
		}
	}
	return false
}
