package gql

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
	T "github.com/verikono/glocalize/types/go"
)

// The graph's contract is kept as three SDL documents beside this file: the
// generated type system, the queries the browser issues, and the mutations it
// may call. They are embedded so the binary carries its own contract and the
// documents can be served or diffed without a filesystem read.
var (
	//go:embed riot.graphql
	riotSDL string

	//go:embed queries.graphql
	queriesSDL string

	//go:embed mutations.graphql
	mutationsSDL string
)

// Documents returns the three SDL documents the graph is built from, keyed by
// file name. The executable schema is assembled in schema.go; these are the
// human-readable contract it implements.
func Documents() map[string]string {
	return map[string]string{
		"riot.graphql":      riotSDL,
		"queries.graphql":   queriesSDL,
		"mutations.graphql": mutationsSDL,
	}
}

// RiotGraph is the GraphQL surface over the riot orm client.
type RiotGraph struct {
	db     *orm.DB
	schema graphql.Schema
}

// NewRiotGraph builds the executable schema and returns a serving-ready graph.
func NewRiotGraph(options ...T.OptionerFn) (*RiotGraph, error) {
	inst := &RiotGraph{}
	for _, option := range options {
		if err := option(inst); err != nil {
			return nil, err
		}
	}
	if inst.db == nil {
		return nil, fmt.Errorf("gql: an orm client is required (WithDB)")
	}
	schema, err := buildSchema(inst.db)
	if err != nil {
		return nil, err
	}
	inst.schema = schema
	return inst, nil
}

// WithDB injects the orm client the resolvers read from.
func WithDB(db *orm.DB) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(*RiotGraph); ok {
				inst.db = db
			}
		}
		return nil
	}
}

// Handler returns the primary GraphQL HTTP handler.
func (g *RiotGraph) Handler() http.Handler {
	return handler.New(&handler.Config{
		Schema:   &g.schema,
		Pretty:   true,
		GraphiQL: false,
	})
}

// PlaygroundHandler returns the basic bundled GraphiQL playground.
func (g *RiotGraph) PlaygroundHandler() http.Handler {
	return handler.New(&handler.Config{
		Schema:   &g.schema,
		Pretty:   true,
		GraphiQL: true,
	})
}

// VoyagerHandler returns an interactive GraphQL Voyager node-graph visualization UI.
// The endpoint argument should point to your GraphQL query endpoint (e.g., "/query" or "/graphql").
func (g *RiotGraph) VoyagerHandler(endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>GraphQL Voyager</title>
  <style>
    html, body, #voyager {
      height: 100%%;
      width: 100%%;
      margin: 0;
      padding: 0;
      overflow: hidden;
    }
  </style>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-voyager/dist/voyager.css" />
  <script src="https://cdn.jsdelivr.net/npm/graphql-voyager/dist/voyager.standalone.js"></script>
</head>
<body>
  <div id="voyager">Loading GraphQL Voyager...</div>
  <script>
    GraphQLVoyager.renderVoyager(document.getElementById('voyager'), {
      introspection: fetch('%s', {
        method: 'post',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query: GraphQLVoyager.voyagerIntrospectionQuery }),
        credentials: 'same-origin',
      }).then(function(response) {
        return response.json();
      })
    });
  </script>
</body>
</html>`, endpoint)
	})
}

// ExtendedGraphiQLHandler serves GraphiQL with the GraphiQL Explorer sidebar plugin.
func (g *RiotGraph) ExtendedGraphiQLHandler(endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>Extended GraphiQL</title>
  <link rel="stylesheet" href="https://unpkg.com/graphiql@3.2.0/graphiql.min.css" />
  <link rel="stylesheet" href="https://unpkg.com/@graphiql/plugin-explorer@2.0.0/dist/style.css" />
  <style>
    body { height: 100vh; margin: 0; overflow: hidden; }
    #graphiql { height: 100%%; }
  </style>
  <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
  <script src="https://unpkg.com/graphiql@3.2.0/graphiql.min.js"></script>
  <script src="https://unpkg.com/@graphiql/plugin-explorer@2.0.0/dist/index.umd.js"></script>
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
      })
    );
  </script>
</body>
</html>`, endpoint)
	})
}

// Schema exposes the executable schema (for tests and introspection).
func (g *RiotGraph) Schema() graphql.Schema {
	return g.schema
}
