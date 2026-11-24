package errantly_test

import (
	"testing"

	"github.com/verikono/glocalize/packages/go/errantly"
)

func TestDeveloperErrors(t *testing.T) {

	t.Run("NotImplemented", func(t *testing.T) {

		t.Run("Instantiation :: with no options", func(t *testing.T) {
			err := errantly.NotImplemented()
			if err == nil {
				t.Fatal("expected non-nil NotImplemented error")
			}
			if err.Error() != "An error has occured" {
				t.Fatalf("expected default message, got %q", err.Error())
			}
			if err.GetSourceLocation() == nil {
				t.Fatal("expected auto-captured source location")
			}
		})

		t.Run("Instantiation :: with custom message", func(t *testing.T) {
			err := errantly.NotImplemented(
				errantly.WithMessage("feature X is not yet implemented"),
			)
			if err.Error() != "feature X is not yet implemented" {
				t.Fatalf("expected custom message, got %q", err.Error())
			}
		})

		t.Run("Instantiation :: with cause", func(t *testing.T) {
			cause := errantly.NewErrantlyError(
				errantly.WithMessage("stub called"),
			)
			err := errantly.NotImplemented(
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

		t.Run("source location is auto-captured at call site", func(t *testing.T) {
			err := errantly.NotImplemented()
			loc := err.GetSourceLocation()
			if loc == nil {
				t.Fatal("expected auto-captured source location")
			}
			if loc.FilePath == "" {
				t.Fatal("expected non-empty file path")
			}
			if loc.Line <= 0 {
				t.Fatalf("expected positive line, got %d", loc.Line)
			}
		})
	})

	t.Run("MustFailure", func(t *testing.T) {

		t.Run("Instantiation :: with no options", func(t *testing.T) {
			err := errantly.MustFailure()
			if err == nil {
				t.Fatal("expected non-nil MustFailure error")
			}
			if err.Error() != "An error has occured" {
				t.Fatalf("expected default message, got %q", err.Error())
			}
			if err.GetSourceLocation() == nil {
				t.Fatal("expected auto-captured source location")
			}
		})

		t.Run("Instantiation :: with custom message", func(t *testing.T) {
			err := errantly.MustFailure(
				errantly.WithMessage("assertion failed: expected non-nil value"),
			)
			if err.Error() != "assertion failed: expected non-nil value" {
				t.Fatalf("expected custom message, got %q", err.Error())
			}
		})

		t.Run("Instantiation :: with cause", func(t *testing.T) {
			cause := errantly.NewErrantlyError(
				errantly.WithMessage("nil pointer dereference"),
			)
			err := errantly.MustFailure(
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

		t.Run("source location is auto-captured at call site", func(t *testing.T) {
			err := errantly.MustFailure()
			loc := err.GetSourceLocation()
			if loc == nil {
				t.Fatal("expected auto-captured source location")
			}
			if loc.FilePath == "" {
				t.Fatal("expected non-empty file path")
			}
			if loc.Line <= 0 {
				t.Fatalf("expected positive line, got %d", loc.Line)
			}
		})
	})
}
