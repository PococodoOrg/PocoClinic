package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/google/uuid"
)

func TestMemoryRepository_ListForUserIncludesActorAndTarget(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()
	userID := uuid.New()
	targetID := uuid.New()

	if err := repo.Create(ctx, &domain.Event{
		EventType: domain.EventLoginSuccess,
		UserID:    &userID,
		Success:   true,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("create actor event: %v", err)
	}
	if err := repo.Create(ctx, &domain.Event{
		EventType:    domain.EventUserUpdated,
		UserID:       &targetID,
		ResourceType: "user",
		ResourceID:   userID.String(),
		Success:      true,
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		t.Fatalf("create target event: %v", err)
	}

	events, total, err := repo.ListForUser(ctx, userID.String(), 1, 10)
	if err != nil {
		t.Fatalf("list for user: %v", err)
	}
	if total != 2 || len(events) != 2 {
		t.Fatalf("expected 2 events, got total=%d len=%d", total, len(events))
	}
}

func TestMemoryRepository_ListRecentFiltersEventType(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()
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
	if total != 2 || len(events) != 2 {
		t.Fatalf("expected 2 login events, got total=%d len=%d", total, len(events))
	}
}
