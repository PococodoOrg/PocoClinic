package backup

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dksch/pococlinic/internal/pkg/database"
)

func TestVerifyIntegritySQLitePayload(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "source.db")
	db, err := database.Connect(ctx, dbPath)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	_ = db.Close()

	payload, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	sum := sha256Hex(payload)
	archive := &Archive{
		Manifest: Manifest{
			DatabaseFormat: "sqlite",
			Checksums: map[string]string{
				sqlitePayloadName: sum,
			},
			Summary: &ManifestSummary{
				TableRows:     map[string]int{"patients": 0},
				DocumentFiles: 0,
			},
		},
		Files: map[string][]byte{
			sqlitePayloadName: payload,
		},
	}
	result := archive.VerifyIntegrity()
	if !result.Valid {
		t.Fatalf("expected valid sqlite archive, issues: %v", result.Issues)
	}
}

func TestVerifyIntegrityRejectsCorruptSQLite(t *testing.T) {
	payload := []byte("SQLite format 3\x00not-a-real-database")
	archive := &Archive{
		Manifest: Manifest{
			DatabaseFormat: "sqlite",
			Checksums: map[string]string{
				sqlitePayloadName: sha256Hex(payload),
			},
		},
		Files: map[string][]byte{
			sqlitePayloadName: payload,
		},
	}
	result := archive.VerifyIntegrity()
	if result.Valid {
		t.Fatal("expected integrity failure for corrupt sqlite payload")
	}
}

func TestVerifyIntegrityMissingSQLitePayload(t *testing.T) {
	archive := &Archive{
		Manifest: Manifest{
			DatabaseFormat: "sqlite",
			Checksums:       map[string]string{},
		},
		Files: map[string][]byte{},
	}
	result := archive.VerifyIntegrity()
	if result.Valid {
		t.Fatal("expected integrity failure without sqlite payload")
	}
}
