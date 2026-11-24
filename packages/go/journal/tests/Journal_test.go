package journal_tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/verikono/glocalize/packages/go/journal"
)

func TestJournalInstantiation(t *testing.T) {
	t.Run("creates a journal with a valid name", func(t *testing.T) {
		name := journal.JournalName(uniqueName(t, "test"))
		j, err := journal.NewJournal(name)
		if err != nil {
			t.Fatalf("failed to create journal: %v", err)
		}
		if j.Name != name {
			t.Fatalf("expected %q, got %q", name, j.Name)
		}
	})

	t.Run("rejects duplicate journal names", func(t *testing.T) {
		name := journal.JournalName(uniqueName(t, "dup"))
		_, err := journal.NewJournal(name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = journal.NewJournal(name)
		if err == nil {
			t.Fatal("expected error for duplicate journal name")
		}
	})

	t.Run("rejects invalid journal names", func(t *testing.T) {
		invalidNames := []string{
			"Test",         // uppercase
			"test journal", // space
			"test_journal", // underscore
			"test@journal", // special char
		}
		for _, name := range invalidNames {
			_, err := journal.NewJournal(journal.JournalName(name))
			if err == nil {
				t.Fatalf("expected error for invalid name %q", name)
			}
		}
	})
}

func TestSeverityMethods(t *testing.T) {
	t.Run("all severity methods execute without panic", func(t *testing.T) {
		name := journal.JournalName(uniqueName(t, "sev"))
		j, err := journal.NewJournal(name)
		if err != nil {
			t.Fatal(err)
		}

		j.Critical("critical message")
		j.Error_("error message")
		j.Warn("warn message")
		j.Notice("notice message")
		j.Info("info message")
		j.Log("log message")
		j.Debug("debug message")
		j.Trace("trace message")
		j.Verbose("verbose message")
		j.Prolix("prolix message")
	})

	t.Run("severity methods format their arguments", func(t *testing.T) {
		var out bytes.Buffer
		name := journal.JournalName(uniqueName(t, "fmt"))
		j, err := journal.NewJournal(name, journal.WithWriter(&out))
		if err != nil {
			t.Fatal(err)
		}
		j.Info("connected to %s (gateway: %s)", "wlan0", "192.168.1.1")
		j.Warn("something happened: %v", "boom")
		if !strings.Contains(out.String(), "connected to wlan0 (gateway: 192.168.1.1)") {
			t.Fatalf("expected formatted info message, got %q", out.String())
		}
		if !strings.Contains(out.String(), "something happened: boom") {
			t.Fatalf("expected formatted warn message, got %q", out.String())
		}
	})
}

func TestChildren(t *testing.T) {
	t.Run("parent sees direct children and grandchildren", func(t *testing.T) {
		prefix := uniqueName(t, "parent")
		parentName := journal.JournalName(prefix)
		// JournalName.Validate allows only [a-z . -], so avoid digits in segments.
		child1Name := journal.JournalName(prefix + ".first")
		child2Name := journal.JournalName(prefix + ".second")
		grandchildName := journal.JournalName(prefix + ".second.grandchild")
		otherName := journal.JournalName(prefix + "x.other")

		parent, err := journal.NewJournal(parentName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = journal.NewJournal(child1Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = journal.NewJournal(child2Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = journal.NewJournal(grandchildName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = journal.NewJournal(otherName)
		if err != nil {
			t.Fatal(err)
		}

		children := parent.GetChildren()

		names := childNames(children)
		if !contains(names, string(child1Name)) {
			t.Errorf("expected child %q in children: %v", child1Name, names)
		}
		if !contains(names, string(child2Name)) {
			t.Errorf("expected child %q in children: %v", child2Name, names)
		}
		if !contains(names, string(grandchildName)) {
			t.Errorf("expected grandchild %q in children: %v", grandchildName, names)
		}
		if contains(names, string(otherName)) {
			t.Errorf("did NOT expect %q in children: %v", otherName, names)
		}
	})

	t.Run("leaf child has no children", func(t *testing.T) {
		prefix := uniqueName(t, "leaf")
		parentName := journal.JournalName(prefix)
		childName := journal.JournalName(prefix + ".child")

		_, err := journal.NewJournal(parentName)
		if err != nil {
			t.Fatal(err)
		}
		child, err := journal.NewJournal(childName)
		if err != nil {
			t.Fatal(err)
		}

		children := child.GetChildren()
		if len(children) != 0 {
			t.Fatalf("expected 0 children for leaf, got %d", len(children))
		}
	})

	t.Run("parent with nested child sees grandchild", func(t *testing.T) {
		prefix := uniqueName(t, "nested")
		parentName := journal.JournalName(prefix)
		// JournalName.Validate allows only [a-z . -], so avoid digits in segments.
		child2Name := journal.JournalName(prefix + ".second")
		grandchildName := journal.JournalName(prefix + ".second.grandchild")

		_, err := journal.NewJournal(parentName)
		if err != nil {
			t.Fatal(err)
		}
		child2, err := journal.NewJournal(child2Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = journal.NewJournal(grandchildName)
		if err != nil {
			t.Fatal(err)
		}

		children := child2.GetChildren()
		if len(children) != 1 {
			t.Fatalf("expected 1 child for child2, got %d: %v", len(children), childNames(children))
		}
		if string(children[0].Name) != string(grandchildName) {
			t.Fatalf("expected grandchild %q, got %q", grandchildName, children[0].Name)
		}
	})
}

func uniqueName(t *testing.T, suffix string) string {
	t.Helper()
	// JournalName.Validate allows only [a-z . -], so strip any other
	// characters (Go subtest names contain underscores/spaces).
	safe := strings.ToLower(t.Name())
	safe = strings.ReplaceAll(safe, "/", "-")
	safe = strings.ReplaceAll(safe, " ", "-")
	safe = strings.ReplaceAll(safe, "_", "-")
	return safe + "." + suffix
}

func childNames(children []*journal.Journal) []string {
	names := make([]string, len(children))
	for i, c := range children {
		names[i] = string(c.Name)
	}
	return names
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
