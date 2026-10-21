package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/verikono/glocalize/apps/riot/internal/web"
)

// TestAgentsDocumentServes proves the agents UI renders its own document, script and stylesheet, and that it does not claim the skills UI's /scripts/ and /styles/ paths.
func TestAgentsDocumentServes(t *testing.T) {
	ui, err := web.NewAgentsWebService()
	if err != nil {
		t.Fatalf("NewAgentsWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindDocumentMux(mux); err != nil {
		t.Fatalf("BindDocumentMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	for _, path := range []string{"/agents/", "/agents/scripts/agents.js", "/agents/styles/agents.css"} {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatalf("get %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, resp.StatusCode, http.StatusOK)
		}
	}

	doc, err := http.Get(server.URL + "/agents/")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	defer doc.Body.Close()

	buf := make([]byte, 8192)
	n, _ := doc.Body.Read(buf)
	body := string(buf[:n])

	for _, want := range []string{`id="agents-list"`, `id="agents-search"`, `id="agents-resizer"`, `id="agent-name"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("document is missing %s: %s", want, body)
		}
	}
}

// TestAgentsDoesNotClaimSkillsPaths proves the agents UI leaves /scripts/ and /styles/ to the skills UI, so both can be mounted on one mux.
func TestAgentsDoesNotClaimSkillsPaths(t *testing.T) {
	ui, err := web.NewAgentsWebService()
	if err != nil {
		t.Fatalf("NewAgentsWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindDocumentMux(mux); err != nil {
		t.Fatalf("BindDocumentMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/scripts/riot.js")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
