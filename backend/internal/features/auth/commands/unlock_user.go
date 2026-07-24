package commands

import (
	"context"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

// UnlockUserCommand clears a lockout after failed sign-in attempts.
type UnlockUserCommand struct {
	UserID    string `json:"-"`
	ActorID   string
	IPAddress string
	UserAgent string
}

// UnlockUserHandler handles administrator unlock of locked staff accounts.
type UnlockUserHandler interface {
	Handle(ctx context.Context, cmd UnlockUserCommand) (*domain.User, error)
}

type unlockUserHandler struct {
	userRepository domain.UserRepository
	auditLogger    auditdomain.Logger
}

// NewUnlockUserHandler creates a handler for account unlock.
func NewUnlockUserHandler(repo domain.UserRepository, auditLogger auditdomain.Logger) UnlockUserHandler {
	return &unlockUserHandler{
		userRepository: repo,
		auditLogger:    auditLogger,
	}
}

// Handle resets failed attempts and lockout for the target user.
func (h *unlockUserHandler) Handle(ctx context.Context, cmd UnlockUserCommand) (*domain.User, error) {
	user, err := h.userRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsLocked() {
		return user, nil
	}

	user.ResetFailedAttempts()

	if err := h.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	h.logUnlock(cmd, user.ID, true)
	return user, nil
}

func (h *unlockUserHandler) logUnlock(cmd UnlockUserCommand, targetUserID uuid.UUID, success bool) {
	if h.auditLogger == nil {
		return
	}

	var actorID *uuid.UUID
	if parsed, err := uuid.Parse(cmd.ActorID); err == nil {
		actorID = &parsed
	}

	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType:    auditdomain.EventUserUnlocked,
		UserID:       actorID,
		ResourceType: "user",
		ResourceID:   targetUserID.String(),
		IPAddress:    cmd.IPAddress,
		UserAgent:    cmd.UserAgent,
		Details: map[string]string{
			"targetUserId": targetUserID.String(),
		},
		Success: success,
	})
}
