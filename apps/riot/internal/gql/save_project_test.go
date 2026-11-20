package gql

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/graphql-go/graphql"
	_ "modernc.org/sqlite"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// newTestSchema opens a throwaway sqlite store, migrates it and returns the executable schema wired to it.
func newTestSchema(t *testing.T) graphql.Schema {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "riot-gql-test.sqlite")
	db, err := orm.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.RunMigrations(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	schema, err := buildSchema(db)
	if err != nil {
		t.Fatalf("buildSchema: %v", err)
	}
	return schema
}

// run executes one operation against the schema and returns the result.
func run(t *testing.T, schema graphql.Schema, query string, variables map[string]interface{}) *graphql.Result {
	t.Helper()
	return graphql.Do(graphql.Params{Schema: schema, RequestString: query, VariableValues: variables})
}

const saveProjectMutation = `mutation SaveProject($input: ProjectInput!) { SaveProject(input: $input) { name location skillsPath agentsPath } }`

const projectsQuery = `query Projects { Projects { name location skillsPath agentsPath } }`

// TestSaveProjectCreate proves a create stores all four fields and the record is visible to the Projects query.
func TestSaveProjectCreate(t *testing.T) {
	schema := newTestSchema(t)

	result := run(t, schema, saveProjectMutation, map[string]interface{}{
		"input": map[string]interface{}{"name": "alpha", "location": "/srv/alpha", "skillsPath": "skills", "agentsPath": "agents", "isNew": true},
	})
	if len(result.Errors) > 0 {
		t.Fatalf("SaveProject create: %v", result.Errors)
	}

	saved, ok := result.Data.(map[string]interface{})["SaveProject"].(map[string]interface{})
	if !ok {
		t.Fatalf("SaveProject returned %T", result.Data.(map[string]interface{})["SaveProject"])
	}
	for field, want := range map[string]string{"name": "alpha", "location": "/srv/alpha", "skillsPath": "skills", "agentsPath": "agents"} {
		if got := saved[field]; got != want {
			t.Errorf("SaveProject.%s = %v, want %q", field, got, want)
		}
	}

	listed := run(t, schema, projectsQuery, nil)
	if len(listed.Errors) > 0 {
		t.Fatalf("Projects: %v", listed.Errors)
	}
	rows, ok := listed.Data.(map[string]interface{})["Projects"].([]interface{})
	if !ok || len(rows) != 1 {
		t.Fatalf("Projects returned %v, want one row", listed.Data.(map[string]interface{})["Projects"])
	}
	if name := rows[0].(map[string]interface{})["name"]; name != "alpha" {
		t.Errorf("Projects[0].name = %v, want alpha", name)
	}
}

// TestSaveProjectUpdate proves an update rewrites the three mutable fields and leaves exactly one row with the name.
func TestSaveProjectUpdate(t *testing.T) {
	schema := newTestSchema(t)

	create := run(t, schema, saveProjectMutation, map[string]interface{}{
		"input": map[string]interface{}{"name": "alpha", "location": "/srv/alpha", "skillsPath": "skills", "agentsPath": "agents", "isNew": true},
	})
	if len(create.Errors) > 0 {
		t.Fatalf("SaveProject create: %v", create.Errors)
	}

	update := run(t, schema, saveProjectMutation, map[string]interface{}{
		"input": map[string]interface{}{"name": "alpha", "location": "/srv/moved", "skillsPath": "skills2", "agentsPath": "agents2"},
	})
	if len(update.Errors) > 0 {
		t.Fatalf("SaveProject update: %v", update.Errors)
	}

	saved, ok := update.Data.(map[string]interface{})["SaveProject"].(map[string]interface{})
	if !ok {
		t.Fatalf("SaveProject returned %T", update.Data.(map[string]interface{})["SaveProject"])
	}
	for field, want := range map[string]string{"name": "alpha", "location": "/srv/moved", "skillsPath": "skills2", "agentsPath": "agents2"} {
		if got := saved[field]; got != want {
			t.Errorf("SaveProject.%s = %v, want %q", field, got, want)
		}
	}

	listed := run(t, schema, projectsQuery, nil)
	if len(listed.Errors) > 0 {
		t.Fatalf("Projects: %v", listed.Errors)
	}
	rows, ok := listed.Data.(map[string]interface{})["Projects"].([]interface{})
	if !ok || len(rows) != 1 {
		t.Fatalf("Projects returned %v, want exactly one row", listed.Data.(map[string]interface{})["Projects"])
	}
}

// TestSaveProjectDuplicateCreate proves a create whose name already exists surfaces an error.
func TestSaveProjectDuplicateCreate(t *testing.T) {
	schema := newTestSchema(t)

	input := map[string]interface{}{"name": "alpha", "location": "/srv/alpha", "skillsPath": "skills", "agentsPath": "agents", "isNew": true}
	if result := run(t, schema, saveProjectMutation, map[string]interface{}{"input": input}); len(result.Errors) > 0 {
		t.Fatalf("SaveProject create: %v", result.Errors)
	}

	duplicate := run(t, schema, saveProjectMutation, map[string]interface{}{"input": input})
	if len(duplicate.Errors) == 0 {
		t.Fatalf("SaveProject duplicate create returned no error")
	}
}

// TestSaveProjectEmptyName proves an empty name is rejected.
func TestSaveProjectEmptyName(t *testing.T) {
	schema := newTestSchema(t)

	result := run(t, schema, saveProjectMutation, map[string]interface{}{
		"input": map[string]interface{}{"name": "", "location": "/srv/alpha", "skillsPath": "skills", "agentsPath": "agents", "isNew": true},
	})
	if len(result.Errors) == 0 {
		t.Fatalf("SaveProject empty name returned no error")
	}
}
