package mcp

import (
	"strings"
	"testing"
)

func TestRenderAgentFile_WithTools(t *testing.T) {
	out, err := renderAgentFile("golang-developer", "Focused Go agent.", `["create-agent","update-agent"]`)
	if err != nil {
		t.Fatalf("renderAgentFile() error = %v", err)
	}
	if !strings.Contains(out, "name: golang-developer\n") {
		t.Fatalf("missing name line:\n%s", out)
	}
	if !strings.Contains(out, "description: Focused Go agent.\n") {
		t.Fatalf("missing description line:\n%s", out)
	}
	if !strings.Contains(out, "allowed-tools:\n  - create-agent\n  - update-agent\n") {
		t.Fatalf("missing allowed-tools block:\n%s", out)
	}
	if !strings.Contains(out, "CHECKIN TO RIOT-MCP WITH THIS AGENT NAME.") {
		t.Fatalf("missing check-in directive:\n%s", out)
	}
}

func TestRenderAgentFile_NoTools(t *testing.T) {
	out, err := renderAgentFile("golang-developer", "Focused Go agent.", "[]")
	if err != nil {
		t.Fatalf("renderAgentFile() error = %v", err)
	}
	if strings.Contains(out, "allowed-tools") {
		t.Fatalf("allowed-tools should be omitted for an empty list:\n%s", out)
	}
}

func TestRenderAgentFile_NonArrayTools(t *testing.T) {
	out, err := renderAgentFile("golang-developer", "Focused Go agent.", `{}`)
	if err != nil {
		t.Fatalf("renderAgentFile() error = %v", err)
	}
	if strings.Contains(out, "allowed-tools") {
		t.Fatalf("allowed-tools should be omitted for a non-array value:\n%s", out)
	}
}

func TestParseToolList(t *testing.T) {
	tools, err := parseToolList(`["a","b"]`)
	if err != nil {
		t.Fatalf("parseToolList() error = %v", err)
	}
	if len(tools) != 2 || tools[0] != "a" || tools[1] != "b" {
		t.Fatalf("parseToolList() = %v, want [a b]", tools)
	}

	empty, err := parseToolList("[]")
	if err != nil {
		t.Fatalf("parseToolList([]) error = %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("parseToolList([]) = %v, want empty", empty)
	}

	nonArray, err := parseToolList(`{"x":1}`)
	if err != nil {
		t.Fatalf("parseToolList(obj) error = %v", err)
	}
	if len(nonArray) != 0 {
		t.Fatalf("parseToolList(obj) = %v, want empty", nonArray)
	}
}
