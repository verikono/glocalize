package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/verikono/glocalize/apps/riot/internal/gql"
	"github.com/verikono/glocalize/apps/riot/internal/mcp"
	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
	"github.com/verikono/glocalize/apps/riot/internal/utils"
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
}

// RiotServer serves the GraphQL and MCP surfaces over one HTTP port, backed by the phi-generated orm client.
type RiotServer struct {
	config RiotConfig
	db     *orm.DB
	graph  *gql.RiotGraph
	mcps   *mcp.RiotMcp
	mux    *http.ServeMux
}

// NewRiotServer opens the store, runs migrations, builds the GraphQL schema and the MCP tool surface, and wires both onto a single mux. It returns a fully-wired server (or an error) — never a half-built one.
func NewRiotServer(options ...T.OptionerFn) (*RiotServer, error) {

	config := RiotConfig{
		Port:   utils.ResolvePort(),
		GqlURL: "/graph",
		McpURL: "/mcp",
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

	inst := &RiotServer{config: config, db: db, graph: graph, mcps: mcps}

	inst.mux = http.NewServeMux()
	inst.mux.Handle(config.GqlURL, graph.Handler())
	inst.mux.Handle(config.McpURL, mcps.Handler())
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
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// WithGqlPath overrides the GraphQL mount point.
func WithGqlPath(value string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.GqlURL = value
			}
		}
		return nil
	}
}

// WithMcpPath overrides the MCP mount point.
func WithMcpPath(value string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.McpURL = value
			}
		}
		return nil
	}
}

// WithPort overrides the serve port.
func WithPort(port int) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.Port = port
			}
		}
		return nil
	}
}

// WithPlaygroundPath mounts the GraphQL playground at the given path when
// non-empty.
func WithPlaygroundPath(value string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*RiotConfig); ok {
				config.Playground = value
			}
		}
		return nil
	}
}

func (c *RiotConfig) applyEnvVars() error {

	if os.Getenv("RIOT_GRAPHL_PLAYGROUND") != "" {
		c.Playground = os.Getenv("RIOT_GRAPHL_PLAYGROUND")
	}

	return nil
}
