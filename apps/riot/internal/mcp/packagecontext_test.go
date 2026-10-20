package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFile writes content to path, creating parent directories.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// TestLoadPackageContext_RiotAgentFile proves the tool resolves the real
// apps/riot/agent.md from two different locations inside the riot package.
func TestLoadPackageContext_RiotAgentFile(t *testing.T) {
	server := newTestMcp(t)
	root := filepath.Join("..", "..")

	cases := []struct {
		name string
		path string
	}{
		{name: "package root", path: root},
		{name: "nested file", path: filepath.Join(root, "internal", "mcp", "handlers.go")},
		{name: "nested directory", path: filepath.Join(root, "internal", "orm")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := call(t, server.onLoadPackageContext, map[string]any{"path": tc.path})
			if result.IsError {
				t.Fatalf("unexpected error: %s", resultText(t, result))
			}
			text := resultText(t, result)
			if text == packageContextUnavailable {
				t.Fatalf("no guidance found for %s", tc.path)
			}
			if !strings.Contains(text, "riot serve") {
				t.Fatalf("expected apps/riot/agent.md contents, got:\n%s", text)
			}
		})
	}
}

// TestLoadPackageContext_DirectAgentFile proves a guidance file in the target
// directory is returned when that directory is also a package root.
func TestLoadPackageContext_DirectAgentFile(t *testing.T) {
	server := newTestMcp(t)
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/x\n")
	writeFile(t, filepath.Join(dir, ".agents.md"), "direct guidance")

	result := call(t, server.onLoadPackageContext, map[string]any{"path": dir})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != "direct guidance" {
		t.Fatalf("got %q, want %q", got, "direct guidance")
	}
}

// TestLoadPackageContext_SkipsMarkerlessStartDir proves a guidance file in the
// starting directory is skipped when that directory holds no package marker,
// and resolution continues upward to the nearest real package root.
func TestLoadPackageContext_SkipsMarkerlessStartDir(t *testing.T) {
	server := newTestMcp(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/x\n")
	writeFile(t, filepath.Join(root, ".agents.md"), "package guidance")
	writeFile(t, filepath.Join(root, "sub", ".agents.md"), "stray guidance")

	result := call(t, server.onLoadPackageContext, map[string]any{"path": filepath.Join(root, "sub")})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != "package guidance" {
		t.Fatalf("got %q, want %q", got, "package guidance")
	}
}

// TestLoadPackageContext_WalksUpToPackageRoot proves a parent directory is
// only accepted when it holds both a guidance file and a package marker.
func TestLoadPackageContext_WalksUpToPackageRoot(t *testing.T) {
	server := newTestMcp(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/x\n")
	writeFile(t, filepath.Join(root, ".agents.md"), "package guidance")
	nested := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	result := call(t, server.onLoadPackageContext, map[string]any{"path": nested})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != "package guidance" {
		t.Fatalf("got %q, want %q", got, "package guidance")
	}
}

// TestLoadPackageContext_IgnoresMarkerlessParent proves a parent holding a
// guidance file but no package marker is skipped.
func TestLoadPackageContext_IgnoresMarkerlessParent(t *testing.T) {
	server := newTestMcp(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agents.md"), "markerless guidance")
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	result := call(t, server.onLoadPackageContext, map[string]any{"path": nested})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != packageContextUnavailable {
		t.Fatalf("got %q, want %q", got, packageContextUnavailable)
	}
}

// TestLoadPackageContext_RequiresPath proves a missing path is rejected.
func TestLoadPackageContext_RequiresPath(t *testing.T) {
	server := newTestMcp(t)
	result := call(t, server.onLoadPackageContext, map[string]any{})
	if !result.IsError {
		t.Fatalf("expected an error result, got %q", resultText(t, result))
	}
}
