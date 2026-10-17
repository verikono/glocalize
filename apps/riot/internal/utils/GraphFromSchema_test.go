package utils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

func TestGraphFromSchema(t *testing.T) {
	schema, err := utils.GraphFromSchema()
	if err != nil {
		t.Fatalf("GraphFromSchema() error = %v", err)
	}
	if strings.TrimSpace(schema) == "" {
		t.Fatal("GraphFromSchema() returned an empty schema")
	}

	target := filepath.Join(utils.GetRiotGraphPath(), "riot.graphql")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read %s: %v", target, err)
	}
	if string(data) != schema {
		t.Fatalf("riot.graphql does not match the returned schema")
	}
}
