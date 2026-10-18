package web

import (
	"fmt"
	"net/http"
)

var Service *RiotWebUi

type WebServer struct {
	http.Handler
}

func init() {
	var err error
	Service, err = NewWebService()
	if err != nil {
		panic(fmt.Errorf("failed starting web service"))
	}
}

func (ws *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Service.Render(w, r)
}
