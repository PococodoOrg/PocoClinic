package service

import (
	"context"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/google/uuid"
)

// AuditService records audit events synchronously for durability.
type AuditService struct {
	repo   domain.Repository
	logger *logging.Logger
}

// NewAuditService creates an audit logger service.
func NewAuditService(repo domain.Repository, logger *logging.Logger) *AuditService {
	return &AuditService{repo: repo, logger: logger}
}

// Log persists an audit event before returning to the caller.
func (s *AuditService) Log(ctx context.Context, event domain.Event) {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	if err := s.repo.Create(ctx, &event); err != nil {
		s.logger.Error("Failed to write audit log", err,
			"eventType", event.EventType,
			"resourceType", event.ResourceType,
			"resourceId", event.ResourceID,
		)
	}
}
