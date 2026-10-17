package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	G "github.com/verikono/glocalize/packages/go/data/prisma/graphql"
	P "github.com/verikono/glocalize/packages/go/puck"
)

// GraphFromSchema reads every *.prisma file under internal/schema, converts the
// combined document into a GraphQL schema, and writes it to internal/gql/
// riot.graphql. It returns the generated schema text.
func GraphFromSchema() (string, error) {
	schemaGlob := filepath.Join(GetRiotSchemaPath(), "*.prisma")
	sourceSchema, err := P.From(schemaGlob).Read()
	if err != nil {
		return "", fmt.Errorf("GraphFromSchema: read %s: %w", schemaGlob, err)
	}
	targetSchema, err := G.ParsePrisma(sourceSchema)
	if err != nil {
		return "", fmt.Errorf("GraphFromSchema: parse prisma: %w", err)
	}
	if targetSchema == nil || *targetSchema == "" {
		return "", fmt.Errorf("GraphFromSchema: empty schema")
	}
	if err := WriteRiotGraph(*targetSchema); err != nil {
		return "", err
	}
	return *targetSchema, nil
}

// WriteRiotGraph writes the generated GraphQL schema to internal/gql/riot.graphql.
func WriteRiotGraph(schema string) error {
	targetPath := filepath.Join(GetRiotGraphPath(), "riot.graphql")
	if err := os.WriteFile(targetPath, []byte(schema), 0o644); err != nil {
		return fmt.Errorf("GraphFromSchema: write %s: %w", targetPath, err)
	}
	return nil
}


func GetRiotSchemaPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "schema")
}

func GetRiotGraphPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "gql")
}

func GetRiotMcpPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "mcp")
}