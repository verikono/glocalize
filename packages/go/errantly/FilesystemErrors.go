package errantly

type MissingFileError struct {
	Path string
	ErrantlyError
}

func MissingFile(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := MissingFileError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("file %s does not exist", path)
	}
	return &inst
}

type MissingDirectoryError struct {
	Path string
	ErrantlyError
}

func MissingDirectory(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := MissingDirectoryError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("directory %s does not exist", path)
	}
	return &inst
}

type FileNotReadableError struct {
	Path string
	ErrantlyError
}

func FileNotReadable(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := FileNotReadableError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("file %s is not readable", path)
	}
	return &inst
}

type FileNotWritableError struct {
	Path string
	ErrantlyError
}

func FileNotWritable(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := FileNotWritableError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("file %s is not writable", path)
	}
	return &inst
}

type FileNotExecutableError struct {
	Path string
	ErrantlyError
}

func FileNotExecutable(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := FileNotExecutableError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("file %s is not executable", path)
	}
	return &inst
}

type PermissionDeniedError struct {
	Path string
	ErrantlyError
}

func PermissionDenied(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := PermissionDeniedError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("permission denied for %s", path)
	}
	return &inst
}

type FileAlreadyExistsError struct {
	Path string
	ErrantlyError
}

func FileAlreadyExists(path string, options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	options = append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := FileAlreadyExistsError{
		Path:          path,
		ErrantlyError: NewErrantlyError(options...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("file %s already exists", path)
	}
	return &inst
}
