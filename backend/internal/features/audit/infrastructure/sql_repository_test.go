package infrastructure

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
	"github.com/google/uuid"
)

func setupAuditSQL(t *testing.T) (*SQLRepository, context.Context) {
	t.Helper()
	ctx := context.Background()
	db, err := database.Connect(ctx, filepath.Join(t.TempDir(), "audit.db"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewSQLRepository(db), ctx
}

func TestSQLRepository_CreateAndListForUser(t *testing.T) {
	repo, ctx := setupAuditSQL(t)
	userID := uuid.New()
	targetID := uuid.New()

	actorEvent := &domain.Event{
		EventType: domain.EventLoginSuccess,
		UserID:    &userID,
		Success:   true,
		Details:   map[string]string{"method": "staff"},
		CreatedAt: time.Now().UTC(),
	}
	if err := repo.Create(ctx, actorEvent); err != nil {
		t.Fatalf("create actor event: %v", err)
	}

	targetEvent := &domain.Event{
		EventType:    domain.EventUserUpdated,
		UserID:       &targetID,
		ResourceType: "user",
		ResourceID:   userID.String(),
		Success:      true,
		CreatedAt:    time.Now().UTC(),
	}
	if err := repo.Create(ctx, targetEvent); err != nil {
		t.Fatalf("create target event: %v", err)
	}

	events, total, err := repo.ListForUser(ctx, userID.String(), 1, 10)
	if err != nil {
		t.Fatalf("list for user: %v", err)
	}
	if total < 2 {
		t.Fatalf("expected at least 2 events for user, got total=%d events=%d", total, len(events))
	}
}

func TestSQLRepository_ListRecentFiltersByEventType(t *testing.T) {
	repo, ctx := setupAuditSQL(t)
	userID := uuid.New()

	for _, eventType := range []string{domain.EventLoginSuccess, domain.EventLogout, domain.EventLoginSuccess} {
		uid := userID
		if err := repo.Create(ctx, &domain.Event{
			EventType: eventType,
			UserID:    &uid,
			Success:   true,
			CreatedAt: time.Now().UTC(),
		}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	events, total, err := repo.ListRecent(ctx, domain.ListFilter{
		EventType: domain.EventLoginSuccess,
		Page:      1,
		PageSize:  10,
	})
	if err != nil {
		t.Fatalf("list recent: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 login events, got total=%d", total)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}
