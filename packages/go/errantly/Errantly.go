package errantly

import (
	"fmt"
	"strings"
)

type Errantly struct {
	Message  string
	Location *NotatedCall
	Cause    error
	Handlers []ErrantlyHandler
}

func (e *Errantly) GetCause() error {
	return e.Cause
}

func NewErrantlyError(options ...OptionerFn) ErrantlyError {
	inst := &Errantly{
		Message:  "An error has occured",
		Handlers: RegisteredHandlers(),
	}
	for _, option := range options {
		err := option(inst)
		if err != nil {
			panic(fmt.Sprintf("failed building error - %s", err))
		}
	}
	return inst
}

func (e *Errantly) GetMessage() string {
	return e.Message
}

func (e *Errantly) SetMessage(message string, objs ...any) {
	e.Message = fmt.Sprintf(message, objs...)
}

// Record hands this error to every registered handler.
//
// Handlers are attached at construction from the package-level registry, so a
// caller that never imports a recorder pays nothing. A handler failure is
// returned rather than swallowed: the caller asked for the error to be
// recorded, and a silent failure would look like success.
// @see .agents.md "## handlers"
func (e *Errantly) Record() error {
	var failures []string
	for _, handler := range e.Handlers {
		if err := handler.Record(e.GetMessage(), ErrantlyMessageMeta{
			Message: e.GetMessage(),
		}); err != nil {
			failures = append(failures, err.Error())
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("errantly: %d handler(s) failed to record: %s", len(failures), strings.Join(failures, "; "))
	}
	return nil
}

func (e *Errantly) GetSourceLocation() *NotatedCall {
	return e.Location
}

func (e *Errantly) Error() string {
	return e.Message
}
