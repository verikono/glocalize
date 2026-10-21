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

	for _, want := range []string{`id="agents-list"`, `id="agents-search"`, `id="agents-resizer"`, `id="agent-name"`, `id="agent-graph"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("document is missing %s: %s", want, body)
		}
	}
}

// TestAgentsDocumentServesD3 proves the agents page loads d3, which the skill tree needs, and that it is served from the page's own path rather than the skills UI's /scripts/.
func TestAgentsDocumentServesD3(t *testing.T) {
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

	resp, err := http.Get(server.URL + "/agents/scripts/d3.min.js")
	if err != nil {
		t.Fatalf("get d3: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.ContentLength <= 0 {
		t.Fatalf("content-length = %d, want > 0", resp.ContentLength)
	}

	doc, err := http.Get(server.URL + "/agents/")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	defer doc.Body.Close()

	buf := make([]byte, 8192)
	n, _ := doc.Body.Read(buf)
	if !strings.Contains(string(buf[:n]), `src="/agents/scripts/d3.min.js"`) {
		t.Fatalf("document does not load d3: %s", string(buf[:n]))
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
