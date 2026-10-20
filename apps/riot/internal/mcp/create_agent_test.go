package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// agentSkillPriority returns the priority an agent is bound to a skill at, or
// "" when the pair is not bound.
func agentSkillPriority(t *testing.T, server *RiotMcp, agentName, skillName string) orm.PriorityType {
	t.Helper()
	ctx := context.Background()
	agent, err := server.agentByName(ctx, agentName)
	if err != nil {
		t.Fatalf("find agent %s: %v", agentName, err)
	}
	if agent == nil {
		t.Fatalf("agent %s not found", agentName)
	}
	skill, err := server.skillByName(ctx, skillName)
	if err != nil {
		t.Fatalf("find skill %s: %v", skillName, err)
	}
	if skill == nil {
		t.Fatalf("skill %s not found", skillName)
	}
	row, err := server.db.AgentSkill.FindFirst(
		eq[orm.AgentSkill, string]("agentId", agent.Id),
		eq[orm.AgentSkill, string]("skillId", skill.Id),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("find agent skill: %v", err)
	}
	if row == nil {
		return ""
	}
	return row.Priority
}

// TestCreateAgentBindsParticipatingAgent proves every new agent is bound to
// participating-agent at ESSENTIAL priority without the caller asking.
func TestCreateAgentBindsParticipatingAgent(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedSkill(t, server, defaultAgentSkillName)

	result := call(t, server.onCreateAgent, map[string]any{
		"name":        "fresh",
		"description": "a fresh agent",
	})
	if result.IsError {
		t.Fatalf("create agent: %s", resultText(t, result))
	}

	if got := agentSkillPriority(t, server, "fresh", defaultAgentSkillName); got != orm.Priority_ESSENTIAL {
		t.Fatalf("participating-agent priority = %q, want ESSENTIAL", got)
	}

	// the binding must show up in the check-in document too.
	doc := resultText(t, call(t, server.onAgentCheckin, map[string]any{"agent": "fresh", "project": "all"}))
	if !strings.Contains(doc, "ESSENTIAL") || !strings.Contains(doc, defaultAgentSkillName) {
		t.Fatalf("check-in document does not list participating-agent as ESSENTIAL:\n%s", doc)
	}
}

// TestCreateAgentReportsMissingDefault proves a missing default skill does not
// block the create: the agent exists and the result names what was skipped.
func TestCreateAgentReportsMissingDefault(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")

	result := call(t, server.onCreateAgent, map[string]any{
		"name":        "orphan",
		"description": "an agent created before the default skill exists",
	})
	if result.IsError {
		t.Fatalf("create agent: %s", resultText(t, result))
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(resultText(t, result)), &payload); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if payload.Name != "orphan" {
		t.Fatalf("agent name = %q, want orphan", payload.Name)
	}
	if len(result.Content) < 2 {
		t.Fatalf("result has %d content blocks, want the record plus a warning", len(result.Content))
	}
	warning, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("second content block is not text: %T", result.Content[1])
	}
	if !strings.Contains(warning.Text, defaultAgentSkillName) {
		t.Fatalf("warning = %q, want it to name %s", warning.Text, defaultAgentSkillName)
	}
}

// TestAgentCheckinBindsParticipatingAgent proves an agent auto-created by
// check-in gets the same default binding as one created by create-agent.
func TestAgentCheckinBindsParticipatingAgent(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedSkill(t, server, defaultAgentSkillName)

	result := call(t, server.onAgentCheckin, map[string]any{"agent": "checkin-only", "project": "all"})
	if result.IsError {
		t.Fatalf("check-in: %s", resultText(t, result))
	}

	if got := agentSkillPriority(t, server, "checkin-only", defaultAgentSkillName); got != orm.Priority_ESSENTIAL {
		t.Fatalf("participating-agent priority = %q, want ESSENTIAL", got)
	}
}
