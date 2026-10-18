package web

import (
	"fmt"
	"html/template"
	"strings"
)

type RiotWebSection struct {
	name     string
	element  string // eg div, section etc.
	classes  []string
	template string
}

// render returns the section as an HTML element, defaulting to a <div> and
// using the section name as the element id.
func (r *RiotWebSection) render() string {
	element := r.element
	if element == "" {
		element = "div"
	}

	var b strings.Builder
	b.WriteString("<")
	b.WriteString(element)

	if r.name != "" {
		b.WriteString(` id="`)
		b.WriteString(r.name)
		b.WriteString(`"`)
	}

	if len(r.classes) > 0 {
		b.WriteString(` class="`)
		b.WriteString(strings.Join(r.classes, " "))
		b.WriteString(`"`)
	}

	b.WriteString(">")
	// render the template
	b.WriteString(r.getTemplate())
	b.WriteString(fmt.Sprintf("</%s>", element))

	return b.String()
}

// getTemplate returns the section's rendered template content. It resolves the
// template name from r.template, falling back to r.name, and reads
// ./templates/<name>.gotpl. When neither is set it returns a placeholder; when
// the named template is missing or fails to parse it returns an error so the
// caller can 500 the request.
func (r *RiotWebSection) getTemplate() string {
	name := r.template
	if name == "" {
		name = r.name
	}
	if name == "" {
		return `<p class="errored">no content</p>`
	}

	content, err := sectionTemplate(name)
	if err != nil {
		return fmt.Sprintf(`<p class="errored">%s</p>`, template.HTMLEscapeString(err.Error()))
	}
	return content
}
