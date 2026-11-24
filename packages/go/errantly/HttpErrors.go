package errantly

// -- parse errors -- //

// ParseError is the base error type for data parsing failures.
type ParseError struct {
	ErrantlyError
}

// JsonParseError wraps a ParseError with the payload that failed to marshal.
type JsonParseError struct {
	ParseError
	Payload any
}

// JsonParseFailure returns a JsonParseError describing a JSON marshal
// failure for the given payload.
func JsonParseFailure(payload any) ErrantlyError {
	nc := GetNotatedCall(1)
	options := []OptionerFn{WithSourceLocation(*nc)}
	inst := &JsonParseError{
		ParseError: ParseError{
			ErrantlyError: NewErrantlyError(options...),
		},
		Payload: payload,
	}
	inst.SetMessage("json parse failure for payload of type %T", payload)
	return inst
}

// -- http errors -- //

// HttpError is the base error type for HTTP-level failures.
type HttpError struct {
	ErrantlyError
}

// NetworkError is the base error type for network transport failures.
type NetworkError struct {
	ErrantlyError
}

type BytesIssue struct {
	Bytes []byte
	ErrantlyError
}

type URLIssue struct {
	URL URLString
	ErrantlyError
}

// NetworkTransportError represents an HTTP transport failure, composing
// both HttpError and NetworkError context.
type NetworkTransportError struct {
	BytesIssue
	URLIssue
	ErrantlyError
}

// InvalidHttpRequest returns a NetworkTransportError describing a failure
// to construct a valid HTTP request.
func InvalidHttpRequest(options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	opts := append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := &NetworkTransportError{
		ErrantlyError: NewErrantlyError(opts...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("invalid http request")
	}
	return inst
}

// HttpSendFailure returns a NetworkTransportError describing a failure
// to send an HTTP request over the network.
func HttpSendFailure(options ...OptionerFn) ErrantlyError {
	nc := GetNotatedCall(1)
	opts := append([]OptionerFn{WithSourceLocation(*nc)}, options...)
	inst := &NetworkTransportError{
		ErrantlyError: NewErrantlyError(opts...),
	}
	if inst.GetMessage() == "An error has occured" {
		inst.SetMessage("http send failure")
	}
	return inst
}
