package PostgresRecorder

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/verikono/glocalize/packages/go/errantly"
)

const (
	defaultTable  = "ErrorRecord"
	retryInterval = 15 * time.Second
	recordTimeout = 5 * time.Second
)

// insertStatement writes one row into ErrorRecord's columns. system and code are nullable, so an unset option lands as NULL rather than an empty string. @see .agents.md "### PostgresRecorder"
const insertStatement = `INSERT INTO %s ("message", "created", "system", "code") VALUES ($1, $2, $3, $4)`

type PostgresRecorder struct {
	ConnectionString string
	Table            string
	System           string

	mu sync.RWMutex
	db *sql.DB
}

func NewPostGresHandler(opts ...OptionerFn) (*PostgresRecorder, error) {
	inst := &PostgresRecorder{
		ConnectionString: strings.TrimSpace(os.Getenv("PGSQL_URL")),
		Table:            defaultTable,
	}

	for _, option := range opts {
		if err := option(inst); err != nil {
			return nil, err
		}
	}

	if inst.ConnectionString == "" {
		return nil, nil
	}
	if inst.Table != defaultTable {
		return nil, fmt.Errorf("PostgresRecorder: table %q is not supported; the insert is written against %q's columns", inst.Table, defaultTable)
	}

	inst.connect()
	return inst, nil
}

// Close releases the connection pool. The recorder owns the handle, so nothing else can close it. @see .agents.md "### PostgresRecorder"
func (p *PostgresRecorder) Close() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	db := p.db
	p.db = nil
	p.mu.Unlock()
	if db == nil {
		return nil
	}
	return db.Close()
}

// Record writes one error to ErrorRecord. A nil recorder or an unconnected recorder is a no-op. @see .agents.md "### PostgresRecorder"
func (p *PostgresRecorder) Record(message string, meta errantly.ErrantlyMessageMeta) error {
	if p == nil {
		return nil
	}
	db := p.getDB()
	if db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), recordTimeout)
	defer cancel()

	query := fmt.Sprintf(insertStatement, pq.QuoteIdentifier(p.Table))
	if _, err := db.ExecContext(ctx, query, message, time.Now().UTC(), nullable(p.System), nullable(meta.Level)); err != nil {
		return fmt.Errorf("PostgresRecorder: failed recording error: %w", err)
	}
	return nil
}

// connect dials once and falls back to the retry loop when the database does not answer. @see .agents.md "### PostgresRecorder"
func (p *PostgresRecorder) connect() {
	db, err := open(p.ConnectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PostgresRecorder: cannot reach database: %v\n", err)
		go p.retry()
		return
	}
	p.setDB(db)
}

func (p *PostgresRecorder) getDB() *sql.DB {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db
}

// nullable maps a blank string to SQL NULL so an unset column stays absent.
func nullable(value string) any {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		return trimmed
	}
	return nil
}

// open dials the database and proves the connection with a ping, because sql.Open itself never contacts the server. @see .agents.md "### PostgresRecorder"
func open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", withDisableSsl(dsn))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}

// retry dials until the database answers, then hands the handle to the recorder. @see .agents.md "### PostgresRecorder"
func (p *PostgresRecorder) retry() {
	for {
		time.Sleep(retryInterval)
		db, err := open(p.ConnectionString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "PostgresRecorder: cannot reach database: %v\n", err)
			continue
		}
		p.setDB(db)
		return
	}
}

// setDB swaps the handle and closes the one it replaces, so a retry cannot leak a pool. @see .agents.md "### PostgresRecorder"
func (p *PostgresRecorder) setDB(db *sql.DB) {
	p.mu.Lock()
	previous := p.db
	p.db = db
	p.mu.Unlock()
	if previous != nil && previous != db {
		_ = previous.Close()
	}
}

// withDisableSsl appends sslmode=disable when the DSN names no sslmode. lib/pq defaults to require, which an in-cluster server refuses. @see .agents.md "### PostgresRecorder"
func withDisableSsl(dsn string) string {
	if strings.Contains(dsn, "sslmode=") {
		return dsn
	}
	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}
	return dsn + separator + "sslmode=disable"
}
