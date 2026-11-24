package errantly

func Provider() *ErrantlyPackage {
	return &ErrantlyPackage{
		NewErrantlyError:      NewErrantlyError,
		NotImplemented:        NotImplemented,
		MustFailure:           MustFailure,
		WithMessage:           WithMessage,
		WithCause:             WithCause,
		InvalidValue:          InvalidValue,
		MissingValue:          MissingValue,
		MissingOption:         MissingOption,
		MissingFile:           MissingFile,
		MissingDirectory:      MissingDirectory,
		FileNotReadable:       FileNotReadable,
		FileNotWritable:       FileNotWritable,
		FileNotExecutable:     FileNotExecutable,
		PermissionDenied:      PermissionDenied,
		FileAlreadyExists:     FileAlreadyExists,
		JsonParseFailure:      JsonParseFailure,
		YamlParseFailure:      YamlParseFailure,
		InvalidHttpRequest:    InvalidHttpRequest,
		HttpSendFailure:       HttpSendFailure,
		InitializationFailure: InitializationFailure,
		MissingEnvVar:         MissingEnvVar,
		WithURL:               WithURL,
		WithBytes:             WithBytes,
		WithEnvVar:            WithEnvVar,
	}
}

var J *JournalPackage
var Log Journalist
var handlers []ErrantlyHandler

func RegisterHandler(handler ErrantlyHandler) {
	if handler == nil {
		return
	}
	for _, existing := range handlers {
		if existing == handler {
			return
		}
	}
	handlers = append(handlers, handler)
}

func RegisteredHandlers() []ErrantlyHandler {
	out := make([]ErrantlyHandler, len(handlers))
	copy(out, handlers)
	return out
}

func RegisterDeps(deps *ErrantlyDependencies) error {
	J = deps.JournalPackage
	Log = Must(J.NewJournal("riot.errantly"))
	return nil
}
