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

// agentCheckinLine is the sole body every mirrored agent file carries. It is
// deliberately blunt: check in first, and if that cannot be done, stop dead
// rather than working without an identity.
const agentCheckinLine = "CHECKIN TO RIOT-MCP (tool: `agent-checkin`) AS THE AGENT NAMED IN THIS FILE'S `name` FRONTMATTER FIELD. DO THIS AS YOUR FIRST ACTION, BEFORE ANY OTHER TOOL CALL. IF RIOT-MCP IS NOT AVAILABLE, OR THE CHECK-IN FAILS FOR ANY REASON WHATSOEVER, TELL THE USER THAT RIOT-MCP IS UNAVAILABLE OR THAT CHECK-IN FAILED. THEN STOP: DO NOT PROCEED, DO NOT CALL ANY OTHER TOOL, AND WAIT FOR THE USER'S INSTRUCTIONS."

// agentEssentialGate is the second half of the check-in obligation. Checking in
// without loading the ESSENTIAL skills is a check-in that produced no context,
// so the gate is stated in the same breath as the check-in itself.
const agentEssentialGate = "THEN, BEFORE ANY OTHER ACTION: fetch every ESSENTIAL skill the check-in document lists with `get-skill` — every one, by name, in this session. The check-in is not complete until they are all loaded. ESSENTIAL is a precondition of the session, not a menu: it is not yours to judge whether a skill applies, and a skill you skip because it \"does not seem relevant\" is the exact failure that tier exists to prevent. A load is a tool call, not a recollection — if the transcript holds no `get-skill` call for a skill, that skill is not loaded, because skills are edited and your memory of one is stale by definition. Skills marked _(required by …)_ are dependencies and are ESSENTIAL by inheritance. If a `get-skill` call fails, name the skill and the failure, then stop and wait. Your first message must list every skill you loaded, by name: a list that omits an ESSENTIAL skill is a false claim of readiness."

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
	b.WriteString(agentCheckinLine)
	b.WriteString("\n\n")
	b.WriteString(agentEssentialGate)
	b.WriteString("\n")
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
