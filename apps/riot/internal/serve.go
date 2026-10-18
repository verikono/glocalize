package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/verikono/glocalize/apps/riot/internal/gql"
	"github.com/verikono/glocalize/apps/riot/internal/mcp"
	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
	"github.com/verikono/glocalize/apps/riot/internal/utils"
	"github.com/verikono/glocalize/apps/riot/internal/web"
	T "github.com/verikono/glocalize/types/go"
)

// UrlPath is a URL path mounted on the single HTTP server, e.g. "/graph".
type UrlPath = string

// RiotConfig holds the server's wiring: the port and the mount points for the GraphQL and MCP surfaces (plus an optional playground path).
type RiotConfig struct {
	Port       int
	GqlURL     UrlPath
	McpURL     UrlPath
	Playground UrlPath
	UiURL      UrlPath
}

// RiotServer serves the GraphQL and MCP surfaces over one HTTP port, backed by the phi-generated orm client.
type RiotServer struct {
	config RiotConfig
	db     *orm.DB
	graph  *gql.RiotGraph
	mcps   *mcp.RiotMcp
	ui     *web.RiotWebUi
	mux    *http.ServeMux
	T.HotReloadable
}

// NewRiotServer opens the store, runs migrations, builds the GraphQL schema and the MCP tool surface, and wires both onto a single mux. It returns a fully-wired server (or an error) — never a half-built one.
func NewRiotServer(options ...T.OptionerFn) (*RiotServer, error) {

	config := RiotConfig{
		Port:   utils.ResolvePort(),
		GqlURL: "/graph",
		McpURL: "/mcp",
		UiURL:  "/chart",
	}

	config.applyEnvVars()

	for _, option := range options {
		if err := option(&config); err != nil {
			return nil, err
		}
	}

	provider, dsn, err := utils.ProviderAndDSN()
	if err != nil {
		return nil, err
	}
	if provider == "sqlite" {
		if err := utils.EnsureSQLiteDir(dsn); err != nil {
			return nil, fmt.Errorf("riot: create data dir: %w", err)
		}
	}

	db, err := orm.Open(provider, dsn)
	if err != nil {
		return nil, fmt.Errorf("riot: open %s store: %w", provider, err)
	}
	if err := db.RunMigrations(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("riot: migrate: %w", err)
	}

	graph, err := gql.NewRiotGraph(gql.WithDB(db))
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	mcps, err := mcp.NewRiotMcp(mcp.WithDB(db), mcp.WithProvider(provider))
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	wui, err := web.NewWebService()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	inst := &RiotServer{config: config, db: db, graph: graph, mcps: mcps, ui: wui}

	inst.mux = http.NewServeMux()
	inst.mux.Handle(config.GqlURL, graph.Handler())
	inst.mux.Handle("/voyager", graph.VoyagerHandler(config.GqlURL))
	inst.mux.Handle("/graphiql", graph.ExtendedGraphiQLHandler(config.GqlURL))
	inst.mux.Handle(config.McpURL, mcps.Handler())
	err = wui.BindMux(inst.mux)
	if config.Playground != "" {
		inst.mux.Handle(config.Playground, graph.PlaygroundHandler())
	}
	inst.mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","name":"riot"}`))
	})

	return inst, nil
}

// Handler returns the server's HTTP handler (the seam tests mount).
func (s *RiotServer) Handler() http.Handler {
	return s.mux
}

// GqlURL returns the GraphQL mount point.
func (s *RiotServer) GqlURL() string {
	return s.config.GqlURL
}

// McpURL returns the MCP mount point.
func (s *RiotServer) McpURL() string {
	return s.config.McpURL
}

// Port returns the port the server binds.
func (s *RiotServer) Port() int {
	return s.config.Port
}

// Start serves on the configured port and blocks until the server exits.
func (s *RiotServer) Start() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.mux,
	}
	return srv.ListenAndServe()
}

// Stop releases the underlying store connection.
func (s *RiotServer) Stop() error {
	s.Dispatch(T.HmrRefreshMessage) // tell connected clients to reload before we go away
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// Dispatch forwards a message to the web UI, which owns the websocket fan-out.
func (s *RiotServer) Dispatch(message string) {
	if s.ui != nil {
		s.ui.Dispatch(message)
	}
}

// OnDispatch handles a message addressed to the server. It forwards to the web
// UI so the same dispatch path carries both HMR and non-HMR traffic.
func (s *RiotServer) OnDispatch(message string) {
	if s.ui != nil {
		s.ui.OnDispatch(message)
	}
}
