package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

// DB wraps *sql.DB with Postgres-style $1 placeholders rebound for SQLite.
type DB struct {
	sql  *sql.DB
	Path string
}

// Result matches pgx CommandTag.RowsAffected() (int64, no error).
type Result struct {
	sql.Result
}

func (r Result) RowsAffected() int64 {
	if r.Result == nil {
		return 0
	}
	n, _ := r.Result.RowsAffected()
	return n
}

// Tx wraps *sql.Tx with the same rebind helpers.
type Tx struct {
	tx *sql.Tx
}

// Open opens (or creates) a SQLite database at path.
// path may be a plain filesystem path; "sqlite:" or "file:" prefixes are stripped.
//
// Production settings (applied on every pooled connection via modernc DSN _pragma):
//   - WAL journal for concurrent readers
//   - foreign_keys ON (SQLite defaults to OFF per connection)
//   - busy_timeout 5s (wait instead of immediate SQLITE_BUSY)
//   - synchronous NORMAL (durable enough with WAL; fewer fsyncs than FULL)
//   - temp_store MEMORY
//   - single open connection (avoids write-write lock upgrade deadlocks)
func Open(ctx context.Context, databaseURL string) (*DB, error) {
	path, err := ResolvePath(databaseURL)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	sqlDB, err := sql.Open("sqlite", buildDSN(path))
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	// One connection: SQLite allows one writer; a multi-conn pool can deadlock on
	// read→write lock upgrades even with busy_timeout. Clinic load is fine at 1.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)

	db := &DB{sql: sqlDB, Path: path}
	if err := db.Ping(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	if err := db.verifyProductionPragmas(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

// buildDSN attaches modernc.org/sqlite _pragma / _txlock options so every new
// connection gets the same settings. One-shot Exec("PRAGMA …") only affects the
// first pooled connection and is not production-safe.
func buildDSN(path string) string {
	opts := []string{
		"_pragma=busy_timeout(5000)",
		"_pragma=foreign_keys(ON)",
		"_pragma=journal_mode(WAL)",
		"_pragma=synchronous(NORMAL)",
		"_pragma=temp_store(MEMORY)",
		"_pragma=wal_autocheckpoint(1000)",
		"_txlock=immediate",
	}
	return path + "?" + strings.Join(opts, "&")
}

func (db *DB) verifyProductionPragmas(ctx context.Context) error {
	checks := []struct {
		pragma string
		want   string
	}{
		{"foreign_keys", "1"},
		{"journal_mode", "wal"},
		{"synchronous", "1"}, // NORMAL
	}
	for _, c := range checks {
		var got string
		if err := db.sql.QueryRowContext(ctx, "PRAGMA "+c.pragma).Scan(&got); err != nil {
			return fmt.Errorf("read PRAGMA %s: %w", c.pragma, err)
		}
		if !strings.EqualFold(got, c.want) {
			return fmt.Errorf("PRAGMA %s = %q, want %q (DSN pragmas not applied?)", c.pragma, got, c.want)
		}
	}
	var timeout int
	if err := db.sql.QueryRowContext(ctx, `PRAGMA busy_timeout`).Scan(&timeout); err != nil {
		return fmt.Errorf("read PRAGMA busy_timeout: %w", err)
	}
	if timeout < 5000 {
		return fmt.Errorf("PRAGMA busy_timeout = %d, want >= 5000", timeout)
	}
	return nil
}

// QuickCheck runs PRAGMA quick_check (faster than full integrity_check).
// Returns nil when the database reports "ok".
func (db *DB) QuickCheck(ctx context.Context) error {
	if db == nil || db.sql == nil {
		return fmt.Errorf("database is nil")
	}
	rows, err := db.sql.QueryContext(ctx, `PRAGMA quick_check`)
	if err != nil {
		return fmt.Errorf("quick_check: %w", err)
	}
	defer rows.Close()

	var issues []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return err
		}
		if !strings.EqualFold(line, "ok") {
			issues = append(issues, line)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(issues) > 0 {
		return fmt.Errorf("sqlite quick_check failed: %s", strings.Join(issues, "; "))
	}
	return nil
}

// QuickCheckFile opens a SQLite file read-only, runs quick_check, and closes it.
func QuickCheckFile(ctx context.Context, path string) error {
	path, err := ResolvePath(path)
	if err != nil {
		return err
	}
	dsn := path + "?mode=ro&_pragma=busy_timeout(5000)"
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("open for quick_check: %w", err)
	}
	defer sqlDB.Close()
	sqlDB.SetMaxOpenConns(1)

	db := &DB{sql: sqlDB, Path: path}
	if err := db.Ping(ctx); err != nil {
		return err
	}
	return db.QuickCheck(ctx)
}

// CheckpointTRUNCATE forces a WAL checkpoint (useful after large writes / before copy).
func (db *DB) CheckpointTRUNCATE(ctx context.Context) error {
	_, err := db.sql.ExecContext(ctx, `PRAGMA wal_checkpoint(TRUNCATE)`)
	return err
}

// Connect opens SQLite and applies pending migrations.
func Connect(ctx context.Context, databaseURL string) (*DB, error) {
	db, err := Open(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := RunMigrations(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	return db, nil
}

// OpenPool is an alias for Open (kept for call-site compatibility during the cutover).
func OpenPool(ctx context.Context, databaseURL string) (*DB, error) {
	return Open(ctx, databaseURL)
}

// ResolvePath normalizes DATABASE_URL to a filesystem path.
func ResolvePath(databaseURL string) (string, error) {
	raw := strings.TrimSpace(databaseURL)
	if raw == "" {
		return "", fmt.Errorf("database path is empty")
	}
	lower := strings.ToLower(raw)
	switch {
	case strings.HasPrefix(lower, "postgresql://"), strings.HasPrefix(lower, "postgres://"):
		return "", fmt.Errorf("postgresql URLs are no longer supported; set DATABASE_URL to a SQLite file path (e.g. ./data/pococlinic.db)")
	case strings.HasPrefix(lower, "sqlite://"):
		raw = raw[len("sqlite://"):]
	case strings.HasPrefix(lower, "sqlite:"):
		raw = raw[len("sqlite:"):]
	case strings.HasPrefix(lower, "file:"):
		raw = raw[len("file:"):]
		if i := strings.Index(raw, "?"); i >= 0 {
			raw = raw[:i]
		}
	}
	// Strip accidental query string if callers passed a DSN.
	if i := strings.Index(raw, "?"); i >= 0 {
		raw = raw[:i]
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("database path is empty")
	}
	if !filepath.IsAbs(raw) {
		abs, err := filepath.Abs(raw)
		if err != nil {
			return "", err
		}
		raw = abs
	}
	return raw, nil
}

func (db *DB) Close() error {
	if db == nil || db.sql == nil {
		return nil
	}
	return db.sql.Close()
}

// ReplaceWithSQLiteBytes closes the live DB, replaces the file on disk, and reopens it.
// Callers must ensure no other process (ops-helper vs main EMR) still holds the file open.
func (db *DB) ReplaceWithSQLiteBytes(ctx context.Context, payload []byte) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	path := db.Path
	if path == "" {
		return fmt.Errorf("database path is unknown")
	}

	tmpCheck := path + ".integrity-check"
	if err := os.WriteFile(tmpCheck, payload, 0o640); err != nil {
		return err
	}
	defer os.Remove(tmpCheck)
	if err := QuickCheckFile(ctx, tmpCheck); err != nil {
		return fmt.Errorf("restore rejected: %w", err)
	}

	_ = db.sql.Close()
	_ = os.Remove(path + "-wal")
	_ = os.Remove(path + "-shm")

	tmp := path + ".restore-tmp"
	if err := os.WriteFile(tmp, payload, 0o640); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	reopened, err := Open(ctx, path)
	if err != nil {
		return err
	}
	db.sql = reopened.sql
	db.Path = reopened.Path
	return nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.sql.PingContext(ctx)
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	res, err := db.sql.ExecContext(ctx, rebind(query), args...)
	return Result{res}, err
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.sql.QueryContext(ctx, rebind(query), args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.sql.QueryRowContext(ctx, rebind(query), args...)
}

func (db *DB) Begin(ctx context.Context) (*Tx, error) {
	tx, err := db.sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

func (tx *Tx) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	res, err := tx.tx.ExecContext(ctx, rebind(query), args...)
	return Result{res}, err
}

func (tx *Tx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, rebind(query), args...)
}

func (tx *Tx) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return tx.tx.QueryRowContext(ctx, rebind(query), args...)
}

func (tx *Tx) Commit() error   { return tx.tx.Commit() }
func (tx *Tx) Rollback() error { return tx.tx.Rollback() }

var placeholderRE = regexp.MustCompile(`\$(\d+)`)

// rebind converts Postgres-style $1,$2 placeholders to SQLite numbered ?1,?2.
// Numbered form preserves reuse ($1 appearing thrice stays ?1 thrice). Anonymous
// ? would consume a new argument each time and break filters that reuse $n.
func rebind(query string) string {
	if !strings.Contains(query, "$") {
		return query
	}
	return placeholderRE.ReplaceAllString(query, "?$1")
}

// Row is satisfied by *sql.Row and *sql.Rows for Scan helpers.
type Row interface {
	Scan(dest ...any) error
}

// ErrNoRows is exported for callers migrating off pgx.
var ErrNoRows = sql.ErrNoRows
