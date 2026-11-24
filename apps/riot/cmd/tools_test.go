package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestWriteTools_NamesOnly proves the default listing prints one bare tool
// name per line.
func TestWriteTools_NamesOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := writeTools(&buf, false); err != nil {
		t.Fatalf("writeTools() error = %v", err)
	}

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("no output")
	}
	for _, line := range lines {
		if strings.Contains(line, ":") {
			t.Fatalf("names-only line carries a separator: %q", line)
		}
		if strings.TrimSpace(line) != line {
			t.Fatalf("names-only line is padded: %q", line)
		}
	}
	if lines[0] != "agent-checkin" {
		t.Fatalf("first line = %q, want %q", lines[0], "agent-checkin")
	}
}

// TestWriteTools_ListAlignsSeparator proves every row in the --list output
// puts its separator at the same column.
func TestWriteTools_ListAlignsSeparator(t *testing.T) {
	var buf bytes.Buffer
	if err := writeTools(&buf, true); err != nil {
		t.Fatalf("writeTools() error = %v", err)
	}

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("no output")
	}

	column := -1
	for _, line := range lines {
		index := strings.Index(line, " : ")
		if index < 0 {
			t.Fatalf("line has no separator: %q", line)
		}
		if column < 0 {
			column = index
			continue
		}
		if index != column {
			t.Fatalf("separator at column %d, want %d: %q", index, column, line)
		}
	}
}

// TestWriteTools_ListIncludesDescriptions proves each row carries a
// description after the separator.
func TestWriteTools_ListIncludesDescriptions(t *testing.T) {
	var buf bytes.Buffer
	if err := writeTools(&buf, true); err != nil {
		t.Fatalf("writeTools() error = %v", err)
	}

	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		_, description, found := strings.Cut(line, " : ")
		if !found {
			t.Fatalf("line has no separator: %q", line)
		}
		if strings.TrimSpace(description) == "" {
			t.Fatalf("line has an empty description: %q", line)
		}
	}
}
