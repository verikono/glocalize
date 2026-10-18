package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

// agentFileSuffix is the suffix every mirrored agent file carries.
const agentFileSuffix = ".agent.md"

// agentDir returns the configured agent directory, or "" when mirroring is off.
func agentDir() string {
	return utils.ResolveAgentDir()
}

// agentFilePath returns the path a named agent's file lives at, or "" when
// mirroring is off.
func agentFilePath(name string) string {
	dir := agentDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, name+agentFileSuffix)
}

// writeAgentFile writes (or replaces) the on-disk mirror of an agent record. It
// is a no-op when RIOT_AGENT_DIR is unset.
func writeAgentFile(name, description, allowedTools string) error {
	path := agentFilePath(name)
	if path == "" {
		return nil
	}
	content, err := renderAgentFile(name, description, allowedTools)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("agent file: create dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("agent file: write %s: %w", path, err)
	}
	return nil
}

// removeAgentFile deletes an agent's on-disk mirror. It is a no-op when
// RIOT_AGENT_DIR is unset or the file is already gone.
func removeAgentFile(name string) error {
	path := agentFilePath(name)
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("agent file: remove %s: %w", path, err)
	}
	return nil
}

// renderAgentFile renders the markdown a mirrored agent file carries. The
// allowed-tools list is emitted only when it is non-empty.
func renderAgentFile(name, description, allowedTools string) (string, error) {
	tools, err := parseToolList(allowedTools)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "name: %s\n", name)
	fmt.Fprintf(&b, "description: %s\n", description)
	if len(tools) > 0 {
		b.WriteString("allowed-tools:\n")
		for _, tool := range tools {
			fmt.Fprintf(&b, "  - %s\n", tool)
		}
	}
	b.WriteString("---\n\n")
	b.WriteString("CHECKIN TO RIOT-MCP WITH THIS AGENT NAME. IF RIOT-MCP IS NOT AVAILABLE INFORM THE USER IT IS NOT AVAILABLE.\n")
	return b.String(), nil
}

// parseToolList decodes a JSON-encoded tool array. A non-array or empty value
// yields an empty list, so the allowed-tools block is omitted.
func parseToolList(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var tools []string
	if err := json.Unmarshal([]byte(raw), &tools); err != nil {
		return nil, nil
	}
	return tools, nil
}
