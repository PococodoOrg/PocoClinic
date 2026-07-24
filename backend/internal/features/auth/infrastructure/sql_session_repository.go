package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/google/uuid"
)

// SQLSessionRepository persists sessions in SQLite.
type SQLSessionRepository struct {
	pool *database.DB
}

// NewSQLSessionRepository creates a SQL-backed session repository.
func NewSQLSessionRepository(pool *database.DB) *SQLSessionRepository {
	return &SQLSessionRepository{pool: pool}
}

// Create adds a new session.
func (r *SQLSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sessions (
			id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
		session.CreatedAt,
		session.UpdatedAt,
	)
	return err
}

// Update modifies an existing session.
func (r *SQLSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE sessions SET
			refresh_token = $2,
			user_agent = $3,
			ip_address = $4,
			expires_at = $5,
			updated_at = $6
		WHERE id = $1
	`,
		session.ID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
		session.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSessionNotFoundError
	}
	return nil
}

// Delete removes a session.
func (r *SQLSessionRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSessionNotFoundError
	}
	return nil
}

// GetByID retrieves a session by ID.
func (r *SQLSessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	row := r.pool.QueryRow(ctx, sessionSelectQuery+` WHERE id = $1`, id)
	return scanSession(row)
}

// GetByRefreshToken retrieves a session by refresh token.
func (r *SQLSessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	row := r.pool.QueryRow(ctx, sessionSelectQuery+` WHERE refresh_token = $1`, token)
	session, err := scanSession(row)
	if errors.Is(err, database.ErrNoRows) {
		return nil, domain.ErrSessionNotFoundError
	}
	return session, err
}

// DeleteExpired removes expired sessions.
func (r *SQLSessionRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < $1`, time.Now().UTC())
	return err
}

// DeleteByUserID removes all sessions for a user.
func (r *SQLSessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

// DeleteByUserIDExcept removes all sessions for a user except the given session ID.
func (r *SQLSessionRepository) DeleteByUserIDExcept(ctx context.Context, userID, exceptSessionID string) error {
	if exceptSessionID == "" {
		return r.DeleteByUserID(ctx, userID)
	}
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1 AND id != $2`, userID, exceptSessionID)
	return err
}

const sessionSelectQuery = `
	SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at, updated_at
	FROM sessions
`

func scanSession(row database.Row) (*domain.Session, error) {
	session := &domain.Session{}
	var userID uuid.UUID
	err := row.Scan(
		&session.ID,
		&userID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		database.Time(&session.ExpiresAt),
		database.Time(&session.CreatedAt),
		database.Time(&session.UpdatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, domain.ErrSessionNotFoundError
		}
		return nil, err
	}
	session.UserID = userID
	return session, nil
}

var _ domain.SessionRepository = (*SQLSessionRepository)(nil)
