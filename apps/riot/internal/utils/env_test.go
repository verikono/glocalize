package utils_test

import (
	"testing"

	"github.com/verikono/glocalize/apps/riot/internal/utils"
)

func TestProviderAndDSN_SQLite(t *testing.T) {
	t.Setenv("RIOT_DATAURL", "sqlite:///tmp/riot-test.sqlite")

	provider, dsn, err := utils.ProviderAndDSN()
	if err != nil {
		t.Fatalf("ProviderAndDSN() error = %v", err)
	}
	if provider != "sqlite" {
		t.Fatalf("provider = %q, want sqlite", provider)
	}
	if dsn != "/tmp/riot-test.sqlite" {
		t.Fatalf("dsn = %q, want /tmp/riot-test.sqlite", dsn)
	}
}

func TestProviderAndDSN_RelativeSQLite(t *testing.T) {
	t.Setenv("RIOT_DATAURL", "sqlite://../../some/dir/db.sqlite")

	_, dsn, err := utils.ProviderAndDSN()
	if err != nil {
		t.Fatalf("ProviderAndDSN() error = %v", err)
	}
	if dsn != "../../some/dir/db.sqlite" {
		t.Fatalf("dsn = %q, want ../../some/dir/db.sqlite", dsn)
	}
}

func TestProviderAndDSN_Postgres(t *testing.T) {
	raw := "postgres://user:pass@host:5432/db"
	t.Setenv("RIOT_DATAURL", raw)

	provider, dsn, err := utils.ProviderAndDSN()
	if err != nil {
		t.Fatalf("ProviderAndDSN() error = %v", err)
	}
	if provider != "postgres" {
		t.Fatalf("provider = %q, want postgres", provider)
	}
	if dsn != raw {
		t.Fatalf("dsn = %q, want %q", dsn, raw)
	}
}

func TestProviderAndDSN_Unsupported(t *testing.T) {
	t.Setenv("RIOT_DATAURL", "mysql://x")

	if _, _, err := utils.ProviderAndDSN(); err == nil {
		t.Fatal("ProviderAndDSN() expected an error for an unsupported scheme")
	}
}

func TestResolvePort_Default(t *testing.T) {
	t.Setenv("RIOT_URL", "")
	if got := utils.ResolvePort(); got != utils.DefaultPort {
		t.Fatalf("ResolvePort() = %d, want %d", got, utils.DefaultPort)
	}
}

func TestResolvePort_FromURL(t *testing.T) {
	t.Setenv("RIOT_URL", "http://localhost:6661/")
	if got := utils.ResolvePort(); got != 6661 {
		t.Fatalf("ResolvePort() = %d, want 6661", got)
	}
}
