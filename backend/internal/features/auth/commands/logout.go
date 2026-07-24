package commands

import (
	"context"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

// LogoutCommand represents a logout request
type LogoutCommand struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

// LogoutHandler ends a user session
type LogoutHandler interface {
	Handle(ctx context.Context, cmd LogoutCommand) error
}

type logoutHandler struct {
	sessionRepository domain.SessionRepository
	auditLogger       auditdomain.Logger
}

// NewLogoutHandler creates a logout handler
func NewLogoutHandler(sessionRepo domain.SessionRepository, auditLogger auditdomain.Logger) LogoutHandler {
	return &logoutHandler{
		sessionRepository: sessionRepo,
		auditLogger:       auditLogger,
	}
}

// Handle removes the session associated with the refresh token
func (h *logoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	var userID *uuid.UUID

	if cmd.RefreshToken != "" {
		session, err := h.sessionRepository.GetByRefreshToken(ctx, domain.HashRefreshToken(cmd.RefreshToken))
		if err == nil {
			userID = &session.UserID
			_ = h.sessionRepository.Delete(ctx, session.ID.String())
		}
	}

	h.logLogout(userID, cmd, true)
	return nil
}

func (h *logoutHandler) logLogout(userID *uuid.UUID, cmd LogoutCommand, success bool) {
	if h.auditLogger == nil {
		return
	}
	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType: auditdomain.EventLogout,
		UserID:    userID,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
		Success:   success,
	})
}
