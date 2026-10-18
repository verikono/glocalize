package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	_ "modernc.org/sqlite"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// newTestMcp opens a throwaway sqlite store, migrates it and returns a server
// wired to it.
func newTestMcp(t *testing.T) *RiotMcp {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "riot-test.sqlite")
	db, err := orm.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.RunMigrations(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	server, err := NewRiotMcp(WithDB(db), WithProvider("sqlite"))
	if err != nil {
		t.Fatalf("NewRiotMcp: %v", err)
	}
	return server
}

// call invokes a handler with the given arguments.
func call(t *testing.T, handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error), args map[string]any) *mcp.CallToolResult {
	t.Helper()
	request := mcp.CallToolRequest{}
	request.Params.Arguments = args
	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	return result
}

// resultText returns the first text block of a tool result.
func resultText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatalf("result has no content")
	}
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("result content is not text: %T", result.Content[0])
	}
	return text.Text
}

// TestMutateSkillCreateUpdateDelete walks a skill through its whole lifecycle
// and proves keywords and both association sides are written.
func TestMutateSkillCreateUpdateDelete(t *testing.T) {
	server := newTestMcp(t)
	ctx := context.Background()

	// two prerequisites the subject will depend on.
	for _, name := range []string{"alpha", "beta"} {
		result := call(t, server.onMutateSkill, map[string]any{
			"name": name,
			"mutation": map[string]any{
				"op":   "create",
				"data": map[string]any{"description": name + " skill", "content": "body", "keywords": "seed"},
			},
		})
		if result.IsError {
			t.Fatalf("create %s: %s", name, resultText(t, result))
		}
	}

	// create the subject with a keyword and one dependency.
	result := call(t, server.onMutateSkill, map[string]any{
		"name": "gamma",
		"mutation": map[string]any{
			"op":   "create",
			"data": map[string]any{"description": "gamma skill", "content": "body", "keywords": "searchable token"},
			"dependencies": map[string]any{
				"add": []any{"alpha"},
			},
		},
	})
	if result.IsError {
		t.Fatalf("create gamma: %s", resultText(t, result))
	}

	gamma, err := server.skillByName(ctx, "gamma")
	if err != nil || gamma == nil {
		t.Fatalf("gamma not stored: %v", err)
	}
	if gamma.Keywords != "searchable token" {
		t.Fatalf("keywords = %q, want %q", gamma.Keywords, "searchable token")
	}

	deps, err := server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependentId", gamma.Id)).Exec(ctx)
	if err != nil {
		t.Fatalf("read dependencies: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("dependencies = %d, want 1", len(deps))
	}

	// update: replace the dependency set with beta and change the keywords.
	result = call(t, server.onMutateSkill, map[string]any{
		"name": "gamma",
		"mutation": map[string]any{
			"op":   "update",
			"data": map[string]any{"keywords": "retuned"},
			"dependencies": map[string]any{
				"add":     []any{"beta"},
				"replace": true,
			},
		},
	})
	if result.IsError {
		t.Fatalf("update gamma: %s", resultText(t, result))
	}

	gamma, err = server.skillByName(ctx, "gamma")
	if err != nil || gamma == nil {
		t.Fatalf("gamma missing after update: %v", err)
	}
	if gamma.Keywords != "retuned" {
		t.Fatalf("keywords = %q, want %q", gamma.Keywords, "retuned")
	}
	deps, err = server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependentId", gamma.Id)).Exec(ctx)
	if err != nil {
		t.Fatalf("read dependencies: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("dependencies after replace = %d, want 1", len(deps))
	}
	beta, err := server.skillByName(ctx, "beta")
	if err != nil || beta == nil {
		t.Fatalf("beta missing: %v", err)
	}
	if deps[0].DependencyId != beta.Id {
		t.Fatalf("dependency = %q, want beta %q", deps[0].DependencyId, beta.Id)
	}

	// dependents: alpha now depends on gamma.
	result = call(t, server.onMutateSkill, map[string]any{
		"name": "gamma",
		"mutation": map[string]any{
			"op":         "update",
			"dependents": map[string]any{"add": []any{"alpha"}},
		},
	})
	if result.IsError {
		t.Fatalf("add dependent: %s", resultText(t, result))
	}
	alpha, err := server.skillByName(ctx, "alpha")
	if err != nil || alpha == nil {
		t.Fatalf("alpha missing: %v", err)
	}
	rows, err := server.db.SkillDependency.FindMany(
		eq[orm.SkillDependency, string]("dependentId", alpha.Id),
		eq[orm.SkillDependency, string]("dependencyId", gamma.Id),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("read dependents: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("dependent rows = %d, want 1", len(rows))
	}

	// delete: the skill and every association pointing at it go.
	result = call(t, server.onMutateSkill, map[string]any{
		"name":     "gamma",
		"mutation": map[string]any{"op": "delete"},
	})
	if result.IsError {
		t.Fatalf("delete gamma: %s", resultText(t, result))
	}
	if gone, err := server.skillByName(ctx, "gamma"); err != nil || gone != nil {
		t.Fatalf("gamma still present: %v", err)
	}
	left, err := server.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependencyId", gamma.Id)).Exec(ctx)
	if err != nil {
		t.Fatalf("read dependencies: %v", err)
	}
	if len(left) != 0 {
		t.Fatalf("dangling dependency rows = %d, want 0", len(left))
	}
}

// TestMutateSkillRequiresTarget proves a mutation without an id or name is
// rejected rather than silently doing nothing.
func TestMutateSkillRequiresTarget(t *testing.T) {
	server := newTestMcp(t)
	result := call(t, server.onMutateSkill, map[string]any{
		"mutation": map[string]any{"op": "update"},
	})
	if !result.IsError {
		t.Fatalf("expected an error result, got %s", resultText(t, result))
	}
}

// TestMutateSkillRejectsUnknownOp proves an unknown op is reported.
func TestMutateSkillRejectsUnknownOp(t *testing.T) {
	server := newTestMcp(t)
	result := call(t, server.onMutateSkill, map[string]any{
		"name":     "anything",
		"mutation": map[string]any{"op": "upsert"},
	})
	if !result.IsError {
		t.Fatalf("expected an error result, got %s", resultText(t, result))
	}
}

// TestMutateSkillResultShape proves the result carries the op and the record.
func TestMutateSkillResultShape(t *testing.T) {
	server := newTestMcp(t)
	result := call(t, server.onMutateSkill, map[string]any{
		"name": "delta",
		"mutation": map[string]any{
			"op":   "create",
			"data": map[string]any{"description": "delta skill", "content": "body", "keywords": "k1 k2"},
		},
	})
	if result.IsError {
		t.Fatalf("create delta: %s", resultText(t, result))
	}
	var payload struct {
		Op    string `json:"op"`
		Skill struct {
			Name     string `json:"name"`
			Keywords string `json:"keywords"`
		} `json:"skill"`
	}
	if err := json.Unmarshal([]byte(resultText(t, result)), &payload); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if payload.Op != "create" {
		t.Fatalf("op = %q, want create", payload.Op)
	}
	if payload.Skill.Name != "delta" || payload.Skill.Keywords != "k1 k2" {
		t.Fatalf("skill = %+v", payload.Skill)
	}
}
