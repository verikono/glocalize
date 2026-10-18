package web

import "strings"

type RiotWebSection struct {
	name    string
	element string // eg div, section etc.
	classes []string
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

	b.WriteString("></")
	b.WriteString(element)
	b.WriteString(">")

	return b.String()
}
