package tracer

import (
	"fmt"
)

// WithApp sets the application name shown in the title bar.
func WithApp(app string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if tracer, ok := arg.(*Tracer); ok {
				tracer.app = app
			}
		}
		return nil
	}
}

// WithFormat sets the render format, e.g. "%l | %t | %m". An unknown option is
// an error, so a bad format fails at construction rather than at render time.
func WithFormat(raw string) OptionerFn {
	return func(args ...any) error {
		format, err := parseFormat(raw)
		if err != nil {
			return err
		}
		for _, arg := range args {
			if tracer, ok := arg.(*Tracer); ok {
				tracer.format = format
			}
		}
		return nil
	}
}

// WithSource sets the log source the tracer follows.
func WithSource(source TracerLogSource) OptionerFn {
	return func(args ...any) error {
		if source == nil {
			return fmt.Errorf("tracer: nil log source")
		}
		for _, arg := range args {
			if tracer, ok := arg.(*Tracer); ok {
				tracer.source = source
			}
		}
		return nil
	}
}

// WithNeedle seeds the filter with a needle. It applies to a *Tracer or a
// *TracerLogFilter, so the same option works at construction and in SetLogFilter.
func WithNeedle(needle string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			switch target := arg.(type) {
			case *Tracer:
				if err := target.filter.SetFilterString(needle); err != nil {
					return err
				}
				target.needle = needle
			case *TracerLogFilter:
				if err := target.SetFilterString(needle); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// WithLevel restricts the filter to one level. It applies to a *Tracer or a
// *TracerLogFilter.
func WithLevel(level LogLevel) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			switch target := arg.(type) {
			case *Tracer:
				if err := target.filter.SetLevel(level); err != nil {
					return err
				}
			case *TracerLogFilter:
				if err := target.SetLevel(level); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
