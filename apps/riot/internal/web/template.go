package web

import (
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/chart.gotpl
var chartGotpl string

// chartTemplate is the document shell. The sections are rendered into #root,
// never directly into <body>.
var chartTemplate = template.Must(template.New("chart").Parse(chartGotpl))

// chartTemplateData is the data handed to chartTemplate: Head is the rendered
// <head> element and Root is the rendered sections, placed inside #root.
type chartTemplateData struct {
	Head template.HTML
	Root template.HTML
}

// sectionTemplates is the embedded section template filesystem, rooted at
// templates/, so a section named "main" resolves to templates/main.gotpl.
//
//go:embed templates/*.gotpl
var sectionTemplates embed.FS

// sectionTemplate renders the named section template and returns its HTML.
// The name is resolved against templates/<name>.gotpl; a missing or malformed
// template is returned as an error so the caller can 500 the request.
func sectionTemplate(name string) (string, error) {
	path := "templates/" + name + ".gotpl"
	raw, err := sectionTemplates.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("section template %q: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("section template %q: %w", name, err)
	}

	var b strings.Builder
	if err := tmpl.Execute(&b, nil); err != nil {
		return "", fmt.Errorf("section template %q: %w", name, err)
	}
	return b.String(), nil
}
