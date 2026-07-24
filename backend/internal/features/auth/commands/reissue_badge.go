package commands

import (
	"context"
	"fmt"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

// ReissueBadgeCommand rotates a user's badge key after loss or compromise.
type ReissueBadgeCommand struct {
	UserID    string `json:"userId" binding:"required"`
	ActorID   string
	IPAddress string
	UserAgent string
}

// ReissueBadgeHandler handles badge reissue for administrators.
type ReissueBadgeHandler interface {
	Handle(ctx context.Context, cmd ReissueBadgeCommand) (*domain.User, string, error)
}

type reissueBadgeHandler struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
	auditLogger       auditdomain.Logger
}

// NewReissueBadgeHandler creates a handler for badge reissue.
func NewReissueBadgeHandler(
	repo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	auditLogger auditdomain.Logger,
) ReissueBadgeHandler {
	return &reissueBadgeHandler{
		userRepository:    repo,
		sessionRepository: sessionRepo,
		auditLogger:       auditLogger,
	}
}

// Handle generates a new badge key for the target user.
func (h *reissueBadgeHandler) Handle(ctx context.Context, cmd ReissueBadgeCommand) (*domain.User, string, error) {
	user, err := h.userRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, "", err
	}

	key, keyCred, err := domain.GenerateKey()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate key: %w", err)
	}

	user.SetKeyCredential(keyCred, key)
	user.ResetFailedAttempts()

	if err := h.userRepository.Update(ctx, user); err != nil {
		return nil, "", err
	}

	if err := h.sessionRepository.DeleteByUserID(ctx, cmd.UserID); err != nil {
		return nil, "", err
	}

	h.logReissue(cmd, user.ID, true, nil)
	return user, key, nil
}

func (h *reissueBadgeHandler) logReissue(cmd ReissueBadgeCommand, targetUserID uuid.UUID, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}
	if details == nil {
		details = map[string]string{}
	}
	details["targetUserId"] = targetUserID.String()
	if cmd.ActorID != "" {
		details["actorId"] = cmd.ActorID
	}

	var actorID *uuid.UUID
	if parsed, err := uuid.Parse(cmd.ActorID); err == nil {
		actorID = &parsed
	}

	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType:    auditdomain.EventBadgeReissued,
		UserID:       actorID,
		ResourceType: "user",
		ResourceID:   targetUserID.String(),
		IPAddress:    cmd.IPAddress,
		UserAgent:    cmd.UserAgent,
		Details:      details,
		Success:      success,
	})
}
