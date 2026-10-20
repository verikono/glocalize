package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTailOffset(t *testing.T) {
	tests := []struct {
		name    string
		content string
		n       int
		want    string
	}{
		{"last two of four", "a\nb\nc\nd\n", 2, "c\nd\n"},
		{"last one of four", "a\nb\nc\nd\n", 1, "d\n"},
		{"more than present", "a\nb\n", 10, "a\nb\n"},
		{"no trailing newline", "a\nb\nc", 2, "b\nc"},
		{"single line", "only\n", 1, "only\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "riot.log")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}

			offset, err := tailOffset(path, tt.n)
			if err != nil {
				t.Fatalf("tailOffset() error = %v", err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			if got := string(data[offset:]); got != tt.want {
				t.Fatalf("tail = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultLogPath(t *testing.T) {
	got := DefaultLogPath()
	if filepath.Base(got) != "riot.log" {
		t.Fatalf("DefaultLogPath() = %q, want a path ending in riot.log", got)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("DefaultLogPath() = %q, want an absolute path", got)
	}
}
