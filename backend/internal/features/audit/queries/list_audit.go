package queries

import (
	"context"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/pagination"
)

// ListAuditQuery retrieves recent audit events across the system.
type ListAuditQuery struct {
	Page      int
	PageSize  int
	EventType string
	UserID    string
	SinceDays int
}

// ListAuditHandler handles system-wide audit retrieval.
type ListAuditHandler interface {
	Handle(ctx context.Context, query ListAuditQuery) (*PaginatedAuditEntries, error)
}

type listAuditHandler struct {
	auditRepository auditdomain.Repository
}

func NewListAuditHandler(repo auditdomain.Repository) ListAuditHandler {
	return &listAuditHandler{auditRepository: repo}
}

func (h *listAuditHandler) Handle(ctx context.Context, query ListAuditQuery) (*PaginatedAuditEntries, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 25
	}
	query.PageSize = pagination.Clamp(query.PageSize, 25)

	events, totalCount, err := h.auditRepository.ListRecent(ctx, auditdomain.ListFilter{
		Page:      query.Page,
		PageSize:  query.PageSize,
		EventType: query.EventType,
		UserID:    query.UserID,
		SinceDays: query.SinceDays,
	})
	if err != nil {
		return nil, err
	}

	result := Paginate(totalCount, query.Page, query.PageSize)
	result.Events = MapEventsToEntries(events)
	return &result, nil
}
