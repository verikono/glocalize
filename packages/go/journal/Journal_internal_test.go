package journal

import (
	"bytes"
	"strings"
	"testing"
)

// Journal_internal_test.go exists because verifying the presence of the
// console handler requires access to the unexported handlers slice and
// console field on Journal.

func TestJournalConsoleHandlerInternal(t *testing.T) {
	t.Run("console handler is added to handlers slice by default", func(t *testing.T) {
		j, err := NewJournal("internal.test.console-default")
		if err != nil {
			t.Fatal(err)
		}
		if len(j.handlers) == 0 {
			t.Fatal("expected at least one handler, got none")
		}
		foundConsole := false
		for _, h := range j.handlers {
			if h.GetKind() == ConsoleJournalHandlerKind {
				foundConsole = true
				break
			}
		}
		if !foundConsole {
			t.Fatal("expected console handler in handlers slice")
		}
	})

	t.Run("console field is true by default", func(t *testing.T) {
		j, err := NewJournal("internal.test.console-true")
		if err != nil {
			t.Fatal(err)
		}
		if !j.console {
			t.Fatal("expected console to be true by default")
		}
	})

	t.Run("WithoutConsole sets console to false and removes handlers", func(t *testing.T) {
		j, err := NewJournal("internal.test.console-false", WithoutConsole())
		if err != nil {
			t.Fatal(err)
		}
		if j.console {
			t.Fatal("expected console to be false with WithoutConsole")
		}
	})

	t.Run("logging dispatches formatted entries to the console handler", func(t *testing.T) {
		var out bytes.Buffer
		j, err := NewJournal("internal.test.handler-write", WithWriter(&out))
		if err != nil {
			t.Fatal(err)
		}
		j.Info("hello %s", "world")
		if !strings.Contains(out.String(), "hello world") {
			t.Fatalf("expected formatted message in output, got %q", out.String())
		}
	})
}
