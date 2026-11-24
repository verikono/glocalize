package journal

import "io"

func WithHandler(handler string) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(Journalist); ok {
				if !HandlerExists(handler) {
					return E.NewErrantlyError(
						E.WithMessage("Handler '%s' is not registered", handler),
					)
				}
				inst.SetHandler(JournalHandlerKind(handler))
			}
		}
		return nil
	}
}

func WithoutConsole() OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*Journal); ok {
				inst.console = false
				return nil
			}
		}
		return nil
	}
}

// WithWriter sets the output writer on a ConsoleJournalHandler. Useful in
// tests to capture journal output.
func WithWriter(writer io.Writer) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if h, ok := target.(*ConsoleJournalHandler); ok {
				h.writer = writer
			}
		}
		return nil
	}
}
