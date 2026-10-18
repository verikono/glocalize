package gql

import (
	"fmt"
	"net/http"
)

func GraphiQLExtendedHandler(endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>Extended GraphiQL</title>
  
  <!-- GraphiQL Core & React CSS -->
  <link rel="stylesheet" href="https://unpkg.com/graphiql/graphiql.min.css" />
  <link rel="stylesheet" href="https://unpkg.com/@graphiql/plugin-explorer/dist/style.css" />
  
  <style>
    body { height: 100vh; margin: 0; overflow: hidden; }
    #graphiql { height: 100%%; }
  </style>

  <!-- React & GraphiQL Core JS -->
  <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
  <script src="https://unpkg.com/graphiql/graphiql.min.js"></script>
  
  <!-- GraphiQL Explorer Plugin -->
  <script src="https://unpkg.com/@graphiql/plugin-explorer/dist/index.umd.js"></script>
</head>
<body>
  <div id="graphiql">Loading...</div>
  <script>
    const root = ReactDOM.createRoot(document.getElementById('graphiql'));
    const fetcher = GraphiQL.createFetcher({ url: '%s' });
    const explorerPlugin = GraphiQLPluginExplorer.explorerPlugin();

    root.render(
      React.createElement(GraphiQL, {
        fetcher: fetcher,
        plugins: [explorerPlugin],
        defaultQuery: 'query { hello }',
      })
    );
  </script>
</body>
</html>`, endpoint)
	}
}

func main() {
	// Mount the extended GraphiQL route alongside your GraphQL server
	http.HandleFunc("/graphiql", GraphiQLExtendedHandler("/query"))
	http.ListenAndServe(":8080", nil)
}
