package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func writeGoMod(t *testing.T, dir, module string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(module+"\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
}

func TestAssertRiotModule(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dir := t.TempDir()
		writeGoMod(t, dir, riotModuleDecl)
		if err := assertRiotModule(dir); err != nil {
			t.Fatalf("assertRiotModule() error = %v", err)
		}
	})

	t.Run("wrong module", func(t *testing.T) {
		dir := t.TempDir()
		writeGoMod(t, dir, "module github.com/someone/else")
		err := assertRiotModule(dir)
		if err == nil || !strings.Contains(err.Error(), "wrong directory for Glocalize::Riot") {
			t.Fatalf("assertRiotModule() error = %v, want wrong-directory error", err)
		}
	})

	t.Run("missing go.mod", func(t *testing.T) {
		if err := assertRiotModule(t.TempDir()); err == nil {
			t.Fatal("assertRiotModule() expected error for missing go.mod")
		}
	})
}

func TestResolveDevDir(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv(riotDevDirEnv, "")
		got, err := resolveDevDir()
		if err != nil || got != "" {
			t.Fatalf("resolveDevDir() = %q, %v; want \"\", nil", got, err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		dir := t.TempDir()
		writeGoMod(t, dir, riotModuleDecl)
		t.Setenv(riotDevDirEnv, dir)
		got, err := resolveDevDir()
		if err != nil {
			t.Fatalf("resolveDevDir() error = %v", err)
		}
		if filepath.Clean(got) != filepath.Clean(dir) {
			t.Fatalf("resolveDevDir() = %q, want %q", got, dir)
		}
	})

	t.Run("wrong directory", func(t *testing.T) {
		dir := t.TempDir()
		writeGoMod(t, dir, "module nope/nope")
		t.Setenv(riotDevDirEnv, dir)
		if _, err := resolveDevDir(); err == nil {
			t.Fatal("resolveDevDir() expected error for wrong directory")
		}
	})
}

func TestChangeType(t *testing.T) {
	cases := map[fsnotify.Op]string{
		fsnotify.Write:  "WRITE",
		fsnotify.Create: "CREATE",
		fsnotify.Remove: "REMOVE",
		fsnotify.Rename: "RENAME",
		fsnotify.Chmod:  "CHMOD",
	}
	for op, want := range cases {
		if got := changeType(op); got != want {
			t.Errorf("changeType(%v) = %q, want %q", op, got, want)
		}
	}
	if got := changeType(fsnotify.Op(0)); got != "CHANGE" {
		t.Errorf("changeType(0) = %q, want CHANGE", got)
	}
}

func TestIgnoredPath(t *testing.T) {
	devDir := "/srv/riot"
	cases := []struct {
		path string
		want bool
	}{
		{"/srv/riot/.git/HEAD", true},
		{"/srv/riot/.git/refs/x", true},
		{"/srv/riot/dist/riot", true},
		{"/srv/riot/internal/utils/HotService.go", false},
		{"/srv/riot/main.go", false},
		{"/srv/riot", true},
		{"/other/riot/x.go", false},
	}
	for _, c := range cases {
		if got := ignoredPath(devDir, c.path); got != c.want {
			t.Errorf("ignoredPath(%q, %q) = %v, want %v", devDir, c.path, got, c.want)
		}
	}
}
