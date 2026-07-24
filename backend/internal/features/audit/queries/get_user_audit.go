package queries

import (
	"context"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/pagination"
)

// GetUserAuditQuery retrieves audit history for a staff member.
type GetUserAuditQuery struct {
	UserID   string `json:"-"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=20"`
}

// GetUserAuditHandler handles user audit retrieval.
type GetUserAuditHandler interface {
	Handle(ctx context.Context, query GetUserAuditQuery) (*PaginatedAuditEntries, error)
}

type getUserAuditHandler struct {
	auditRepository auditdomain.Repository
}

// NewGetUserAuditHandler creates a handler for user audit history.
func NewGetUserAuditHandler(repo auditdomain.Repository) GetUserAuditHandler {
	return &getUserAuditHandler{auditRepository: repo}
}

// Handle returns audit events for a user.
func (h *getUserAuditHandler) Handle(ctx context.Context, query GetUserAuditQuery) (*PaginatedAuditEntries, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	query.PageSize = pagination.Clamp(query.PageSize, 20)

	events, totalCount, err := h.auditRepository.ListForUser(ctx, query.UserID, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	result := Paginate(totalCount, query.Page, query.PageSize)
	result.Events = MapEventsToEntries(events)
	return &result, nil
}
