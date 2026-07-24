package commands

import (
	"context"
	"fmt"
	"time"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/google/uuid"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
)

// UpdateUserCommand updates staff profile fields.
type UpdateUserCommand struct {
	UserID    string      `json:"-"`
	Email     string      `json:"email" binding:"required,email"`
	Name      string      `json:"name" binding:"required"`
	Role      domain.Role `json:"role" binding:"required"`
	IsActive  bool        `json:"isActive"`
	ActorID   string
	IPAddress string
	UserAgent string
}

// UpdateUserHandler handles staff profile updates.
type UpdateUserHandler interface {
	Handle(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error)
}

type updateUserHandler struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
	auditLogger       auditdomain.Logger
}

// NewUpdateUserHandler creates a handler for staff updates.
func NewUpdateUserHandler(
	repo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	auditLogger auditdomain.Logger,
) UpdateUserHandler {
	return &updateUserHandler{
		userRepository:    repo,
		sessionRepository: sessionRepo,
		auditLogger:       auditLogger,
	}
}

// Handle updates a user's profile.
func (h *updateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error) {
	if !domain.IsAssignableStaffRole(cmd.Role) {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid role")
	}

	user, err := h.userRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if user.Role == domain.RoleAdmin && cmd.Role != domain.RoleAdmin {
		adminCount, err := h.userRepository.CountByRole(ctx, domain.RoleAdmin)
		if err != nil {
			return nil, err
		}
		if adminCount <= 1 {
			return nil, domain.ErrLastAdminError
		}
	}

	if user.ID.String() == cmd.ActorID && !cmd.IsActive {
		return nil, domain.NewAuthError(domain.ErrCannotDeleteSelf, "you cannot deactivate your own account")
	}

	if user.Role == domain.RoleAdmin && !cmd.IsActive {
		adminCount, err := h.userRepository.CountByRole(ctx, domain.RoleAdmin)
		if err != nil {
			return nil, err
		}
		if adminCount <= 1 {
			return nil, domain.ErrLastAdminError
		}
	}

	if cmd.Email != user.Email {
		existing, err := h.userRepository.GetByEmail(ctx, cmd.Email)
		if err == nil && existing.ID != user.ID {
			return nil, domain.ErrEmailTakenError(cmd.Email)
		}
		if err != nil {
			if authErr, ok := err.(*domain.AuthError); !ok || authErr.Code != domain.ErrUserNotFound {
				return nil, err
			}
		}
	}

	wasActive := user.IsActive
	previousRole := user.Role

	user.Email = cmd.Email
	user.Name = cmd.Name
	user.Role = cmd.Role
	user.IsActive = cmd.IsActive
	user.UpdatedAt = time.Now()

	if err := h.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	if (!cmd.IsActive && wasActive) || cmd.Role != previousRole {
		if err := h.sessionRepository.DeleteByUserID(ctx, cmd.UserID); err != nil {
			return nil, err
		}
	}

	h.logEvent(cmd, user.ID, true, map[string]string{
		"email":    user.Email,
		"role":     string(user.Role),
		"isActive": fmt.Sprintf("%t", user.IsActive),
	})
	return user, nil
}

func (h *updateUserHandler) logEvent(cmd UpdateUserCommand, targetUserID uuid.UUID, success bool, details map[string]string) {
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
		EventType:    auditdomain.EventUserUpdated,
		UserID:       actorID,
		ResourceType: "user",
		ResourceID:   targetUserID.String(),
		IPAddress:    cmd.IPAddress,
		UserAgent:    cmd.UserAgent,
		Details:      details,
		Success:      success,
	})
}
