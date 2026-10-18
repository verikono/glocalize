package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	T "github.com/verikono/glocalize/types/go"
)

type RiotWebUi struct {
	route    T.UrlPath
	styles   map[T.UrlPath]StyleSheet
	scripts  map[T.UrlPath]WebScript
	sections []RiotWebSection
	server   any // any webserver is fine, it needs to call RiotWebUi.Render(w http.ResponseWriter, r *http.Request) on every request
}

// NewWebService builds the shared riot web UI: the /chart route, its body
// sections, and the in-memory style and script content.
func NewWebService() (*RiotWebUi, error) {
	return &RiotWebUi{
		route: "/chart/",
		sections: []RiotWebSection{
			{
				name:    "main",
				classes: []string{"main"},
			},
			{
				name:    "chooser",
				classes: []string{"selector"},
			},
		},
		styles:  StyleSheets,
		scripts: ScriptFiles,
	}, nil
}

// Handler returns an http.Handler that serves the riot web UI. It lazily
// builds the shared Service and dispatches every request to it.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Service == nil {
			var err error
			Service, err = NewWebService()
			if err != nil {
				http.Error(w, fmt.Sprintf("riot: build web ui: %v", err), http.StatusInternalServerError)
				return
			}
		}
		Service.ServeHTTP(w, r)
	})
}

// ServeHTTP routes the request: a known script urlPath has its contents
// written out as if it were a file, the document route renders the page, and
// everything else is not found.
func (rw *RiotWebUi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if script, ok := rw.scripts[r.URL.Path]; ok {
		rw.writeScript(w, r, script)
		return
	}
	if !rw.isUiPath(r.URL.Path) {
		http.NotFound(w, r)
		return
	}
	rw.Render(w, r)
}

// isUiPath reports whether path is the document route ("/chart" or any path
// beneath it).
func (rw *RiotWebUi) isUiPath(path string) bool {
	return path == strings.TrimSuffix(rw.route, "/") || strings.HasPrefix(path, rw.route)
}

// writeScript writes the embedded script out as though it were a file served
// from disk: content-type and content-length headers, then the source body.
func (rw *RiotWebUi) writeScript(w http.ResponseWriter, r *http.Request, script WebScript) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	http.ServeContent(w, r, script.urlPath, time.Time{}, bytes.NewReader([]byte(script.GetContent())))
}

// Render writes the full HTML document for the route to w.
func (rw *RiotWebUi) Render(w http.ResponseWriter, r *http.Request) {
	rendition, err := rw.renderDocument()
	if err != nil {
		http.Error(w, fmt.Sprintf("riot: render: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(rendition))
}

// renderDocument assembles the full HTML document: doctype, <html>, the
// rendered <head> and <body>.
func (w *RiotWebUi) renderDocument() (string, error) {
	headContent, err := w.renderHead()
	if err != nil {
		return "", err
	}
	bodyContent, err := w.renderBody()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("<!doctype html>\n<html>\n")
	b.WriteString(headContent)
	b.WriteString("\n")
	b.WriteString(bodyContent)
	b.WriteString("\n</html>\n")
	return b.String(), nil
}

// renderHead renders <head> with the styles inlined and the scripts referenced
// by url, so the browser fetches each script from the server.
func (w *RiotWebUi) renderHead() (string, error) {
	var b strings.Builder
	b.WriteString("<head>")

	for _, sheet := range w.styles {
		b.WriteString("<style>")
		b.WriteString(sheet.GetContent())
		b.WriteString("</style>")
	}

	for _, script := range w.scripts {
		b.WriteString(`<script src="`)
		b.WriteString(script.urlPath)
		b.WriteString(`"></script>`)
	}

	b.WriteString("</head>")
	return b.String(), nil
}

// renderBody renders <body> with each section element.
func (w *RiotWebUi) renderBody() (string, error) {
	var b strings.Builder
	b.WriteString("<body>")

	for _, section := range w.sections {
		b.WriteString(section.render())
	}

	b.WriteString("</body>")
	return b.String(), nil
}
