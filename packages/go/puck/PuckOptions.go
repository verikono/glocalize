package puck

import (
	"strings"
)

func From(pattern string) IPuck {
	return NewPuck(from(pattern))
}

func from(pattern string) OptionerFn {
	return func(args ...any) error {
		for _, a := range args {
			puck, ok := a.(*Puck)
			if ok {
				if isPathPattern(pattern) {
					puck.url = "file://" + pattern
					puck.protocol = ProtocolFile
					if isGlobPattern((pattern)) {
						puck.glob = pattern
					}
				}
				// deprecated -- NewPuck will call this: puck.assertValidFormat()
			}
		}
		return nil
	}
}

func isPathPattern(value string) bool {
	return strings.HasPrefix(value, "/")
}

func isGlobPattern(value string) bool {
	return strings.Contains(value, "*")
}

// provide a listener which will be called when change occurs.
func WithWatcher(fn FileSystemWatcher) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if watcher, ok := arg.(*Watcher); ok {
				watcher.fn = fn
			}
		}
		return nil
	}
}
