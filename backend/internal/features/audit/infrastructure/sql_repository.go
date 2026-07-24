package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/google/uuid"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

const auditSelectQuery = `
	SELECT id, event_type, user_id, resource_type, resource_id,
	       ip_address, user_agent, details, success, created_at
	FROM audit_logs
`

// SQLRepository stores audit events in SQLite.
type SQLRepository struct {
	pool *database.DB
}

// NewSQLRepository creates a SQL-backed audit repository.
func NewSQLRepository(pool *database.DB) *SQLRepository {
	return &SQLRepository{pool: pool}
}

// Create inserts an audit event.
func (r *SQLRepository) Create(ctx context.Context, event *domain.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	details, err := json.Marshal(event.Details)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO audit_logs (
			id, event_type, user_id, resource_type, resource_id,
			ip_address, user_agent, details, success, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		event.ID,
		event.EventType,
		event.UserID,
		nullIfEmpty(event.ResourceType),
		nullIfEmpty(event.ResourceID),
		nullIfEmpty(event.IPAddress),
		nullIfEmpty(event.UserAgent),
		details,
		event.Success,
		event.CreatedAt,
	)
	return err
}

// ListForUser returns audit events performed by or targeting a user.
func (r *SQLRepository) ListForUser(ctx context.Context, userID string, page, pageSize int) ([]*domain.Event, int64, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user id")
	}

	filter := `WHERE user_id = $1 OR (resource_type = 'user' AND resource_id = $1)`

	var totalCount int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs `+filter, parsedID).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx,
		auditSelectQuery+filter+` ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		parsedID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	events := make([]*domain.Event, 0)
	for rows.Next() {
		event, err := scanAuditEvent(rows)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}
	return events, totalCount, rows.Err()
}

// ListRecent returns audit events across the system with optional filters.
func (r *SQLRepository) ListRecent(ctx context.Context, filter domain.ListFilter) ([]*domain.Event, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 25
	}

	where := `WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.EventType != "" {
		where += fmt.Sprintf(" AND event_type = $%d", argIndex)
		args = append(args, filter.EventType)
		argIndex++
	}
	if filter.UserID != "" {
		parsedID, err := uuid.Parse(filter.UserID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid user id")
		}
		where += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, parsedID)
		argIndex++
	}
	if filter.SinceDays > 0 {
		where += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, time.Now().UTC().AddDate(0, 0, -filter.SinceDays))
		argIndex++
	}

	countQuery := `SELECT COUNT(*) FROM audit_logs ` + where
	var totalCount int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listArgs := append(append([]interface{}{}, args...), pageSize, offset)
	limitOffset := fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)

	rows, err := r.pool.Query(ctx, auditSelectQuery+where+limitOffset, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	events := make([]*domain.Event, 0)
	for rows.Next() {
		event, err := scanAuditEvent(rows)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}
	return events, totalCount, rows.Err()
}

func scanAuditEvent(row database.Row) (*domain.Event, error) {
	event := &domain.Event{}
	var userID *uuid.UUID
	var resourceType *string
	var resourceID *string
	var ipAddress *string
	var userAgent *string
	var details []byte

	err := row.Scan(
		&event.ID,
		&event.EventType,
		&userID,
		&resourceType,
		&resourceID,
		&ipAddress,
		&userAgent,
		&details,
		&event.Success,
		database.Time(&event.CreatedAt),
	)
	if err != nil {
		return nil, err
	}

	event.UserID = userID
	if resourceType != nil {
		event.ResourceType = *resourceType
	}
	if resourceID != nil {
		event.ResourceID = *resourceID
	}
	if ipAddress != nil {
		event.IPAddress = *ipAddress
	}
	if userAgent != nil {
		event.UserAgent = *userAgent
	}
	if len(details) > 0 {
		_ = json.Unmarshal(details, &event.Details)
	}
	return event, nil
}

func nullIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
