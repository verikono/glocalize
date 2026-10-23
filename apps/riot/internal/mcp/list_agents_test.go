package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// seedAgent creates an agent with the given metadata and binds the named
// skills to it.
func seedAgent(t *testing.T, server *RiotMcp, name, metadata string, skills ...string) {
	t.Helper()
	result := call(t, server.onCreateAgent, map[string]any{
		"name":        name,
		"description": name + " agent",
		"metadata":    metadata,
	})
	if result.IsError {
		t.Fatalf("create agent %s: %s", name, resultText(t, result))
	}
	for _, skill := range skills {
		result := call(t, server.onAssociateAgentSkill, map[string]any{"agent": name, "skill": skill})
		if result.IsError {
			t.Fatalf("bind %s to %s: %s", skill, name, resultText(t, result))
		}
	}
}

// seedSkill creates a skill with the minimum required fields.
func seedSkill(t *testing.T, server *RiotMcp, name string) {
	t.Helper()
	result := call(t, server.onAddSkill, map[string]any{
		"project":     "all",
		"name":        name,
		"description": name + " skill",
		"content":     "body",
	})
	if result.IsError {
		t.Fatalf("create skill %s: %s", name, resultText(t, result))
	}
}

// seedProject creates a project and binds the named agents to it.
func seedProject(t *testing.T, server *RiotMcp, name string, agents ...string) {
	t.Helper()
	ctx := context.Background()
	if _, err := server.db.Project.Create().
		SetName(name).
		SetLocation("N/A").
		SetSkillsPath("skills").
		SetAgentsPath("agents").
		Exec(ctx); err != nil {
		t.Fatalf("create project %s: %v", name, err)
	}
	for _, agent := range agents {
		result := call(t, server.onAssociateAgentProject, map[string]any{"agent": agent, "project": name})
		if result.IsError {
			t.Fatalf("bind %s to %s: %s", agent, name, resultText(t, result))
		}
	}
}

// decodeAgents unmarshals a list-agents result into the reported shape.
func decodeAgents(t *testing.T, result *mcp.CallToolResult) []agentSummary {
	t.Helper()
	var out []agentSummary
	if err := json.Unmarshal([]byte(resultText(t, result)), &out); err != nil {
		t.Fatalf("decode agents: %v", err)
	}
	return out
}

// TestListAgentsReportsShape proves every agent is listed with its id,
// description, decoded metadata and bound skill names.
func TestListAgentsReportsShape(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedSkill(t, server, "alpha")
	seedSkill(t, server, "beta")
	seedAgent(t, server, "zeta", `{"tier":"high"}`, "alpha", "beta")
	seedAgent(t, server, "alpha", `{}`)

	agents := decodeAgents(t, call(t, server.onListAgents, map[string]any{}))
	if len(agents) != 2 {
		t.Fatalf("agents = %d, want 2", len(agents))
	}
	// ordered by name.
	if agents[0].Name != "alpha" || agents[1].Name != "zeta" {
		t.Fatalf("order = %q, %q; want alpha, zeta", agents[0].Name, agents[1].Name)
	}
	if agents[0].Id == "" {
		t.Fatalf("alpha has no id")
	}
	if agents[0].Description != "alpha agent" {
		t.Fatalf("description = %q, want %q", agents[0].Description, "alpha agent")
	}
	if string(agents[0].Metadata) != "{}" {
		t.Fatalf("alpha metadata = %s, want {}", agents[0].Metadata)
	}
	if string(agents[1].Metadata) != `{"tier":"high"}` {
		t.Fatalf("zeta metadata = %s, want the stored object", agents[1].Metadata)
	}
	if len(agents[1].Skills) != 2 {
		t.Fatalf("zeta skills = %d, want 2", len(agents[1].Skills))
	}
	names := map[string]bool{}
	for _, s := range agents[1].Skills {
		names[s.Name] = true
	}
	if !names["alpha"] || !names["beta"] {
		t.Fatalf("zeta skills = %v, want alpha and beta", names)
	}
	if len(agents[0].Skills) != 0 {
		t.Fatalf("alpha skills = %d, want 0", len(agents[0].Skills))
	}
}

// TestListAgentsFiltersByProject proves the project argument narrows the list
// to agents associated with any of the named projects.
func TestListAgentsFiltersByProject(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedProject(t, server, "kaltra")
	seedAgent(t, server, "shared", `{}`)
	seedAgent(t, server, "kaltra-only", `{}`)
	seedAgent(t, server, "unbound", `{}`)

	call(t, server.onAssociateAgentProject, map[string]any{"agent": "shared", "project": "all"})
	call(t, server.onAssociateAgentProject, map[string]any{"agent": "shared", "project": "kaltra"})
	call(t, server.onAssociateAgentProject, map[string]any{"agent": "kaltra-only", "project": "kaltra"})

	agents := decodeAgents(t, call(t, server.onListAgents, map[string]any{
		"project": []any{"kaltra"},
	}))
	if len(agents) != 2 {
		t.Fatalf("kaltra agents = %d, want 2", len(agents))
	}
	for _, a := range agents {
		if a.Name == "unbound" {
			t.Fatalf("unbound agent leaked into the kaltra filter")
		}
	}

	agents = decodeAgents(t, call(t, server.onListAgents, map[string]any{
		"project": []any{"all", "kaltra"},
	}))
	if len(agents) != 2 {
		t.Fatalf("all+kaltra agents = %d, want 2 (shared counted once)", len(agents))
	}

	agents = decodeAgents(t, call(t, server.onListAgents, map[string]any{
		"project": []any{"nope"},
	}))
	if len(agents) != 0 {
		t.Fatalf("unknown project agents = %d, want 0", len(agents))
	}
}

// TestListAgentsDecodesBadMetadata proves a blank or malformed stored metadata
// string is reported as an empty object rather than breaking the response.
func TestListAgentsDecodesBadMetadata(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedAgent(t, server, "broken", `not json`)

	agents := decodeAgents(t, call(t, server.onListAgents, map[string]any{}))
	if len(agents) != 1 {
		t.Fatalf("agents = %d, want 1", len(agents))
	}
	if string(agents[0].Metadata) != "{}" {
		t.Fatalf("metadata = %s, want {}", agents[0].Metadata)
	}
}
