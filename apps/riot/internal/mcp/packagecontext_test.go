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

// riotGuidance returns the guidance file the riot package root carries, so a
// test can assert against whatever file the preference order resolves to
// rather than a hardcoded one.
func riotGuidance(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	file := agentFileIn(root)
	if file == "" {
		t.Fatalf("apps/riot carries no guidance file")
	}
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
	return string(content)
}

// TestLoadPackageContext_RiotAgentFile proves the tool resolves the real
// apps/riot guidance file from three different locations inside the package.
func TestLoadPackageContext_RiotAgentFile(t *testing.T) {
	server := newTestMcp(t)
	root := filepath.Join("..", "..")
	want := riotGuidance(t)

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
			if got := resultText(t, result); got != want {
				t.Fatalf("expected apps/riot guidance, got:\n%s", got)
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

// TestLoadPackageContext_MultiplePaths proves several paths resolving to
// different guidance files are stitched into one content block, each section
// headed by the file it came from.
func TestLoadPackageContext_MultiplePaths(t *testing.T) {
	server := newTestMcp(t)
	root := t.TempDir()

	first := filepath.Join(root, "first")
	writeFile(t, filepath.Join(first, "go.mod"), "module example.com/first\n")
	writeFile(t, filepath.Join(first, ".agents.md"), "first guidance")

	second := filepath.Join(root, "second")
	writeFile(t, filepath.Join(second, "package.json"), `{"name":"second"}`)
	writeFile(t, filepath.Join(second, "agent.md"), "second guidance")

	nested := filepath.Join(first, "internal", "x.go")
	writeFile(t, nested, "package internal\n")

	result := call(t, server.onLoadPackageContext, map[string]any{
		"paths": []any{
			nested,
			second,
		},
	})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}

	got := resultText(t, result)
	want := "# " + filepath.Join(first, ".agents.md") + "\n\nfirst guidance\n\n---\n\n# " +
		filepath.Join(second, "agent.md") + "\n\nsecond guidance"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

// TestLoadPackageContext_DeduplicatesPaths proves two paths resolving to the
// same guidance file contribute that file once.
func TestLoadPackageContext_DeduplicatesPaths(t *testing.T) {
	server := newTestMcp(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/x\n")
	writeFile(t, filepath.Join(root, ".agents.md"), "package guidance")

	nested := filepath.Join(root, "sub", "file.go")
	writeFile(t, nested, "package sub\n")

	result := call(t, server.onLoadPackageContext, map[string]any{
		"paths": []any{root, nested},
	})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != "package guidance" {
		t.Fatalf("got %q, want %q", got, "package guidance")
	}
}

// TestLoadPackageContext_PathAndPaths proves the single `path` argument and the
// `paths` array are combined.
func TestLoadPackageContext_PathAndPaths(t *testing.T) {
	server := newTestMcp(t)
	root := t.TempDir()

	first := filepath.Join(root, "first")
	writeFile(t, filepath.Join(first, "go.mod"), "module example.com/first\n")
	writeFile(t, filepath.Join(first, ".agents.md"), "first guidance")

	second := filepath.Join(root, "second")
	writeFile(t, filepath.Join(second, "go.mod"), "module example.com/second\n")
	writeFile(t, filepath.Join(second, ".agents.md"), "second guidance")

	result := call(t, server.onLoadPackageContext, map[string]any{
		"path":  first,
		"paths": []any{second},
	})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}

	got := resultText(t, result)
	if !strings.Contains(got, "first guidance") || !strings.Contains(got, "second guidance") {
		t.Fatalf("expected both guidance files, got:\n%s", got)
	}
}

// TestGetProjectContext_ByPackageName proves the riot package resolves by its
// directory name and returns its own guidance file.
func TestGetProjectContext_ByPackageName(t *testing.T) {
	server := newTestMcp(t)
	want := riotGuidance(t)

	result := call(t, server.onGetProjectContext, map[string]any{"package-name": "riot"})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != want {
		t.Fatalf("expected apps/riot guidance, got:\n%s", got)
	}
}

// TestGetProjectContext_ByDirectory proves an absolute directory resolves to
// the closest package root above it.
func TestGetProjectContext_ByDirectory(t *testing.T) {
	server := newTestMcp(t)
	want := riotGuidance(t)
	dir, err := filepath.Abs(filepath.Join("..", "orm"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}

	result := call(t, server.onGetProjectContext, map[string]any{"directory": dir})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != want {
		t.Fatalf("expected apps/riot guidance, got:\n%s", got)
	}
}

// TestGetProjectContext_ByFilepath proves an absolute file path resolves to
// the closest package root above its directory.
func TestGetProjectContext_ByFilepath(t *testing.T) {
	server := newTestMcp(t)
	want := riotGuidance(t)
	file, err := filepath.Abs(filepath.Join("..", "orm", "agent.go"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}

	result := call(t, server.onGetProjectContext, map[string]any{"filepath": file})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != want {
		t.Fatalf("expected apps/riot guidance, got:\n%s", got)
	}
}

// TestGetProjectContext_RequiresExactlyOneArgument proves the tool rejects
// both an empty argument set and more than one supplied argument.
func TestGetProjectContext_RequiresExactlyOneArgument(t *testing.T) {
	server := newTestMcp(t)

	if result := call(t, server.onGetProjectContext, map[string]any{}); !result.IsError {
		t.Fatalf("expected an error for no arguments, got %q", resultText(t, result))
	}

	result := call(t, server.onGetProjectContext, map[string]any{
		"package-name": "riot",
		"directory":    "/tmp",
	})
	if !result.IsError {
		t.Fatalf("expected an error for two arguments, got %q", resultText(t, result))
	}
}

// TestGetProjectContext_UnknownPackage proves an unmatched name reports that
// no information is available rather than erroring.
func TestGetProjectContext_UnknownPackage(t *testing.T) {
	server := newTestMcp(t)

	result := call(t, server.onGetProjectContext, map[string]any{"package-name": "no-such-package-anywhere"})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != packageContextUnavailable {
		t.Fatalf("got %q, want %q", got, packageContextUnavailable)
	}
}

// TestGetProjectContext_PackageWithoutGuidance proves a package root that
// carries no guidance file reports that no information is available.
func TestGetProjectContext_PackageWithoutGuidance(t *testing.T) {
	server := newTestMcp(t)
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/x\n")

	result := call(t, server.onGetProjectContext, map[string]any{"directory": dir})
	if result.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, result))
	}
	if got := resultText(t, result); got != packageContextUnavailable {
		t.Fatalf("got %q, want %q", got, packageContextUnavailable)
	}
}

// TestDeclaredPackageName proves each marker dialect yields its declared name.
func TestDeclaredPackageName(t *testing.T) {
	cases := []struct {
		name string
		file string
		body string
		want string
	}{
		{name: "package.json", file: "package.json", body: `{"name":"@scope/pkg"}`, want: "@scope/pkg"},
		{name: "go.mod", file: "go.mod", body: "module example.com/mod\n\ngo 1.27.0\n", want: "example.com/mod"},
		{name: "pyproject.toml", file: "pyproject.toml", body: "[project]\nname = \"pypkg\"\n", want: "pypkg"},
		{name: "cargo.toml", file: "cargo.toml", body: "[package]\nname = \"crate\"\n", want: "crate"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, tc.file), tc.body)
			if got := declaredPackageName(dir); got != tc.want {
				t.Fatalf("declaredPackageName() = %q, want %q", got, tc.want)
			}
		})
	}
}
