package errantly_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/verikono/glocalize/packages/go/errantly"
	types "github.com/verikono/glocalize/types/go"
)

func TestErrantly(t *testing.T) {

	t.Run("Instantiation", func(t *testing.T) {

		t.Run("using default values (no options)", func(t *testing.T) {
			err := errantly.NewErrantlyError()
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if err.Error() != "An error has occured" {
				t.Fatalf("expected default message, got %q", err.Error())
			}
			if err.GetSourceLocation() != nil {
				t.Fatal("expected nil source location with no options")
			}
		})

		t.Run("WithMessage :: custom error message", func(t *testing.T) {
			err := errantly.NewErrantlyError(
				errantly.WithMessage("something went wrong"),
			)
			if err.Error() != "something went wrong" {
				t.Fatalf("expected custom message, got %q", err.Error())
			}
		})

		t.Run("WithMessage :: formatted message with args", func(t *testing.T) {
			err := errantly.NewErrantlyError(
				errantly.WithMessage("user %s not found in %s", "bren", "database"),
			)
			expected := "user bren not found in database"
			if err.Error() != expected {
				t.Fatalf("expected %q, got %q", expected, err.Error())
			}
		})

		t.Run("WithCause :: wrapped error chain", func(t *testing.T) {
			cause := errors.New("network timeout")
			err := errantly.NewErrantlyError(
				errantly.WithMessage("failed to connect"),
				errantly.WithCause(cause),
			)
			if err.Error() != "failed to connect" {
				t.Fatalf("expected %q, got %q", "failed to connect", err.Error())
			}
			if err.GetCause() == nil {
				t.Fatal("expected GetCause() to return the cause")
			}
			if err.GetCause().Error() != "network timeout" {
				t.Fatalf("expected cause %q, got %q", "network timeout", err.GetCause().Error())
			}
		})

		t.Run("WithSourceLocation :: captures caller info", func(t *testing.T) {
			err := errantly.NewErrantlyError(
				errantly.WithSourceLocation(errantly.NotatedCall{
					FilePath: "/some/file.go",
					Line:     42,
					Function: "TestFunc",
				}),
			)
			loc := err.GetSourceLocation()
			if loc == nil {
				t.Fatal("expected non-nil source location")
			}
			if loc.FilePath != "/some/file.go" {
				t.Fatalf("expected /some/file.go, got %s", loc.FilePath)
			}
			if loc.Line != 42 {
				t.Fatalf("expected line 42, got %d", loc.Line)
			}
			if loc.Function != "TestFunc" {
				t.Fatalf("expected TestFunc, got %s", loc.Function)
			}
		})

		t.Run("combining multiple options", func(t *testing.T) {
			cause := errors.New("disk full")
			err := errantly.NewErrantlyError(
				errantly.WithMessage("write failed: %s", "no space"),
				errantly.WithCause(cause),
				errantly.WithSourceLocation(errantly.NotatedCall{
					FilePath: "/app/writer.go",
					Line:     99,
					Function: "WriteFile",
				}),
			)
			if err.Error() != "write failed: no space" {
				t.Fatalf("expected message %q, got %q", "write failed: no space", err.Error())
			}
			if err.GetCause().Error() != "disk full" {
				t.Fatal("expected cause to be reachable via GetCause()")
			}
			if err.GetSourceLocation().Line != 99 {
				t.Fatalf("expected line 99, got %d", err.GetSourceLocation().Line)
			}
		})
	})

	t.Run("GetMessage", func(t *testing.T) {
		err := errantly.NewErrantlyError(
			errantly.WithMessage("hello world"),
		)
		if err.GetMessage() != "hello world" {
			t.Fatalf("expected %q, got %q", "hello world", err.GetMessage())
		}
	})

	t.Run("SetMessage", func(t *testing.T) {
		err := errantly.NewErrantlyError()
		err.SetMessage("updated message")
		if err.GetMessage() != "updated message" {
			t.Fatalf("expected %q, got %q", "updated message", err.GetMessage())
		}
		if err.Error() != "updated message" {
			t.Fatalf("Error() should return the updated message")
		}
	})

	t.Run("SetMessage with format args", func(t *testing.T) {
		err := errantly.NewErrantlyError()
		err.SetMessage("code %d: %s", 500, "internal error")
		if err.GetMessage() != "code 500: internal error" {
			t.Fatalf("expected formatted message, got %q", err.GetMessage())
		}
	})

	t.Run("Error() returns message", func(t *testing.T) {
		err := errantly.NewErrantlyError(
			errantly.WithMessage("plain error"),
		)
		if err.Error() != "plain error" {
			t.Fatalf("expected %q, got %q", "plain error", err.Error())
		}
	})

	t.Run("GetSourceLocation returns nil for defaults", func(t *testing.T) {
		err := errantly.NewErrantlyError()
		if err.GetSourceLocation() != nil {
			t.Fatal("expected nil source location for default error")
		}
	})

	t.Run("Provider returns a fully wired ErrantlyPackage", func(t *testing.T) {
		pkg := errantly.Provider()
		if pkg == nil {
			t.Fatal("expected non-nil ErrantlyPackage")
		}
		if pkg.NewErrantlyError == nil {
			t.Fatal("NewErrantlyError is nil")
		}
		if pkg.NotImplemented == nil {
			t.Fatal("NotImplemented is nil")
		}
		if pkg.MustFailure == nil {
			t.Fatal("MustFailure is nil")
		}
		if pkg.WithMessage == nil {
			t.Fatal("WithMessage is nil")
		}
		if pkg.InvalidValue == nil {
			t.Fatal("InvalidValue is nil")
		}
		if pkg.MissingValue == nil {
			t.Fatal("MissingValue is nil")
		}
		e := pkg.NewErrantlyError(pkg.WithMessage("via provider"))
		if e.Error() != "via provider" {
			t.Fatalf("expected %q, got %q", "via provider", e.Error())
		}
	})
}

