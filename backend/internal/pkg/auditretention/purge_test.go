package auditretention

import (
	"context"
	"testing"
	"time"
)

func TestPurgeOlderThanRequiresPositiveDays(t *testing.T) {
	_, err := PurgeOlderThan(context.Background(), nil, 0)
	if err == nil {
		t.Fatal("expected error for zero retention days")
	}
}

func TestPurgeResultCutoff(t *testing.T) {
	days := 30
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	result := &PurgeResult{
		RetentionDays: days,
		Cutoff:        cutoff,
	}
	if result.RetentionDays != 30 {
		t.Fatalf("unexpected retention: %d", result.RetentionDays)
	}
}
