package commands

import (
	"context"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

// DeleteUserCommand removes a staff account.
type DeleteUserCommand struct {
	UserID    string `json:"-"`
	ActorID   string
	IPAddress string
	UserAgent string
}

// DeleteUserHandler handles staff deletion.
type DeleteUserHandler interface {
	Handle(ctx context.Context, cmd DeleteUserCommand) error
}

type deleteUserHandler struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
	auditLogger       auditdomain.Logger
}

// NewDeleteUserHandler creates a handler for staff deletion.
func NewDeleteUserHandler(
	repo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	auditLogger auditdomain.Logger,
) DeleteUserHandler {
	return &deleteUserHandler{
		userRepository:    repo,
		sessionRepository: sessionRepo,
		auditLogger:       auditLogger,
	}
}

// Handle deletes a staff account when policy checks pass.
func (h *deleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	if cmd.UserID == cmd.ActorID {
		return domain.ErrCannotDeleteSelfError
	}

	user, err := h.userRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleAdmin {
		adminCount, err := h.userRepository.CountByRole(ctx, domain.RoleAdmin)
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return domain.ErrLastAdminError
		}
	}

	if err := h.sessionRepository.DeleteByUserID(ctx, cmd.UserID); err != nil {
		return err
	}
	if err := h.userRepository.Delete(ctx, cmd.UserID); err != nil {
		return err
	}

	h.logEvent(cmd, user.ID, true, map[string]string{
		"email": user.Email,
		"role":  string(user.Role),
	})
	return nil
}

func (h *deleteUserHandler) logEvent(cmd DeleteUserCommand, targetUserID uuid.UUID, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}
	if details == nil {
		details = map[string]string{}
	}
	details["targetUserId"] = targetUserID.String()

	var actorID *uuid.UUID
	if parsed, err := uuid.Parse(cmd.ActorID); err == nil {
		actorID = &parsed
	}

	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType:    auditdomain.EventUserDeleted,
		UserID:       actorID,
		ResourceType: "user",
		ResourceID:   targetUserID.String(),
		IPAddress:    cmd.IPAddress,
		UserAgent:    cmd.UserAgent,
		Details:      details,
		Success:      success,
	})
}
