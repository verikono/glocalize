package errantly

import "fmt"

func WithMessage(message string, objs ...any) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*Errantly); ok {
				inst.Message = fmt.Sprintf(message, objs...)
			}
		}
		return nil
	}
}

func WithCause(cause error, objs ...any) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*Errantly); ok {
				inst.Cause = cause
			}
		}
		return nil
	}
}

func WithSourceLocation(location NotatedCall) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*Errantly); ok {
				inst.Location = &location
			}
		}
		return nil
	}
}
