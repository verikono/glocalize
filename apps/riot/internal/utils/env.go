package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// DefaultPort is the port riot serves on when RIOT_URL carries no port.
	DefaultPort = 6660

	// DefaultURL is the serve URL used when RIOT_URL is empty.
	DefaultURL = "http://localhost:6660/"
)

// ResolveURL returns RIOT_URL, falling back to http://localhost:6660/.
func ResolveURL() string {
	if raw := strings.TrimSpace(os.Getenv("RIOT_URL")); raw != "" {
		return raw
	}
	return DefaultURL
}

// ResolvePort returns the port encoded in RIOT_URL, falling back to DefaultPort
// when the URL is absent or its port is invalid.
func ResolvePort() int {
	u, err := url.Parse(ResolveURL())
	if err != nil {
		return DefaultPort
	}
	if u.Port() == "" {
		return DefaultPort
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port < 1 || port > 65535 {
		return DefaultPort
	}
	return port
}

// ResolveDataURL returns RIOT_DATAURL (a url whose scheme selects the database
// dialect), falling back to the default on-disk SQLite store.
func ResolveDataURL() string {
	if raw := strings.TrimSpace(os.Getenv("RIOT_DATAURL")); raw != "" {
		return raw
	}
	return "sqlite://" + defaultDataPath()
}

// ResolveDirectDataURL returns RIOT_DIRECT_DATAURL (the raw database path with
// no scheme), falling back to the default on-disk SQLite path.
func ResolveDirectDataURL() string {
	if raw := strings.TrimSpace(os.Getenv("RIOT_DIRECT_DATAURL")); raw != "" {
		return raw
	}
	return defaultDataPath()
}

// ProviderAndDSN derives the sql driver name and data source name from
// RIOT_DATAURL. The scheme selects the dialect: sqlite:///abs/path or
// sqlite://rel/path yields ("sqlite", "/abs/path" | "rel/path"), while
// postgres://… or postgresql://… yields ("postgres", <url>). A leading ~ in a
// sqlite path is expanded to the user's home directory.
func ProviderAndDSN() (provider, dsn string, err error) {
	raw := ResolveDataURL()
	switch {
	case strings.HasPrefix(raw, "sqlite://"):
		path := strings.TrimPrefix(raw, "sqlite://")
		path = expandHome(path)
		return "sqlite", path, nil
	case strings.HasPrefix(raw, "postgres://"), strings.HasPrefix(raw, "postgresql://"):
		return "postgres", raw, nil
	default:
		return "", "", fmt.Errorf("unsupported RIOT_DATAURL scheme in %q", raw)
	}
}

// EnsureSQLiteDir creates the parent directory of an on-disk SQLite path so the
// driver can open (and create) the file. Relative paths are left alone.
func EnsureSQLiteDir(dsn string) error {
	dir := filepath.Dir(dsn)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// expandHome expands a leading ~ to the user's home directory.
func expandHome(path string) string {
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return path
	}
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// defaultDataPath returns the on-disk SQLite path riot uses when no URL is
// configured: <user-config>/riot/data.sqlite.
func defaultDataPath() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "riot", "data.sqlite")
	}
	return "riot/data.sqlite"
}
