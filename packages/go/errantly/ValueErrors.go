package errantly

type MissingValueError struct {
	Variable string
	ErrantlyError
}

func MissingValue(variable string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := MissingValueError{
		Variable:      variable,
		ErrantlyError: NewErrantlyError(options...),
	}
	// Only set the default message if the caller didn't provide one.
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("variable %s has no value.", variable)
	}
	return &inst
}

type InvalidValueError struct {
	Variable string
	ErrantlyError
}

func InvalidValue(variable string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := InvalidValueError{
		Variable:      variable,
		ErrantlyError: NewErrantlyError(options...),
	}
	// Only set the default message if the caller didn't provide one.
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("variable %s is invalid", variable)
	}
	return &inst
}

type MissingOptionError struct {
	Option string
	ErrantlyError
}

// MissingOption returns an error indicating a required functional option
// was not provided (e.g. WithBotToken, WithChannel).
func MissingOption(option string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := MissingOptionError{
		Option:        option,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("required option %s is missing", option)
	}
	return &inst
}
