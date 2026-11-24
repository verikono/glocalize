package errantly

// YamlParseError is a YAML-specific parse failure.
type YamlParseError struct {
	FilePath string
	ParseError
}

// YamlParseFailure creates a YamlParseError for the given file path.
func YamlParseFailure(filePath string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := YamlParseError{
		FilePath:    filePath,
		ParseError:  ParseError{ErrantlyError: NewErrantlyError(options...)},
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("failed to parse YAML file: %s", filePath)
	}
	return &inst
}
