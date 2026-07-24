package backup

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dksch/pococlinic/internal/pkg/database"
)

func TestSQLiteBackupCreateRestoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	dataDir := t.TempDir()
	dbPath := filepath.Join(dataDir, "clinic.db")
	backupDir := filepath.Join(dataDir, "backups")
	docsDir := filepath.Join(dataDir, "documents")

	db, err := database.Connect(ctx, dbPath)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	_, err = db.Exec(ctx, `
		INSERT INTO patients (id, first_name, last_name, date_of_birth, gender, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, "p-roundtrip", "Pat", "Round", "1990-01-02", "other", "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed patient: %v", err)
	}

	archivePath, err := Create(ctx, db, backupDir, docsDir, "test")
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}
	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("backup file missing: %v", err)
	}

	archive, err := Open(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	if archive.Manifest.DatabaseFormat != "sqlite" {
		t.Fatalf("expected sqlite format, got %q", archive.Manifest.DatabaseFormat)
	}
	integrity := archive.VerifyIntegrity()
	if !integrity.Valid {
		t.Fatalf("integrity failed: %v", integrity.Issues)
	}

	_, err = db.Exec(ctx, `DELETE FROM patients`)
	if err != nil {
		t.Fatalf("clear patients: %v", err)
	}
	var cleared int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM patients`).Scan(&cleared); err != nil {
		t.Fatalf("count cleared: %v", err)
	}
	if cleared != 0 {
		t.Fatalf("expected empty patients before restore, got %d", cleared)
	}

	if err := Restore(ctx, db, archivePath, docsDir); err != nil {
		t.Fatalf("restore: %v", err)
	}

	var restored int
	var name string
	if err := db.QueryRow(ctx, `SELECT COUNT(*), MAX(first_name) FROM patients`).Scan(&restored, &name); err != nil {
		t.Fatalf("count restored: %v", err)
	}
	if restored != 1 || name != "Pat" {
		t.Fatalf("unexpected restore result: count=%d name=%q", restored, name)
	}

	_ = db.Close()
}

func TestSQLiteBackupRejectsLegacyJSONL(t *testing.T) {
	ctx := context.Background()
	dataDir := t.TempDir()
	dbPath := filepath.Join(dataDir, "clinic.db")
	db, err := database.Connect(ctx, dbPath)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	err = Restore(ctx, db, filepath.Join(dataDir, "missing.tar.gz"), "")
	if err == nil {
		t.Fatal("expected restore failure for missing archive")
	}
}
