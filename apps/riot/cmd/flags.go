package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// runtimeFlags records how the running daemon was started so `riot restart`
// can replay the same flags.
type runtimeFlags struct {
	Port   int    `json:"port"`
	Hot    bool   `json:"hot"`
	DevDir string `json:"dev_dir,omitempty"`
}

// flagsPath returns the state file beneath the user config directory.
func flagsPath() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "riot", "flags")
	}
	return filepath.Join(os.TempDir(), "riot", "flags")
}

func marshalFlags(flags runtimeFlags) ([]byte, error) {
	return json.Marshal(flags)
}

func unmarshalFlags(data []byte) (runtimeFlags, error) {
	var flags runtimeFlags
	if err := json.Unmarshal(data, &flags); err != nil {
		return runtimeFlags{}, err
	}
	return flags, nil
}

// saveRuntimeFlags persists the start flags for `riot restart`.
func saveRuntimeFlags(flags runtimeFlags) error {
	data, err := marshalFlags(flags)
	if err != nil {
		return fmt.Errorf("riot: marshal flags: %w", err)
	}
	path := flagsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("riot: flags dir: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("riot: write flags: %w", err)
	}
	return nil
}

// loadRuntimeFlags reads the persisted start flags.
func loadRuntimeFlags() (runtimeFlags, error) {
	data, err := os.ReadFile(flagsPath())
	if err != nil {
		return runtimeFlags{}, fmt.Errorf("riot: read flags: %w", err)
	}
	flags, err := unmarshalFlags(data)
	if err != nil {
		return runtimeFlags{}, fmt.Errorf("riot: parse flags: %w", err)
	}
	return flags, nil
}
