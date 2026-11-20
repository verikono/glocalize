package web

import (
	_ "embed"

	T "github.com/verikono/glocalize/types/go"
)

//go:embed scripts/projects.js
var projectsJs string

//go:embed style/projects.css
var projectsCss string

// ProjectScriptFiles and ProjectStyleSheets are keyed by their own /projects/ paths, so the projects UI can be mounted beside the agents and skills UIs without colliding on /scripts/ and /styles/.
var ProjectScriptFiles = map[T.UrlPath]WebScript{
	"/projects/scripts/projects.js": {name: "projects", route: "/projects/", urlPath: "/projects/scripts/projects.js", path: "./scripts/projects.js", language: CommonJavascript, content: Utf8String(projectsJs)},
	"/projects/scripts/d3.min.js":   {name: "d3", route: "/projects/", urlPath: "/projects/scripts/d3.min.js", path: "./scripts/d3.min.js", language: CommonJavascript, content: Utf8String(d3Js)},
}

var ProjectStyleSheets = map[T.UrlPath]StyleSheet{
	"/projects/styles/projects.css": {name: "projects", route: "/projects/", urlPath: "/projects/styles/projects.css", path: "./style/projects.css", content: CssContent(projectsCss)},
}

// NewProjectsWebService builds the /projects/ UI: a searchable project list on the left, a draggable separator, and the selected project's tree on the right.
func NewProjectsWebService() (*RiotWebUi, error) {
	ui := &RiotWebUi{
		route: "/projects/",
		sections: []RiotWebSection{
			{name: "navigation", classes: []string{"navigation"}, template: "navigation-projects"},
			{name: "projects", classes: []string{"projects"}, template: "projects"},
		},
		styles:   ProjectStyleSheets,
		scripts:  ProjectScriptFiles,
		dispatch: make(chan string, dispatchBuffer),
		hub:      newHub(),
	}
	go ui.run()
	return ui, nil
}
