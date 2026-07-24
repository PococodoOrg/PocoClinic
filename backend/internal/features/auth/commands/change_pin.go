package commands

import (
	"context"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	"github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

// ChangePINCommand updates the authenticated user's PIN.
type ChangePINCommand struct {
	UserID     string `json:"-"`
	SessionID  string `json:"-"`
	CurrentPIN string `json:"currentPin" binding:"required,len=4"`
	NewPIN     string `json:"newPin" binding:"required,len=4"`
	IPAddress  string
	UserAgent  string
}

// ChangePINHandler handles self-service PIN changes.
type ChangePINHandler interface {
	Handle(ctx context.Context, cmd ChangePINCommand) (*domain.User, error)
}

type changePINHandler struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
	auditLogger       auditdomain.Logger
}

// NewChangePINHandler creates a handler for PIN changes.
func NewChangePINHandler(
	repo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	auditLogger auditdomain.Logger,
) ChangePINHandler {
	return &changePINHandler{
		userRepository:    repo,
		sessionRepository: sessionRepo,
		auditLogger:       auditLogger,
	}
}

// Handle verifies the current PIN and stores a new one.
func (h *changePINHandler) Handle(ctx context.Context, cmd ChangePINCommand) (*domain.User, error) {
	if err := domain.ValidatePINFormat(cmd.CurrentPIN); err != nil {
		return nil, domain.NewAuthError(domain.ErrWeakPIN, err.Error())
	}
	if err := domain.ValidatePINFormat(cmd.NewPIN); err != nil {
		return nil, domain.NewAuthError(domain.ErrWeakPIN, err.Error())
	}
	if cmd.CurrentPIN == cmd.NewPIN {
		return nil, domain.ErrWeakPINError
	}

	user, err := h.userRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if user.IsLocked() {
		h.logChange(cmd, false, map[string]string{"reason": "account_locked"})
		return nil, domain.ErrAccountLockedError
	}

	if !user.ValidatePIN(cmd.CurrentPIN) {
		user.RecordFailedAttempt()
		_ = h.userRepository.Update(ctx, user)
		h.logChange(cmd, false, map[string]string{"reason": "invalid_current_pin"})
		if user.IsLocked() {
			return nil, domain.ErrAccountLockedError
		}
		return nil, domain.ErrInvalidPINError
	}

	pinCred, err := domain.GeneratePINCredential(cmd.NewPIN)
	if err != nil {
		return nil, domain.NewAuthError(domain.ErrWeakPIN, err.Error())
	}

	user.SetPINCredential(pinCred)
	user.ClearMustChangePIN()
	user.ResetFailedAttempts()
	if err := h.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	if h.sessionRepository != nil {
		if err := h.sessionRepository.DeleteByUserIDExcept(ctx, cmd.UserID, cmd.SessionID); err != nil {
			return nil, err
		}
	}

	h.logChange(cmd, true, nil)
	return user, nil
}

func (h *changePINHandler) logChange(cmd ChangePINCommand, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}
	if details == nil {
		details = map[string]string{}
	}

	var userID *uuid.UUID
	if parsed, err := uuid.Parse(cmd.UserID); err == nil {
		userID = &parsed
	}

	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType: auditdomain.EventPINChanged,
		UserID:    userID,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
		Details:   details,
		Success:   success,
	})
}
