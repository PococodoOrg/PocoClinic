package database_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

func TestConnectAndMigrate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pococlinic.db")
	ctx := context.Background()

	db, err := database.Connect(ctx, path)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	applied, err := database.AppliedMigrations(ctx, db)
	if err != nil {
		t.Fatalf("applied: %v", err)
	}
	if len(applied) == 0 {
		t.Fatal("expected at least one applied migration")
	}

	var n int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM form_groups`).Scan(&n); err != nil {
		t.Fatalf("query form_groups: %v", err)
	}
	if n < 1 {
		t.Fatalf("expected default form group, got %d", n)
	}

	pending, err := database.PendingMigrations(ctx, db)
	if err != nil {
		t.Fatalf("pending: %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("expected no pending migrations, got %v", pending)
	}

	if err := db.QuickCheck(ctx); err != nil {
		t.Fatalf("quick_check: %v", err)
	}
}

func TestProductionPragmasApplied(t *testing.T) {
	ctx := context.Background()
	db, err := database.Open(ctx, filepath.Join(t.TempDir(), "pragmas.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	var fk, journal, sync string
	if err := db.QueryRow(ctx, `PRAGMA foreign_keys`).Scan(&fk); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(ctx, `PRAGMA journal_mode`).Scan(&journal); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(ctx, `PRAGMA synchronous`).Scan(&sync); err != nil {
		t.Fatal(err)
	}
	if fk != "1" {
		t.Fatalf("foreign_keys=%q", fk)
	}
	if journal != "wal" {
		t.Fatalf("journal_mode=%q", journal)
	}
	if sync != "1" {
		t.Fatalf("synchronous=%q (want NORMAL=1)", sync)
	}
	var timeout int
	if err := db.QueryRow(ctx, `PRAGMA busy_timeout`).Scan(&timeout); err != nil {
		t.Fatal(err)
	}
	if timeout < 5000 {
		t.Fatalf("busy_timeout=%d", timeout)
	}
}

func TestRejectPostgresURL(t *testing.T) {
	_, err := database.ResolvePath("postgresql://root@localhost:26257/pococlinic")
	if err == nil {
		t.Fatal("expected error for postgres URL")
	}
}
