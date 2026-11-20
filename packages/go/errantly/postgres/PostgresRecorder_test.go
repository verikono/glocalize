package PostgresRecorder_test

import (
	"os"
	"testing"

	"github.com/verikono/glocalize/packages/go/errantly"
	PostgresRecorder "github.com/verikono/glocalize/packages/go/errantly/postgres"
)

// TestRecordWritesToErrorRecord proves the recorder persists a real row. Gated on PGSQL_URL: a plain `go test ./...` must never touch a database.
func TestRecordWritesToErrorRecord(t *testing.T) {
	if os.Getenv("PGSQL_URL") == "" {
		t.Skip("PGSQL_URL unset — skipping the recorder write test")
	}

	recorder, err := PostgresRecorder.NewPostGresHandler(PostgresRecorder.WithSystem("errantly-test"))
	if err != nil {
		t.Fatalf("NewPostGresHandler() error = %v", err)
	}
	if recorder == nil {
		t.Fatal("expected a recorder when PGSQL_URL is set")
	}

	if err := recorder.Record("errantly recorder test message", errantly.ErrantlyMessageMeta{Message: "errantly recorder test message", Level: "test"}); err != nil {
		t.Fatalf("Record() error = %v", err)
	}
}

// TestBlankConnectionStringIsNil proves recording is opt-in: no DSN, no recorder, and no error.
func TestBlankConnectionStringIsNil(t *testing.T) {
	t.Setenv("PGSQL_URL", "")

	recorder, err := PostgresRecorder.NewPostGresHandler()
	if err != nil {
		t.Fatalf("NewPostGresHandler() error = %v", err)
	}
	if recorder != nil {
		t.Fatal("expected a nil recorder when PGSQL_URL is blank")
	}
}

// TestRecordOnNilRecorderIsNoop proves a nil recorder is safe to call, so a service with no database needs no nil check at the call site.
func TestRecordOnNilRecorderIsNoop(t *testing.T) {
	var recorder *PostgresRecorder.PostgresRecorder
	if err := recorder.Record("ignored", errantly.ErrantlyMessageMeta{}); err != nil {
		t.Fatalf("Record() on a nil recorder error = %v", err)
	}
}