func TestErrantly_CauseChain(t *testing.T) {
	root := errors.New("root cause")
	mid := errantly.NewErrantlyError(
		errantly.WithMessage("middle"),
		errantly.WithCause(root),
	)
	top := errantly.NewErrantlyError(
		errantly.WithMessage("top level"),
		errantly.WithCause(mid),
	)
	if top.Error() != "top level" {
		t.Fatalf("expected top level message")
	}
	if top.GetCause() == nil {
		t.Fatal("expected mid-level cause")
	}
	if top.GetCause().Error() != "middle" {
		t.Fatalf("expected middle message, got %q", top.GetCause().Error())
	}
	midErr, ok := top.GetCause().(errantly.ErrantlyError)
	if !ok {
		t.Fatal("expected GetCause() to return ErrantlyError")
	}
	if midErr.GetCause() == nil {
		t.Fatal("expected root cause")
	}
	if midErr.GetCause().Error() != "root cause" {
		t.Fatalf("expected root cause, got %q", midErr.GetCause().Error())
	}
}

func TestErrantly_GetNotatedCall(t *testing.T) {
	nc := errantly.GetNotatedCall(0)
	if nc == nil {
		t.Fatal("expected non-nil NotatedCall")
	}
	if !strings.Contains(nc.Function, "GetNotatedCall") {
		t.Fatalf("expected function name to contain GetNotatedCall, got %q", nc.Function)
	}
	if nc.FilePath == "" {
		t.Fatal("expected non-empty FilePath")
	}
	if nc.Line <= 0 {
		t.Fatalf("expected positive line number, got %d", nc.Line)
	}
}

func TestMust(t *testing.T) {
	t.Run("returns value on success", func(t *testing.T) {
		result := errantly.Must(42, error(nil))
		if result != 42 {
			t.Fatalf("expected 42, got %v", result)
		}
	})
	t.Run("panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, but none occurred")
			}
		}()
		errantly.Must("boom", errantly.NewErrantlyError(
			errantly.WithMessage("something failed"),
		))
	})
	t.Run("returns string on success", func(t *testing.T) {
		result := errantly.Must("hello", error(nil))
		if result != "hello" {
			t.Fatalf("expected hello, got %v", result)
		}
	})
}

func TestRegisterDeps(t *testing.T) {
	t.Run("sets up journal and log", func(t *testing.T) {
		var journalCalled string
		core := &errantly.ErrantlyDependencies{
			JournalPackage: &errantly.JournalPackage{
				NewJournal: func(name string, options ...errantly.OptionerFn) (errantly.Journalist, error) {
					journalCalled = name
					return &testJournalist{}, nil
				},
			},
		}
		err := errantly.RegisterDeps(core)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if journalCalled != "kaltra.io.errantly" {
			t.Fatalf("expected journal name %q, got %q", "kaltra.io.errantly", journalCalled)
		}
	})
}

type testJournalist struct{}

func (j *testJournalist) GetKind() types.JournalHandlerKind           { return "test" }
func (j *testJournalist) GetName() types.JournalName                  { return "test-journal" }
func (j *testJournalist) SetHandler(h types.JournalHandlerKind) error { return nil }
func (j *testJournalist) Critical(format string, args ...any)         {}
func (j *testJournalist) Error_(format string, args ...any)           {}
func (j *testJournalist) Warn(format string, args ...any)             {}
func (j *testJournalist) Notice(format string, args ...any)           {}
func (j *testJournalist) Info(format string, args ...any)             {}
func (j *testJournalist) Debug(format string, args ...any)            {}
func (j *testJournalist) Trace(format string, args ...any)            {}
func (j *testJournalist) Verbose(format string, args ...any)          {}
func (j *testJournalist) Prolix(format string, args ...any)           {}
func (j *testJournalist) Log(format string, args ...any)              {}

func TestNewErrantlyError_PanicsOnFailingOption(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from failing option, but none occurred")
		}
	}()
	failingOption := func(targets ...any) error {
		return errantly.NewErrantlyError(errantly.WithMessage("option failed"))
	}
	errantly.NewErrantlyError(failingOption)
}
