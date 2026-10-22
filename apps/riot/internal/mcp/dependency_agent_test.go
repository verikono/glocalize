package mcp

import (
	"context"
	"testing"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// linkSkillForAgent writes an edge scoped to one agent.
func linkSkillForAgent(t *testing.T, server *RiotMcp, dependent, dependency, agent string, recursive bool) {
	t.Helper()
	result := call(t, server.onMutateSkill, map[string]any{
		"name":  dependent,
		"agent": agent,
		"mutation": map[string]any{
			"op": "update",
			"dependencies": map[string]any{
				"add": []any{map[string]any{"key": dependency, "recursive": recursive}},
			},
		},
	})
	if result.IsError {
		t.Fatalf("link %s -> %s for %s: %s", dependent, dependency, agent, resultText(t, result))
	}
}

// TestAgentEdgeIsScopedToItsAgent proves an agent-scoped edge is stored under
// that agent and does not leak into the global edge set.
func TestAgentEdgeIsScopedToItsAgent(t *testing.T) {
	server := newDependencyMcp(t)
	seedSkillRecord(t, server, "root")
	seedSkillRecord(t, server, "middle")

	linkSkillForAgent(t, server, "root", "middle", "scoped-agent", true)

	ctx := context.Background()
	root, err := server.skillByName(ctx, "root")
	if err != nil || root == nil {
		t.Fatalf("root missing: %v", err)
	}
	rows, err := server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependentId", root.Id)).Exec(ctx)
	if err != nil {
		t.Fatalf("read dependencies: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("edges = %d, want 1", len(rows))
	}
	if rows[0].Agent != "scoped-agent" {
		t.Fatalf("agent = %q, want scoped-agent", rows[0].Agent)
	}
}

// TestGlobalEdgeDefaultsToEmptyAgent proves an edge written without an agent
// lands on the global "" scope.
func TestGlobalEdgeDefaultsToEmptyAgent(t *testing.T) {
	server := newDependencyMcp(t)
	seedSkillRecord(t, server, "root")
	seedSkillRecord(t, server, "middle")

	linkSkill(t, server, "root", "middle", true)

	ctx := context.Background()
	root, err := server.skillByName(ctx, "root")
	if err != nil || root == nil {
		t.Fatalf("root missing: %v", err)
	}
	rows, err := server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependentId", root.Id)).Exec(ctx)
	if err != nil {
		t.Fatalf("read dependencies: %v", err)
	}
	if len(rows) != 1 || rows[0].Agent != "" {
		t.Fatalf("edges = %+v, want one global edge", rows)
	}
}

// TestAgentEdgeShadowsGlobalEdge proves an agent's own edge overrides the
// global edge for the same skill pair, in either insertion order.
func TestAgentEdgeShadowsGlobalEdge(t *testing.T) {
	for _, order := range []string{"global-first", "agent-first"} {
		t.Run(order, func(t *testing.T) {
			server := newDependencyMcp(t)
			root := seedSkillRecord(t, server, "root")
			seedSkillRecord(t, server, "middle")
			seedSkillRecord(t, server, "leaf")

			// the global edge is recursive, so leaf would be pulled in.
			// the agent's own edge is not, so leaf must not be.
			if order == "global-first" {
				linkSkill(t, server, "root", "middle", true)
				linkSkillForAgent(t, server, "root", "middle", "shadow-agent", false)
			} else {
				linkSkillForAgent(t, server, "root", "middle", "shadow-agent", false)
				linkSkill(t, server, "root", "middle", true)
			}
			linkSkill(t, server, "middle", "leaf", true)

			agentID := bindAgent(t, server, "shadow-agent", root.Id)
			got := boundNames(t, server, agentID)
			if !got["root"] || !got["middle"] {
				t.Fatalf("root and middle must be bound, got %v", got)
			}
			if got["leaf"] {
				t.Fatalf("the agent's non-recursive edge must shadow the global one, got %v", got)
			}
		})
	}
}

// TestGlobalEdgeAppliesToEveryAgent proves an agent with no edge of its own
// still inherits the global edge.
func TestGlobalEdgeAppliesToEveryAgent(t *testing.T) {
	server := newDependencyMcp(t)
	root := seedSkillRecord(t, server, "root")
	seedSkillRecord(t, server, "middle")

	linkSkill(t, server, "root", "middle", true)

	agentID := bindAgent(t, server, "inheriting-agent", root.Id)
	got := boundNames(t, server, agentID)
	if !got["root"] || !got["middle"] {
		t.Fatalf("the global edge must apply to every agent, got %v", got)
	}
}
