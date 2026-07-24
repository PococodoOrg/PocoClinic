package infrastructure

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/google/uuid"
)

// MemoryRepository stores audit events in memory for development without a database.
type MemoryRepository struct {
	events []*domain.Event
}

// NewMemoryRepository creates an in-memory audit repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{events: make([]*domain.Event, 0)}
}

// Create appends an audit event.
func (r *MemoryRepository) Create(ctx context.Context, event *domain.Event) error {
	r.events = append(r.events, event)
	return nil
}

// ListForUser returns audit events performed by or targeting a user.
func (r *MemoryRepository) ListForUser(ctx context.Context, userID string, page, pageSize int) ([]*domain.Event, int64, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user id")
	}

	matched := make([]*domain.Event, 0)
	for _, event := range r.events {
		if event.UserID != nil && *event.UserID == parsedID {
			matched = append(matched, event)
			continue
		}
		if event.ResourceType == "user" && event.ResourceID == userID {
			matched = append(matched, event)
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})

	totalCount := int64(len(matched))
	start := (page - 1) * pageSize
	if start >= len(matched) {
		return []*domain.Event{}, totalCount, nil
	}
	end := start + pageSize
	if end > len(matched) {
		end = len(matched)
	}

	return matched[start:end], totalCount, nil
}

// ListRecent returns audit events across the system with optional filters.
func (r *MemoryRepository) ListRecent(ctx context.Context, filter domain.ListFilter) ([]*domain.Event, int64, error) {
	_ = ctx
	var since time.Time
	if filter.SinceDays > 0 {
		since = time.Now().AddDate(0, 0, -filter.SinceDays)
	}
	matched := make([]*domain.Event, 0, len(r.events))
	for _, event := range r.events {
		if filter.EventType != "" && event.EventType != filter.EventType {
			continue
		}
		if filter.UserID != "" {
			if event.UserID == nil || event.UserID.String() != filter.UserID {
				continue
			}
		}
		if !since.IsZero() && event.CreatedAt.Before(since) {
			continue
		}
		matched = append(matched, event)
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})

	totalCount := int64(len(matched))
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 25
	}

	start := (page - 1) * pageSize
	if start >= len(matched) {
		return []*domain.Event{}, totalCount, nil
	}
	end := start + pageSize
	if end > len(matched) {
		end = len(matched)
	}
	return matched[start:end], totalCount, nil
}

// Events returns a copy of stored events (for tests).
func (r *MemoryRepository) Events() []*domain.Event {
	copied := make([]*domain.Event, len(r.events))
	copy(copied, r.events)
	return copied
}

// Len returns the number of stored events.
func (r *MemoryRepository) Len() int {
	return len(r.events)
}

var _ domain.Repository = (*MemoryRepository)(nil)
