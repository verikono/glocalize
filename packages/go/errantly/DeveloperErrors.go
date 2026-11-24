package errantly

type NotImplementedError struct {
	ErrantlyError
}

func NotImplemented(options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	opts := append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := NotImplementedError{
		ErrantlyError: NewErrantlyError(opts...),
	}
	return &inst
}

type MustAssertionError struct {
	ErrantlyError
}

func MustFailure(options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := NotImplementedError{
		ErrantlyError: NewErrantlyError(options...),
	}
	return &inst
}
