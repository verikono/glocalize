package errantly

type InitializationFailureError struct {
	URLIssue
	ErrantlyError
}

func InitializationFailure(options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	opts := append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := &InitializationFailureError{
		ErrantlyError: NewErrantlyError(opts...),
	}
	applyIssueOptions(inst, options)
	// Only set the default message if the caller didn't provide one.
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("initialization error")
	}
	return inst
}

// MissingEnvVarError indicates a required environment variable is not set.
type MissingEnvVarError struct {
	Variable string
	ErrantlyError
}

// MissingEnvVar returns an error indicating a required environment variable is missing. Attach the variable name with WithEnvVar.
func MissingEnvVar(options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	opts := append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := &MissingEnvVarError{
		ErrantlyError: NewErrantlyError(opts...),
	}
	applyIssueOptions(inst, options)
	// Only set the default message if the caller didn't provide one.
	if inst.GetMessage() == "An error has occured" {
		if inst.Variable != "" {
			inst.SetMessage("environment variable %s is required but not set", inst.Variable)
		} else {
			inst.SetMessage("required environment variable is missing")
		}
	}
	return inst
}

// WithEnvVar attaches the environment variable name to an ErrantlyError.
func WithEnvVar(variable string) OptionerFn {
	return func(targets ...any) error {
		for _, target := range targets {
			if inst, ok := target.(*MissingEnvVarError); ok {
				inst.Variable = variable
			}
		}
		return nil
	}
}
