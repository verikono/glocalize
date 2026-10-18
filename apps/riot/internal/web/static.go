package web

import "net/http"

// StaticHandler serves the in-memory scripts at their urlPath. It is mounted
// under /scripts/ by the server, so the document can reference each script by
// URL instead of inlining its source.
func StaticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		script, ok := ScriptFiles[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = w.Write([]byte(script.GetContent()))
	})
}
