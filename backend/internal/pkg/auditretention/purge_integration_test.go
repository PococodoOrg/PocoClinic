package auditretention

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/google/uuid"
)

func TestPurgeOlderThanDeletesOldRows(t *testing.T) {
	ctx := context.Background()
	db, err := database.Connect(ctx, filepath.Join(t.TempDir(), "audit-purge.db"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	oldTime := time.Now().UTC().AddDate(0, 0, -60)
	newTime := time.Now().UTC().AddDate(0, 0, -1)

	for _, createdAt := range []time.Time{oldTime, oldTime, newTime} {
		_, err := db.Exec(ctx, `
			INSERT INTO audit_logs (id, event_type, success, created_at)
			VALUES ($1, $2, $3, $4)
		`, uuid.New(), "test.event", true, createdAt)
		if err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	before, err := CountRows(ctx, db)
	if err != nil {
		t.Fatalf("count before: %v", err)
	}
	if before != 3 {
		t.Fatalf("expected 3 rows before purge, got %d", before)
	}

	result, err := PurgeOlderThan(ctx, db, 30)
	if err != nil {
		t.Fatalf("purge: %v", err)
	}
	if result.DeletedRows != 2 {
		t.Fatalf("expected 2 deleted rows, got %d", result.DeletedRows)
	}

	after, err := CountRows(ctx, db)
	if err != nil {
		t.Fatalf("count after: %v", err)
	}
	if after != 1 {
		t.Fatalf("expected 1 row after purge, got %d", after)
	}
}
