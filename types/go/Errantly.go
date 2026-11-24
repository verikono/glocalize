package types

type ErrantlyDependencies struct {
	CoreDependencies
}

type ErrantlyError interface {
	GetMessage() string
	GetCause() error
	SetMessage(message string, objs ...any)
	GetSourceLocation() *NotatedCall
	Error() string
}

type ErrantlyPackage struct {
	NewErrantlyError      func(options ...OptionerFn) ErrantlyError
	WithMessage           func(message string, objs ...any) OptionerFn
	WithCause             func(cause error, objs ...any) OptionerFn
	NotImplemented        func(options ...OptionerFn) ErrantlyError
	MustFailure           func(options ...OptionerFn) ErrantlyError
	MissingValue          func(variable string, options ...OptionerFn) ErrantlyError
	InvalidValue          func(variable string, options ...OptionerFn) ErrantlyError
	MissingOption         func(option string, options ...OptionerFn) ErrantlyError
	MissingFile           func(path string, options ...OptionerFn) ErrantlyError
	MissingDirectory      func(path string, options ...OptionerFn) ErrantlyError
	FileNotReadable       func(path string, options ...OptionerFn) ErrantlyError
	FileNotWritable       func(path string, options ...OptionerFn) ErrantlyError
	FileNotExecutable     func(path string, options ...OptionerFn) ErrantlyError
	PermissionDenied      func(path string, options ...OptionerFn) ErrantlyError
	FileAlreadyExists     func(path string, options ...OptionerFn) ErrantlyError
	JsonParseFailure      func(payload any) ErrantlyError
	YamlParseFailure      func(filePath string, options ...OptionerFn) ErrantlyError
	InvalidHttpRequest    func(options ...OptionerFn) ErrantlyError
	HttpSendFailure       func(options ...OptionerFn) ErrantlyError
	InitializationFailure func(options ...OptionerFn) ErrantlyError
	MissingEnvVar         func(options ...OptionerFn) ErrantlyError
	WithURL               func(url string) OptionerFn
	WithBytes             func(data []byte) OptionerFn
	WithEnvVar            func(variable string) OptionerFn
}

type ErrantlyMessageMeta struct {
	Message string
	Level   string
}

type ErrantlyHandler interface {
	Record(message string, meta ErrantlyMessageMeta) error
}
