package mcp

import (
	"strings"
	"testing"
)

// indexOf returns the position of needle in haystack, or -1.
func indexOf(haystack, needle string) int {
	return strings.Index(haystack, needle)
}

// bindEssential binds a skill to an agent at ESSENTIAL priority through the
// public handler, so the test exercises the same path an operator would.
func bindEssential(t *testing.T, server *RiotMcp, agentName, skillName string) {
	t.Helper()
	result := call(t, server.onAssociateAgentSkill, map[string]any{
		"agent":    agentName,
		"skill":    skillName,
		"priority": "ESSENTIAL",
	})
	if result.IsError {
		t.Fatalf("bind %s to %s: %s", skillName, agentName, resultText(t, result))
	}
}

// TestCheckinOrdersSkillsClosestFirst proves the check-in document lists a
// skill bound directly to the agent before the skills it depends on, and those
// before their own dependencies — so the agent can resolve a conflict by
// taking the first rule it reads.
func TestCheckinOrdersSkillsClosestFirst(t *testing.T) {
	server := newDependencyMcp(t)
	root := seedSkillRecord(t, server, "root")
	seedSkillRecord(t, server, "middle")
	seedSkillRecord(t, server, "leaf")

	// root -> middle -> leaf, all recursive, so the whole chain is bound.
	linkSkill(t, server, "root", "middle", true)
	linkSkill(t, server, "middle", "leaf", true)

	bindAgent(t, server, "ordering-agent", root.Id)
	// lift root to ESSENTIAL so the chain inherits ESSENTIAL and lands in one
	// bucket, making the ordering within that bucket the thing under test.
	bindEssential(t, server, "ordering-agent", "root")

	doc := resultText(t, call(t, server.onAgentCheckin, map[string]any{"agent": "ordering-agent", "project": "all"}))

	rootAt := indexOf(doc, "**root**")
	middleAt := indexOf(doc, "**middle**")
	leafAt := indexOf(doc, "**leaf**")
	if rootAt < 0 || middleAt < 0 || leafAt < 0 {
		t.Fatalf("document is missing a skill:\n%s", doc)
	}
	if !(rootAt < middleAt && middleAt < leafAt) {
		t.Fatalf("skills are not ordered closest-first (root=%d middle=%d leaf=%d):\n%s", rootAt, middleAt, leafAt, doc)
	}
}

// TestCheckinOrdersDirectBeforeDependencies proves a skill bound directly to
// the agent outranks a dependency of the same priority, even when the
// dependency's name sorts first alphabetically.
func TestCheckinOrdersDirectBeforeDependencies(t *testing.T) {
	server := newDependencyMcp(t)
	parent := seedSkillRecord(t, server, "zulu")
	seedSkillRecord(t, server, "alpha")

	linkSkill(t, server, "zulu", "alpha", true)

	bindAgent(t, server, "direct-first-agent", parent.Id)
	bindEssential(t, server, "direct-first-agent", "zulu")

	doc := resultText(t, call(t, server.onAgentCheckin, map[string]any{"agent": "direct-first-agent", "project": "all"}))

	zuluAt := indexOf(doc, "**zulu**")
	alphaAt := indexOf(doc, "**alpha**")
	if zuluAt < 0 || alphaAt < 0 {
		t.Fatalf("document is missing a skill:\n%s", doc)
	}
	if zuluAt > alphaAt {
		t.Fatalf("the directly bound skill must come first (zulu=%d alpha=%d):\n%s", zuluAt, alphaAt, doc)
	}
}

// TestCheckinDocumentExplainsConflictResolution proves the document tells the
// agent how to use the ordering to settle conflicting rules.
func TestCheckinDocumentExplainsConflictResolution(t *testing.T) {
	server := newDependencyMcp(t)
	root := seedSkillRecord(t, server, "root")
	bindAgent(t, server, "explained-agent", root.Id)

	doc := resultText(t, call(t, server.onAgentCheckin, map[string]any{"agent": "explained-agent", "project": "all"}))
	if !strings.Contains(doc, "conflict") {
		t.Fatalf("document does not explain conflict resolution:\n%s", doc)
	}
	if !strings.Contains(doc, "FIFO") {
		t.Fatalf("document does not state the FIFO rule:\n%s", doc)
	}
}

// TestCheckinDocumentMakesEssentialLoadsMandatory proves the document states
// the ESSENTIAL tier as a precondition of the session rather than a
// recommendation, and closes on that obligation. An agent that reads the
// document top to bottom must be left in no doubt that skipping a load is a
// protocol violation, not a judgement call.
func TestCheckinDocumentMakesEssentialLoadsMandatory(t *testing.T) {
	server := newDependencyMcp(t)
	root := seedSkillRecord(t, server, "root")
	bindAgent(t, server, "gated-agent", root.Id)
	bindEssential(t, server, "gated-agent", "root")

	doc := resultText(t, call(t, server.onAgentCheckin, map[string]any{"agent": "gated-agent", "project": "all"}))

	// the gate must be stated as a precondition, not a suggestion
	for _, want := range []string{
		"not complete until every ESSENTIAL skill",
		"a load is a tool call, not a recollection",
		"ESSENTIAL by inheritance",
		"stops the session",
	} {
		if !strings.Contains(doc, want) {
			t.Fatalf("document does not state %q:\n%s", want, doc)
		}
	}

	// the document must end on the obligation, after the skill list
	gateAt := indexOf(doc, "## Before you act")
	rootAt := indexOf(doc, "**root**")
	if gateAt < 0 {
		t.Fatalf("document has no closing obligation section:\n%s", doc)
	}
	if rootAt < 0 || gateAt < rootAt {
		t.Fatalf("the closing obligation must follow the skill list (root=%d gate=%d):\n%s", rootAt, gateAt, doc)
	}
	if !strings.Contains(doc[gateAt:], "Fetch all **1** ESSENTIAL skills") {
		t.Fatalf("closing section does not count the ESSENTIAL skills:\n%s", doc[gateAt:])
	}
}
