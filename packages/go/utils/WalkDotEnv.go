package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

// WalkDotEnv loads every `.env` file found between the repository root and the
// directory of the calling source file, so a package inherits the environment
// declared by its ancestors. Files are applied root-first: a `.env` closer to
// the caller wins over one further up the tree. Variables already present in
// the process environment are never overwritten, and a missing `.env` is not an
// error. When no repository root marker (`go.work` or `.git`) is found above
// the caller, only the caller's own directory is considered.
func WalkDotEnv() error {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("walk dotenv: resolve caller")
	}
	callerDir := filepath.Dir(filename)

	dirs, err := pathBetween(findRepoRoot(callerDir), callerDir)
	if err != nil {
		return fmt.Errorf("walk dotenv: %w", err)
	}

	merged := map[string]string{}
	for _, dir := range dirs {
		values, err := readDotEnv(filepath.Join(dir, ".env"))
		if err != nil {
			return fmt.Errorf("walk dotenv: %w", err)
		}
		for key, value := range values {
			merged[key] = value
		}
	}
	return applyEnv(merged)
}

// applyEnv sets every key that the process environment does not already define.
func applyEnv(values map[string]string) error {
	for key, value := range values {
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
	}
	return nil
}

// findRepoRoot walks up from dir until it finds a `go.work` or `.git` marker,
// returning dir itself when no marker exists.
func findRepoRoot(dir string) string {
	for current := dir; ; {
		for _, marker := range []string{"go.work", ".git"} {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			return dir
		}
		current = parent
	}
}

// pathBetween returns the directories from root down to dir, inclusive. When
// dir is not nested under root, only dir is returned.
func pathBetween(root, dir string) ([]string, error) {
	relative, err := filepath.Rel(root, dir)
	if err != nil {
		return nil, fmt.Errorf("relate %s to %s: %w", dir, root, err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return []string{dir}, nil
	}

	dirs := []string{dir}
	for current := dir; current != root; {
		current = filepath.Dir(current)
		dirs = append(dirs, current)
	}
	for left, right := 0, len(dirs)-1; left < right; left, right = left+1, right-1 {
		dirs[left], dirs[right] = dirs[right], dirs[left]
	}
	return dirs, nil
}

// readDotEnv parses one `.env` file, treating a missing file as empty.
func readDotEnv(path string) (map[string]string, error) {
	values, err := godotenv.Read(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return values, nil
}
