package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

// LoginCommand represents the login command
type LoginCommand struct {
	Email     string `json:"email" binding:"required,email"`
	Key       string `json:"key" binding:"required"`
	PIN       string `json:"pin" binding:"required,len=4"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User        *domain.User `json:"user"`
	AccessToken string       `json:"accessToken"`
}

// LoginHandler handles user login
type LoginHandler interface {
	Handle(ctx context.Context, cmd LoginCommand) (*LoginResponse, error)
}

// loginHandler implements LoginHandler
type loginHandler struct {
	userRepository    domain.ValidateUserRepository
	sessionRepository domain.CreateSessionRepository
	tokenConfig       domain.TokenConfig
}

// NewLoginHandler creates a new handler for user login
func NewLoginHandler(
	userRepo domain.ValidateUserRepository,
	sessionRepo domain.CreateSessionRepository,
	tokenConfig domain.TokenConfig,
) LoginHandler {
	return &loginHandler{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		tokenConfig:       tokenConfig,
	}
}

// Handle processes the login command
func (h *loginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResponse, error) {
	// Get user by email
	user, err := h.userRepository.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is locked
	if user.IsLocked() {
		return nil, fmt.Errorf("account is locked")
	}

	// Validate credentials
	if !user.ValidateCredentials(cmd.Key, cmd.PIN) {
		user.RecordFailedAttempt()
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create new session
	session := domain.NewSession(
		user.ID,
		cmd.UserAgent,
		cmd.IPAddress,
		time.Now().Add(24*time.Hour),
	)

	// Generate tokens
	accessToken, _, err := session.GenerateTokens(user, h.tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save session
	err = h.sessionRepository.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Record successful login
	user.RecordLogin()

	return &LoginResponse{
		User:        user,
		AccessToken: accessToken,
	}, nil
}
