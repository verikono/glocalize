package web

import (
	"fmt"
	"net/http"
)

// NewWebService builds the shared riot web UI: the /chart route, its body sections, and the in-memory style and script content.
func NewWebService() (*RiotWebUi, error) {
	ui := &RiotWebUi{
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
		styles:   StyleSheets,
		scripts:  ScriptFiles,
		dispatch: make(chan string, dispatchBuffer),
		hub:      newHub(),
	}
	go ui.run()
	return ui, nil
}

// Handler returns an http.Handler that serves the riot web UI. It lazily builds the shared Service and dispatches every request to it.
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

// BindMux mounts this instance's routes on mux. It registers the instance
// itself, not the package-level Service, so dispatch reaches the same hub the
// websocket clients are attached to.
func (w *RiotWebUi) BindMux(mux *http.ServeMux) error {
	routes, err := GetHandledRoutes(w.route)
	if err != nil {
		return err
	}
	for _, route := range routes {
		mux.Handle(route, w)
	}
	mux.Handle(wsPath, w.websocketHandler())
	return nil
}

// HandleWebsocket returns the handler that upgrades /ws to a websocket.
func (w *RiotWebUi) HandleWebsocket() http.Handler {
	return w.websocketHandler()
}

func GetHandledRoutes(basepath string) ([]string, error) {
	return []string{basepath, "/scripts/", "/styles/"}, nil
}
