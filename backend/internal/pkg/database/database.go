package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

const schemaMigrationsDDL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	name TEXT PRIMARY KEY,
	applied_at TEXT NOT NULL DEFAULT (datetime('now'))
)
`

// RunMigrations applies pending SQL migrations tracked in schema_migrations.
func RunMigrations(ctx context.Context, db *DB) error {
	if _, err := db.Exec(ctx, schemaMigrationsDDL); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		applied, err := isMigrationApplied(ctx, db, entry.Name())
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return err
		}

		if err := applyMigration(ctx, db, entry.Name(), string(content)); err != nil {
			return fmt.Errorf("migration %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func isMigrationApplied(ctx context.Context, db *DB, name string) (bool, error) {
	var exists int
	err := db.QueryRow(ctx, `
		SELECT COUNT(1) FROM schema_migrations WHERE name = $1
	`, name).Scan(&exists)
	return exists > 0, err
}

func applyMigration(ctx context.Context, db *DB, name, sqlText string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range splitSQLStatements(sqlText) {
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", truncateSQL(stmt), err)
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO schema_migrations (name) VALUES ($1)
	`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func splitSQLStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		stmt := strings.TrimSpace(p)
		if stmt == "" {
			continue
		}
		// Skip pure comment blocks
		lines := strings.Split(stmt, "\n")
		nonComment := false
		for _, line := range lines {
			trim := strings.TrimSpace(line)
			if trim != "" && !strings.HasPrefix(trim, "--") {
				nonComment = true
				break
			}
		}
		if !nonComment {
			continue
		}
		out = append(out, stmt)
	}
	return out
}

func truncateSQL(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 80 {
		return s[:80] + "…"
	}
	return s
}

// PendingMigrations returns migration file names not yet applied.
func PendingMigrations(ctx context.Context, db *DB) ([]string, error) {
	if _, err := db.Exec(ctx, schemaMigrationsDDL); err != nil {
		return nil, err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	pending := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		applied, err := isMigrationApplied(ctx, db, entry.Name())
		if err != nil {
			return nil, err
		}
		if !applied {
			pending = append(pending, entry.Name())
		}
	}
	return pending, nil
}

// Ping verifies the database connection is alive.
func Ping(ctx context.Context, db *DB) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	return db.Ping(ctx)
}

// AppliedMigrations returns migration names already applied, sorted.
func AppliedMigrations(ctx context.Context, db *DB) ([]string, error) {
	if _, err := db.Exec(ctx, schemaMigrationsDDL); err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, `SELECT name FROM schema_migrations ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied = append(applied, name)
	}
	return applied, rows.Err()
}
