package mcp

import (
	"strings"
	"testing"
)

// seedReference creates a static reference with the given body.
func seedReference(t *testing.T, server *RiotMcp, name, content string) {
	t.Helper()
	result := call(t, server.onAddReference, map[string]any{
		"name":         name,
		"description":  name + " reference",
		"accessType":   "static",
		"resourceType": "literature",
		"content":      content,
	})
	if result.IsError {
		t.Fatalf("create reference %s: %s", name, resultText(t, result))
	}
}

// TestGetReferenceReturnsBody proves a reference's content is readable through
// the MCP surface, which is what makes a skill able to keep its detail outside
// its own body.
func TestGetReferenceReturnsBody(t *testing.T) {
	server := newTestMcp(t)
	seedReference(t, server, "dynamodb-card", "on-demand capacity is the default")

	result := call(t, server.onGetReference, map[string]any{"name": "dynamodb-card"})
	if result.IsError {
		t.Fatalf("get-reference: %s", resultText(t, result))
	}
	document := resultText(t, result)
	for _, want := range []string{"name: dynamodb-card", "resourceType: literature", "on-demand capacity is the default"} {
		if !strings.Contains(document, want) {
			t.Fatalf("document missing %q:\n%s", want, document)
		}
	}
}

// TestGetReferenceUnknownNameIsAnError proves a missing reference is reported
// rather than returned as an empty document.
func TestGetReferenceUnknownNameIsAnError(t *testing.T) {
	server := newTestMcp(t)
	result := call(t, server.onGetReference, map[string]any{"name": "nope"})
	if !result.IsError {
		t.Fatalf("expected an error for an unknown reference, got: %s", resultText(t, result))
	}
}

// TestGetSkillListsAttachedReferences proves an agent reading a skill is told
// which references it may fetch, since a reference it is never told about is a
// reference it will never load.
func TestGetSkillListsAttachedReferences(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedSkill(t, server, "aws-database")
	seedReference(t, server, "aws-database-dynamodb", "card body")

	if result := call(t, server.onAssociateSkillReference, map[string]any{
		"skill": "aws-database", "reference": "aws-database-dynamodb",
	}); result.IsError {
		t.Fatalf("associate: %s", resultText(t, result))
	}

	result := call(t, server.onGetSkill, map[string]any{"project": "all", "name": "aws-database"})
	if result.IsError {
		t.Fatalf("get-skill: %s", resultText(t, result))
	}
	document := resultText(t, result)
	if !strings.Contains(document, "## References") {
		t.Fatalf("skill document has no references section:\n%s", document)
	}
	if !strings.Contains(document, "aws-database-dynamodb") {
		t.Fatalf("skill document does not name the attached reference:\n%s", document)
	}
	// the reference body must NOT be inlined: the point of a reference is that
	// it is fetched on demand.
	if strings.Contains(document, "card body") {
		t.Fatalf("skill document inlined the reference body:\n%s", document)
	}
}

// TestGetSkillWithoutReferencesHasNoSection proves the references section is
// absent rather than empty when a skill draws on nothing.
func TestGetSkillWithoutReferencesHasNoSection(t *testing.T) {
	server := newTestMcp(t)
	seedProject(t, server, "all")
	seedSkill(t, server, "plain-skill")

	result := call(t, server.onGetSkill, map[string]any{"project": "all", "name": "plain-skill"})
	if result.IsError {
		t.Fatalf("get-skill: %s", resultText(t, result))
	}
	if strings.Contains(resultText(t, result), "## References") {
		t.Fatalf("unexpected references section:\n%s", resultText(t, result))
	}
}
