package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	T "github.com/verikono/glocalize/types/go"
)

// Document Base struct - pairs to routes - eg Request("/chart/") <- RiotWebUi
type RiotWebUi struct {
	route    T.UrlPath
	styles   map[T.UrlPath]StyleSheet
	scripts  map[T.UrlPath]WebScript
	sections []RiotWebSection
	server   WebServer
	dispatch chan string
	hub      *hub
	T.HotReloadable
}

// ServeHTTP routes the request: a known script or stylesheet urlPath has its contents written out as if it were a file, the document route renders the page, and everything else is not found.
func (rw *RiotWebUi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if script, ok := rw.scripts[r.URL.Path]; ok {
		rw.writeScript(w, r, script)
		return
	}
	if sheet, ok := rw.styles[r.URL.Path]; ok {
		rw.writeStyle(w, r, sheet)
		return
	}
	if !rw.isUiPath(r.URL.Path) {
		http.NotFound(w, r)
		return
	}
	rw.Render(w, r)
}

// isUiPath reports whether path is the document route ("/chart" or any path beneath it).
func (rw *RiotWebUi) isUiPath(path string) bool {
	return path == strings.TrimSuffix(rw.route, "/") || strings.HasPrefix(path, rw.route)
}

// writeScript writes the embedded script out as though it were a file served from disk: content-type and content-length headers, then the source body.
func (rw *RiotWebUi) writeScript(w http.ResponseWriter, r *http.Request, script WebScript) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	http.ServeContent(w, r, script.urlPath, time.Time{}, bytes.NewReader([]byte(script.GetContent())))
}

// writeStyle writes the embedded stylesheet out as though it were a css file served from disk: content-type and content-length headers, then the source body.
func (rw *RiotWebUi) writeStyle(w http.ResponseWriter, r *http.Request, sheet StyleSheet) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	http.ServeContent(w, r, sheet.urlPath, time.Time{}, bytes.NewReader([]byte(sheet.GetContent())))
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

// renderDocument assembles the full HTML document from the chart template:
// the rendered <head> and the sections, which are placed inside #root.
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
	if err := chartTemplate.Execute(&b, chartTemplateData{Head: template.HTML(headContent), Root: template.HTML(bodyContent)}); err != nil {
		return "", fmt.Errorf("execute chart template: %w", err)
	}
	return b.String(), nil
}

// renderHead renders <head> with the styles and scripts referenced by url, so
// the browser fetches each one from the server.
func (w *RiotWebUi) renderHead() (string, error) {
	var b strings.Builder
	b.WriteString("<head>")

	for _, sheet := range w.styles {
		b.WriteString(`<link rel="stylesheet" href="`)
		b.WriteString(sheet.urlPath)
		b.WriteString(`">`)
	}

	for _, script := range w.scripts {
		if script.language == EcmascriptModule {
			b.WriteString(fmt.Sprintf(`<script src="%s" type="module"></script>`, script.urlPath))
		} else {
			b.WriteString(fmt.Sprintf(`<script src="%s"></script>`, script.urlPath))
		}
	}

	b.WriteString("</head>")
	return b.String(), nil
}

// renderBody renders the sections as the document's root content. The sections
// are placed inside #root by the chart template, not directly in <body>.
func (w *RiotWebUi) renderBody() (string, error) {
	var b strings.Builder

	for _, section := range w.sections {
		b.WriteString(section.render())
	}

	return b.String(), nil
}
