package errantly

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
)

func GetNotatedCall(skip int) *NotatedCall {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return &NotatedCall{
			Line:     -1,
			FilePath: "unknown",
			Function: "unknown",
		}
	}
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		return &NotatedCall{
			Line:     line,
			FilePath: filepath.Base(file),
			Function: fn.Name(),
		}
	}
	return &NotatedCall{
		Line:     line,
		FilePath: filepath.Base(file),
		Function: "unknown",
	}
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(MustFailure())
	}
	return value
}

// findEmbeddedIssue returns a pointer to the first value field of type I
// embedded in target (a pointer to a struct), falling back to target itself
// when it is already of type *I. Issue options (WithURL, WithBytes,
// WithEnvVar, ...) use it so they can enrich both standalone issues and
// composite error types that embed them.
func findEmbeddedIssue[I any](target any) *I {
	if issue, ok := target.(*I); ok {
		return issue
	}
	v := reflect.ValueOf(target)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}
	sv := v.Elem()
	if sv.Kind() != reflect.Struct {
		return nil
	}
	issueType := reflect.TypeOf((*I)(nil)).Elem()
	for i := 0; i < sv.NumField(); i++ {
		field := sv.Field(i)
		if field.Type() == issueType && field.CanAddr() {
			return field.Addr().Interface().(*I)
		}
	}
	return nil
}

// applyIssueOptions runs options against a composite error so issue options
// (WithURL, WithBytes, WithEnvVar, ...) can reach embedded issue fields that
// the inner Errantly cannot see.
func applyIssueOptions(inst any, options []OptionerFn) {
	for _, opt := range options {
		if err := opt(inst); err != nil {
			panic(fmt.Sprintf("failed building error - %s", err))
		}
	}
}
