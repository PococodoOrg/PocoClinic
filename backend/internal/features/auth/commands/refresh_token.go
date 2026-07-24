package commands

import (
	"context"
	"fmt"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/google/uuid"
)

// RefreshTokenCommand represents a token refresh request
type RefreshTokenCommand struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

// RefreshTokenResponse contains a new access token and optional rotated refresh token
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"-"`
	User         *domain.User
}

// RefreshTokenHandler refreshes an access token using a refresh token
type RefreshTokenHandler interface {
	Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResponse, error)
}

type refreshTokenHandler struct {
	userRepository    domain.RefreshUserRepository
	sessionRepository domain.RefreshSessionRepository
	tokenConfig       domain.TokenConfig
	auditLogger       auditdomain.Logger
}

// NewRefreshTokenHandler creates a refresh token handler
func NewRefreshTokenHandler(
	userRepo domain.RefreshUserRepository,
	sessionRepo domain.RefreshSessionRepository,
	tokenConfig domain.TokenConfig,
	auditLogger auditdomain.Logger,
) RefreshTokenHandler {
	return &refreshTokenHandler{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		tokenConfig:       tokenConfig,
		auditLogger:       auditLogger,
	}
}

// Handle validates the refresh token and issues a new access token
func (h *refreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResponse, error) {
	if cmd.RefreshToken == "" {
		return nil, domain.ErrInvalidTokenError
	}

	claims, err := domain.ValidateToken(cmd.RefreshToken, domain.TokenTypeRefresh, h.tokenConfig.RefreshTokenSecret)
	if err != nil {
		h.logRefresh(nil, cmd, false, map[string]string{"reason": "invalid_token"})
		return nil, domain.ErrInvalidTokenError
	}

	session, err := h.sessionRepository.GetByRefreshToken(ctx, domain.HashRefreshToken(cmd.RefreshToken))
	if err != nil {
		h.logRefresh(nil, cmd, false, map[string]string{"reason": "session_not_found"})
		return nil, domain.ErrInvalidTokenError
	}

	if session.IsExpired() {
		h.logRefresh(&session.UserID, cmd, false, map[string]string{"reason": "session_expired"})
		return nil, domain.ErrSessionExpiredError
	}

	if claims.UserID != session.UserID.String() {
		h.logRefresh(&session.UserID, cmd, false, map[string]string{"reason": "token_user_mismatch"})
		return nil, domain.ErrInvalidTokenError
	}

	user, err := h.userRepository.GetByID(ctx, claims.UserID)
	if err != nil {
		h.logRefresh(&session.UserID, cmd, false, map[string]string{"reason": "user_not_found"})
		return nil, domain.ErrInvalidTokenError
	}

	if !user.IsActive {
		h.logRefresh(&user.ID, cmd, false, map[string]string{"reason": "account_inactive"})
		return nil, domain.ErrSessionExpiredError
	}

	accessToken, refreshToken, err := session.GenerateTokens(user, h.tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	session.Refresh(h.tokenConfig.SessionInactivityTTL)
	if err := h.sessionRepository.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	h.logRefresh(&user.ID, cmd, true, nil)

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (h *refreshTokenHandler) logRefresh(userID *uuid.UUID, cmd RefreshTokenCommand, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}
	h.auditLogger.Log(context.Background(), auditdomain.Event{
		EventType: auditdomain.EventTokenRefresh,
		UserID:    userID,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
		Details:   details,
		Success:   success,
	})
}
