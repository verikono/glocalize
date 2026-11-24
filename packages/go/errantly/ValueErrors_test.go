package errantly_test

import (
	"testing"

	"github.com/verikono/glocalize/packages/go/errantly"
)

func TestValueErrors(t *testing.T) {

	t.Run("MissingValue", func(t *testing.T) {

		t.Run("Instantiation :: with variable name only", func(t *testing.T) {
			err := errantly.MissingValue("configPath")
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			expected := "variable configPath has no value."
			if err.Error() != expected {
				t.Fatalf("expected %q, got %q", expected, err.Error())
			}
			if err.GetSourceLocation() == nil {
				t.Fatal("expected auto-captured source location")
			}
		})

		t.Run("Instantiation :: with custom message override", func(t *testing.T) {
			err := errantly.MissingValue("apiKey",
				errantly.WithMessage("the API key must be set in the environment"),
			)
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if err.Error() != "the API key must be set in the environment" {
				t.Fatalf("expected overridden message, got %q", err.Error())
			}
		})

		t.Run("Instantiation :: with cause", func(t *testing.T) {
			cause := errantly.NewErrantlyError(
				errantly.WithMessage("env file not found"),
			)
			err := errantly.MissingValue("DATABASE_URL",
				errantly.WithCause(cause),
			)
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if err.Error() != "variable DATABASE_URL has no value." {
				t.Fatalf("expected default message, got %q", err.Error())
			}
			caught := err.GetCause()
			if caught == nil {
				t.Fatal("expected cause via GetCause()")
			}
		})

		t.Run("source location is auto-captured", func(t *testing.T) {
			err := errantly.MissingValue("testVar")
			loc := err.GetSourceLocation()
			if loc == nil {
				t.Fatal("expected auto-captured source location")
			}
			if loc.FilePath == "" {
				t.Fatal("expected non-empty file path")
			}
			if loc.Line <= 0 {
				t.Fatalf("expected positive line number, got %d", loc.Line)
			}
			if loc.Function == "" {
				t.Fatal("expected non-empty function name")
			}
		})
	})

	t.Run("InvalidValue", func(t *testing.T) {

		t.Run("Instantiation :: with variable name only", func(t *testing.T) {
			err := errantly.InvalidValue("port")
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			expected := "variable port is invalid"
			if err.Error() != expected {
				t.Fatalf("expected %q, got %q", expected, err.Error())
			}
			if err.GetSourceLocation() == nil {
				t.Fatal("expected auto-captured source location")
			}
		})

		t.Run("Instantiation :: with custom message override", func(t *testing.T) {
			err := errantly.InvalidValue("timeout",
				errantly.WithMessage("timeout must be between 1 and 60 seconds"),
			)
			if err.Error() != "timeout must be between 1 and 60 seconds" {
				t.Fatalf("expected overridden message, got %q", err.Error())
			}
		})

		t.Run("Instantiation :: with cause", func(t *testing.T) {
			cause := errantly.NewErrantlyError(
				errantly.WithMessage("parsing failed"),
			)
			err := errantly.InvalidValue("duration",
				errantly.WithCause(cause),
			)
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			caught := err.GetCause()
			if caught == nil {
				t.Fatal("expected cause via GetCause()")
			}
		})

		t.Run("source location is auto-captured", func(t *testing.T) {
			err := errantly.InvalidValue("testVar")
			loc := err.GetSourceLocation()
			if loc == nil {
				t.Fatal("expected auto-captured source location")
			}
			if loc.Line <= 0 {
				t.Fatalf("expected positive line number, got %d", loc.Line)
			}
		})
	})
}
