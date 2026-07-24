package queries

import (
	"time"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
)

// AuditEntry is the API representation of an audit event.
type AuditEntry struct {
	ID           string            `json:"id"`
	EventType    string            `json:"eventType"`
	UserID       *string           `json:"userId,omitempty"`
	ResourceType string            `json:"resourceType,omitempty"`
	ResourceID   string            `json:"resourceId,omitempty"`
	IPAddress    string            `json:"ipAddress,omitempty"`
	UserAgent    string            `json:"userAgent,omitempty"`
	Details      map[string]string `json:"details,omitempty"`
	Success      bool              `json:"success"`
	CreatedAt    time.Time         `json:"createdAt"`
}

// PaginatedAuditEntries represents paginated audit history.
type PaginatedAuditEntries struct {
	Events      []AuditEntry `json:"events"`
	TotalCount  int64        `json:"totalCount"`
	CurrentPage int          `json:"currentPage"`
	PageSize    int          `json:"pageSize"`
	TotalPages  int          `json:"totalPages"`
}

func MapEventsToEntries(events []*auditdomain.Event) []AuditEntry {
	entries := make([]AuditEntry, 0, len(events))
	for _, event := range events {
		entry := AuditEntry{
			ID:           event.ID.String(),
			EventType:    event.EventType,
			ResourceType: event.ResourceType,
			ResourceID:   event.ResourceID,
			IPAddress:    event.IPAddress,
			UserAgent:    event.UserAgent,
			Details:      event.Details,
			Success:      event.Success,
			CreatedAt:    event.CreatedAt,
		}
		if event.UserID != nil {
			id := event.UserID.String()
			entry.UserID = &id
		}
		entries = append(entries, entry)
	}
	return entries
}

func Paginate(totalCount int64, page, pageSize int) PaginatedAuditEntries {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	return PaginatedAuditEntries{
		TotalCount:  totalCount,
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}
}
