package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/verikono/glocalize/apps/riot/internal/web"
)

// TestStylesheetServedAsFile proves a stylesheet is served as a css file with
// file-like headers, and that the document links it rather than inlining it.
func TestStylesheetServedAsFile(t *testing.T) {
	ui, err := web.NewWebService()
	if err != nil {
		t.Fatalf("NewWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindMux(mux); err != nil {
		t.Fatalf("BindMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/styles/main.css")
	if err != nil {
		t.Fatalf("get stylesheet: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/css") {
		t.Fatalf("content-type = %q, want text/css", ct)
	}
	if resp.ContentLength <= 0 {
		t.Fatalf("content-length = %d, want > 0", resp.ContentLength)
	}

	doc, err := http.Get(server.URL + "/chart/")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	defer doc.Body.Close()

	buf := make([]byte, 4096)
	n, _ := doc.Body.Read(buf)
	body := string(buf[:n])

	if !strings.Contains(body, `<link rel="stylesheet" href="/styles/main.css">`) {
		t.Fatalf("document does not link the stylesheet: %s", body)
	}
	if strings.Contains(body, "<style>") {
		t.Fatalf("document still inlines a <style> block: %s", body)
	}
}

// TestUnknownStylesheetIsNotFound proves an unknown css path 404s rather than
// falling through to the document renderer.
func TestUnknownStylesheetIsNotFound(t *testing.T) {
	ui, err := web.NewWebService()
	if err != nil {
		t.Fatalf("NewWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindMux(mux); err != nil {
		t.Fatalf("BindMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/styles/nope.css")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
