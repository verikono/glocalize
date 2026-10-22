package mcp

import (
	"context"
	"testing"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// seedSkillRecord creates a skill and returns the stored record.
func seedSkillRecord(t *testing.T, server *RiotMcp, name string) *orm.Skill {
	t.Helper()
	seedSkill(t, server, name)
	skill, err := server.skillByName(context.Background(), name)
	if err != nil || skill == nil {
		t.Fatalf("%s not stored: %v", name, err)
	}
	return skill
}

// newDependencyMcp opens a store with the "all" project seeded, which
// add-skill requires.
func newDependencyMcp(t *testing.T) *RiotMcp {
	t.Helper()
	server := newTestMcp(t)
	seedProject(t, server, "all")
	return server
}

// linkSkill makes dependent depend on dependency with the given recursion.
func linkSkill(t *testing.T, server *RiotMcp, dependent, dependency string, recursive bool) {
	t.Helper()
	result := call(t, server.onMutateSkill, map[string]any{
		"name": dependent,
		"mutation": map[string]any{
			"op": "update",
			"dependencies": map[string]any{
				"add": []any{map[string]any{"key": dependency, "recursive": recursive}},
			},
		},
	})
	if result.IsError {
		t.Fatalf("link %s -> %s: %s", dependent, dependency, resultText(t, result))
	}
}

// boundNames returns the names of the skills bound to an agent.
func boundNames(t *testing.T, server *RiotMcp, agentID string) map[string]bool {
	t.Helper()
	skills, err := server.boundSkills(context.Background(), agentID)
	if err != nil {
		t.Fatalf("boundSkills: %v", err)
	}
	out := make(map[string]bool, len(skills))
	for _, s := range skills {
		out[s.Name] = true
	}
	return out
}

// bindAgent creates an agent with the named skill bound directly to it.
func bindAgent(t *testing.T, server *RiotMcp, name string, skillIDs ...string) string {
	t.Helper()
	ctx := context.Background()
	agent, err := server.db.Agent.Create().SetName(name).SetDescription(name).SetAllowedTools("[]").SetMetadata("{}").Exec(ctx)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	for _, skillID := range skillIDs {
		if _, err := server.db.AgentSkill.Create().SetAgentId(agent.Id).SetSkillId(skillID).SetPriority(orm.Priority_NORMAL).Exec(ctx); err != nil {
			t.Fatalf("bind skill: %v", err)
		}
	}
	return agent.Id
}

// TestBoundSkillsRecursiveEdgePullsWholeChain proves a recursive edge (the
// default) hands the agent the dependency and everything that dependency needs.
func TestBoundSkillsRecursiveEdgePullsWholeChain(t *testing.T) {
	server := newDependencyMcp(t)
	root := seedSkillRecord(t, server, "root")
	seedSkillRecord(t, server, "middle")
	seedSkillRecord(t, server, "leaf")

	linkSkill(t, server, "root", "middle", true)
	linkSkill(t, server, "middle", "leaf", true)

	agentID := bindAgent(t, server, "recursive-agent", root.Id)
	got := boundNames(t, server, agentID)
	for _, name := range []string{"root", "middle", "leaf"} {
		if !got[name] {
			t.Fatalf("%q missing from bound skills %v", name, got)
		}
	}
}

// TestBoundSkillsNonRecursiveEdgeStops proves a non-recursive edge loads the
// dependency and stops: the dependency's own dependency is not loaded.
func TestBoundSkillsNonRecursiveEdgeStops(t *testing.T) {
	server := newDependencyMcp(t)
	root := seedSkillRecord(t, server, "root")
	seedSkillRecord(t, server, "middle")
	seedSkillRecord(t, server, "leaf")

	linkSkill(t, server, "root", "middle", false)
	linkSkill(t, server, "middle", "leaf", true)

	agentID := bindAgent(t, server, "stopping-agent", root.Id)
	got := boundNames(t, server, agentID)
	if !got["root"] || !got["middle"] {
		t.Fatalf("root and middle must be bound, got %v", got)
	}
	if got["leaf"] {
		t.Fatalf("leaf must not be bound through a non-recursive edge, got %v", got)
	}
}

// TestLinkSkillDefaultIsRecursive proves an add without an explicit recursive
// flag stores true, so an existing caller keeps the old whole-chain behaviour.
func TestLinkSkillDefaultIsRecursive(t *testing.T) {
	server := newDependencyMcp(t)
	seedSkillRecord(t, server, "alpha")
	beta := seedSkillRecord(t, server, "beta")

	result := call(t, server.onMutateSkill, map[string]any{
		"name": "alpha",
		"mutation": map[string]any{
			"op":           "update",
			"dependencies": map[string]any{"add": []any{"beta"}},
		},
	})
	if result.IsError {
		t.Fatalf("link alpha -> beta: %s", resultText(t, result))
	}

	rows, err := server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependencyId", beta.Id)).Exec(context.Background())
	if err != nil {
		t.Fatalf("read dependency: %v", err)
	}
	if len(rows) != 1 || !rows[0].Recursive {
		t.Fatalf("edges = %+v, want one recursive edge", rows)
	}
}

// TestLinkSkillRestatesRecursive proves re-adding an existing edge with a
// different recursive value patches it rather than being ignored.
func TestLinkSkillRestatesRecursive(t *testing.T) {
	server := newDependencyMcp(t)
	alpha := seedSkillRecord(t, server, "alpha")
	seedSkillRecord(t, server, "beta")

	linkSkill(t, server, "alpha", "beta", true)
	linkSkill(t, server, "alpha", "beta", false)

	rows, err := server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependentId", alpha.Id)).Exec(context.Background())
	if err != nil {
		t.Fatalf("read dependencies: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("edges = %d, want 1", len(rows))
	}
	if rows[0].Recursive {
		t.Fatalf("recursive = true, want false after the restatement")
	}
}