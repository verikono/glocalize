package gql

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
	T "github.com/verikono/glocalize/types/go"
)

// RiotGraph is the read-only GraphQL surface over the riot orm client. It
// builds an executable schema from the generated GraphQL document and serves it
// over HTTP.
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

// Handler returns the GraphQL HTTP handler (query execution only, no
// playground — Altair is an external client).
func (g *RiotGraph) Handler() http.Handler {
	return handler.New(&handler.Config{
		Schema:   &g.schema,
		Pretty:   true,
		GraphiQL: false,
	})
}

// PlaygroundHandler returns the bundled GraphiQL playground.
func (g *RiotGraph) PlaygroundHandler() http.Handler {
	return handler.New(&handler.Config{
		Schema:   &g.schema,
		Pretty:   true,
		GraphiQL: true,
	})
}

// Schema exposes the executable schema (for tests and introspection).
func (g *RiotGraph) Schema() graphql.Schema {
	return g.schema
}
