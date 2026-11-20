package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/verikono/glocalize/apps/riot/internal/web"
)

// TestProjectsDocumentServes proves the projects UI renders its own document, script and stylesheet, and that it does not claim the skills UI's /scripts/ and /styles/ paths.
func TestProjectsDocumentServes(t *testing.T) {
	ui, err := web.NewProjectsWebService()
	if err != nil {
		t.Fatalf("NewProjectsWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindDocumentMux(mux); err != nil {
		t.Fatalf("BindDocumentMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	for _, path := range []string{"/projects/", "/projects/scripts/projects.js", "/projects/styles/projects.css"} {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatalf("get %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, resp.StatusCode, http.StatusOK)
		}
	}

	doc, err := http.Get(server.URL + "/projects/")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	defer doc.Body.Close()

	buf := make([]byte, 8192)
	n, _ := doc.Body.Read(buf)
	body := string(buf[:n])

	for _, want := range []string{`id="projects-list"`, `id="projects-search"`, `id="projects-resizer"`, `id="project-name"`, `id="project-tree"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("document is missing %s: %s", want, body)
		}
	}
}

// TestProjectsDocumentServesD3 proves the projects page loads d3, which the project tree needs, and that it is served from the page's own path rather than the skills UI's /scripts/.
func TestProjectsDocumentServesD3(t *testing.T) {
	ui, err := web.NewProjectsWebService()
	if err != nil {
		t.Fatalf("NewProjectsWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindDocumentMux(mux); err != nil {
		t.Fatalf("BindDocumentMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/projects/scripts/d3.min.js")
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

	doc, err := http.Get(server.URL + "/projects/")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	defer doc.Body.Close()

	buf := make([]byte, 8192)
	n, _ := doc.Body.Read(buf)
	if !strings.Contains(string(buf[:n]), `src="/projects/scripts/d3.min.js"`) {
		t.Fatalf("document does not load d3: %s", string(buf[:n]))
	}
}

// TestProjectsDoesNotClaimSkillsPaths proves the projects UI leaves /scripts/ and /styles/ to the skills UI, so both can be mounted on one mux.
func TestProjectsDoesNotClaimSkillsPaths(t *testing.T) {
	ui, err := web.NewProjectsWebService()
	if err != nil {
		t.Fatalf("NewProjectsWebService: %v", err)
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

// TestNavigationListsProjectsFirst proves the nav order is projects, agents, skills — the order is deliberate, so presence alone is not enough.
func TestNavigationListsProjectsFirst(t *testing.T) {
	ui, err := web.NewProjectsWebService()
	if err != nil {
		t.Fatalf("NewProjectsWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindDocumentMux(mux); err != nil {
		t.Fatalf("BindDocumentMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	doc, err := http.Get(server.URL + "/projects/")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	defer doc.Body.Close()

	buf := make([]byte, 8192)
	n, _ := doc.Body.Read(buf)
	body := string(buf[:n])

	projects := strings.Index(body, `href="/projects/"`)
	agents := strings.Index(body, `href="/agents/"`)
	skills := strings.Index(body, `href="/skills/"`)
	if projects < 0 || agents < 0 || skills < 0 {
		t.Fatalf("document is missing a nav link (projects=%d agents=%d skills=%d): %s", projects, agents, skills, body)
	}
	if !(projects < agents && agents < skills) {
		t.Fatalf("nav order is not projects, agents, skills (projects=%d agents=%d skills=%d): %s", projects, agents, skills, body)
	}
}
