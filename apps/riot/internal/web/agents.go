package web

import (
	_ "embed"
	"net/http"

	T "github.com/verikono/glocalize/types/go"
)

//go:embed scripts/agents.js
var agentsJs string

//go:embed style/agents.css
var agentsCss string

// AgentScriptFiles and AgentStyleSheets are keyed by their own /agents/ paths, so the agents UI can be mounted beside the skills UI without colliding on /scripts/ and /styles/.
var AgentScriptFiles = map[T.UrlPath]WebScript{
	"/agents/scripts/agents.js": {name: "agents", route: "/agents/", urlPath: "/agents/scripts/agents.js", path: "./scripts/agents.js", language: CommonJavascript, content: Utf8String(agentsJs)},
	"/agents/scripts/d3.min.js": {name: "d3", route: "/agents/", urlPath: "/agents/scripts/d3.min.js", path: "./scripts/d3.min.js", language: CommonJavascript, content: Utf8String(d3Js)},
}

var AgentStyleSheets = map[T.UrlPath]StyleSheet{
	"/agents/styles/agents.css": {name: "agents", route: "/agents/", urlPath: "/agents/styles/agents.css", path: "./style/agents.css", content: CssContent(agentsCss)},
}

// NewAgentsWebService builds the /agents/ UI: a searchable agent list on the left, a draggable separator, and the selected agent's detail on the right.
func NewAgentsWebService() (*RiotWebUi, error) {
	ui := &RiotWebUi{
		route: "/agents/",
		sections: []RiotWebSection{
			{name: "navigation", classes: []string{"navigation"}, template: "navigation-agents"},
			{name: "agents", classes: []string{"agents"}, template: "agents"},
		},
		styles:   AgentStyleSheets,
		scripts:  AgentScriptFiles,
		dispatch: make(chan string, dispatchBuffer),
		hub:      newHub(),
	}
	go ui.run()
	return ui, nil
}

// BindDocumentMux mounts only this instance's document route and its own script and style paths, so a second UI can share the server without claiming /scripts/ and /styles/.
func (w *RiotWebUi) BindDocumentMux(mux *http.ServeMux) error {
	mux.Handle(w.route, w)
	for path := range w.scripts {
		mux.Handle(path, w)
	}
	for path := range w.styles {
		mux.Handle(path, w)
	}
	return nil
}
