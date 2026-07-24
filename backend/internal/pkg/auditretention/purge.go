package auditretention

import (
	"context"
	"fmt"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

// PurgeResult describes deleted audit rows.
type PurgeResult struct {
	RetentionDays int
	DeletedRows   int64
	Cutoff        time.Time
}

// PurgeOlderThan deletes audit log entries older than retentionDays.
func PurgeOlderThan(ctx context.Context, pool *database.DB, retentionDays int) (*PurgeResult, error) {
	if pool == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	if retentionDays <= 0 {
		return nil, fmt.Errorf("retention days must be positive")
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays)
	tag, err := pool.Exec(ctx, `DELETE FROM audit_logs WHERE created_at < $1`, cutoff)
	if err != nil {
		return nil, err
	}

	return &PurgeResult{
		RetentionDays: retentionDays,
		DeletedRows:   tag.RowsAffected(),
		Cutoff:        cutoff,
	}, nil
}

// CountRows returns total audit log rows (for reporting before purge).
func CountRows(ctx context.Context, pool *database.DB) (int64, error) {
	var count int64
	err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs`).Scan(&count)
	return count, err
}
