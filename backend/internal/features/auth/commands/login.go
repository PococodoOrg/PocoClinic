package commands

import (
	"context"
	"fmt"
	"time"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

const (
	LoginModeStaff = "staff"
	LoginModeAdmin = "admin"
)

// LoginResponse represents a successful authentication response
type LoginResponse struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"-"`
}

type loginAttemptContext struct {
	Mode      string
	Email     string
	IPAddress string
	UserAgent string
}

// StaffLoginCommand authenticates via badge QR key + PIN (daily staff sign-in)
type StaffLoginCommand struct {
	Key       string `json:"key" binding:"required"`
	PIN       string `json:"pin" binding:"required,len=4"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
}

// AdminLoginCommand authenticates via email + key + PIN (bootstrap and administrator setup)
type AdminLoginCommand struct {
	Email     string `json:"email" binding:"required,email"`
	Key       string `json:"key" binding:"required"`
	PIN       string `json:"pin" binding:"required,len=4"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
}

// LoginHandler handles staff and administrator login flows
type LoginHandler interface {
	HandleStaff(ctx context.Context, cmd StaffLoginCommand) (*LoginResponse, error)
	HandleAdmin(ctx context.Context, cmd AdminLoginCommand) (*LoginResponse, error)
}

type loginHandler struct {
	userRepository    domain.LoginUserRepository
	sessionRepository domain.CreateSessionRepository
	tokenConfig       domain.TokenConfig
	auditLogger       auditdomain.Logger
}

// NewLoginHandler creates a new handler for user login
func NewLoginHandler(
	userRepo domain.LoginUserRepository,
	sessionRepo domain.CreateSessionRepository,
	tokenConfig domain.TokenConfig,
	auditLogger auditdomain.Logger,
) LoginHandler {
	return &loginHandler{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		tokenConfig:       tokenConfig,
		auditLogger:       auditLogger,
	}
}

// HandleStaff processes badge + PIN login
func (h *loginHandler) HandleStaff(ctx context.Context, cmd StaffLoginCommand) (*LoginResponse, error) {
	user, err := h.userRepository.FindByKey(ctx, cmd.Key)
	if err != nil {
		h.logAttempt(nil, loginAttemptContext{
			Mode: LoginModeStaff, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent,
		}, auditdomain.EventLoginFailure, false, map[string]string{"reason": "invalid_credentials"})
		return nil, domain.ErrInvalidCredentialsError
	}

	return h.completeLogin(ctx, user, cmd.Key, cmd.PIN, loginAttemptContext{
		Mode: LoginModeStaff, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent,
	})
}

// HandleAdmin processes email + key + PIN login for initial setup and break-glass access
func (h *loginHandler) HandleAdmin(ctx context.Context, cmd AdminLoginCommand) (*LoginResponse, error) {
	user, err := h.userRepository.GetByEmail(ctx, cmd.Email)
	if err != nil {
		h.logAttempt(nil, loginAttemptContext{
			Mode: LoginModeAdmin, Email: cmd.Email, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent,
		}, auditdomain.EventLoginFailure, false, map[string]string{"reason": "invalid_credentials"})
		return nil, domain.ErrInvalidCredentialsError
	}

	if user.Role != domain.RoleAdmin {
		h.logAttempt(&user.ID, loginAttemptContext{
			Mode: LoginModeAdmin, Email: cmd.Email, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent,
		}, auditdomain.EventLoginFailure, false, map[string]string{"reason": "invalid_credentials"})
		return nil, domain.ErrInvalidCredentialsError
	}

	return h.completeLogin(ctx, user, cmd.Key, cmd.PIN, loginAttemptContext{
		Mode: LoginModeAdmin, Email: cmd.Email, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent,
	})
}

func (h *loginHandler) completeLogin(
	ctx context.Context,
	user *domain.User,
	key, pin string,
	attempt loginAttemptContext,
) (*LoginResponse, error) {
	if !user.IsActive {
		h.logAttempt(&user.ID, attempt, auditdomain.EventLoginFailure, false, map[string]string{"reason": "account_inactive"})
		return nil, domain.ErrAccountInactiveError
	}

	if user.IsLocked() {
		h.logAttempt(&user.ID, attempt, auditdomain.EventLoginLocked, false, map[string]string{"reason": "account_locked"})
		return nil, domain.ErrAccountLockedError
	}

	if !user.ValidateCredentials(key, pin) {
		user.RecordFailedAttempt()
		_ = h.userRepository.Update(ctx, user)
		h.logAttempt(&user.ID, attempt, auditdomain.EventLoginFailure, false, map[string]string{"reason": "invalid_credentials"})
		return nil, domain.ErrInvalidCredentialsError
	}

	session := domain.NewSession(
		user.ID,
		attempt.UserAgent,
		attempt.IPAddress,
		time.Now().Add(h.tokenConfig.SessionInactivityTTL),
	)

	accessToken, refreshToken, err := session.GenerateTokens(user, h.tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	if err := h.sessionRepository.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	user.RecordLogin()
	_ = h.userRepository.Update(ctx, user)
	h.logAttempt(&user.ID, attempt, auditdomain.EventLoginSuccess, true, map[string]string{"mode": attempt.Mode})

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *loginHandler) logAttempt(userID *uuid.UUID, attempt loginAttemptContext, eventType string, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}
	if details == nil {
		details = map[string]string{}
	}
	details["loginMode"] = attempt.Mode
	if attempt.Email != "" {
		details["email"] = attempt.Email
	}

	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType: eventType,
		UserID:    userID,
		IPAddress: attempt.IPAddress,
		UserAgent: attempt.UserAgent,
		Details:   details,
		Success:   success,
	})
}
