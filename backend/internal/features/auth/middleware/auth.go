package middleware

import (
	"context"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/pkg/config"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/pkg/httperr"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
)

// PINCheckUserRepository loads users for PIN-change enforcement.
type PINCheckUserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// AuthSessionRepository loads sessions for access-token validation.
type AuthSessionRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Session, error)
}

// AuthMiddleware provides authentication and authorization middleware
type AuthMiddleware struct {
	tokenConfig domain.TokenConfig
	authConfig  config.AuthConfig
	userRepo    PINCheckUserRepository
	sessionRepo AuthSessionRepository
	logger      *logging.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(
	tokenConfig domain.TokenConfig,
	authConfig config.AuthConfig,
	userRepo PINCheckUserRepository,
	sessionRepo AuthSessionRepository,
	logger *logging.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenConfig: tokenConfig,
		authConfig:  authConfig,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (m *AuthMiddleware) abort(c *gin.Context, err error) {
	httperr.Respond(c, m.logger, err)
	c.Abort()
}

// RequireAuth validates the access token cookie against the live session and user record.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(m.authConfig.AccessTokenCookieName)
		if err != nil || token == "" {
			m.abort(c, pkgerrors.NewAPIError(pkgerrors.ErrUnauthorized, "missing access token"))
			return
		}

		claims, err := domain.ValidateToken(token, domain.TokenTypeAccess, m.tokenConfig.AccessTokenSecret)
		if err != nil {
			m.abort(c, domain.ErrInvalidTokenError)
			return
		}

		if claims.SessionID == "" {
			m.abort(c, domain.ErrInvalidTokenError)
			return
		}

		session, err := m.sessionRepo.GetByID(c.Request.Context(), claims.SessionID)
		if err != nil || session.IsExpired() || session.UserID.String() != claims.UserID {
			m.abort(c, domain.ErrInvalidTokenError)
			return
		}

		user, err := m.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			m.abort(c, err)
			return
		}
		if !user.IsActive {
			m.abort(c, domain.ErrAccountInactiveError)
			return
		}

		c.Set("userID", user.ID.String())
		c.Set("userRole", user.Role)
		c.Set("sessionID", claims.SessionID)
		c.Next()
	}
}

// RequirePINChanged blocks PHI and admin routes until the user changes their default PIN.
func (m *AuthMiddleware) RequirePINChanged() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			m.abort(c, pkgerrors.NewAPIError(pkgerrors.ErrUnauthorized, "missing user ID"))
			return
		}

		user, err := m.userRepo.GetByID(c.Request.Context(), userID.(string))
		if err != nil {
			m.abort(c, err)
			return
		}
		if user.MustChangePIN {
			m.abort(c, pkgerrors.NewAPIError(pkgerrors.ErrForbidden, "PIN change required"))
			return
		}

		c.Next()
	}
}

// RequireRole ensures the user has one of the required roles
func (m *AuthMiddleware) RequireRole(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			m.abort(c, pkgerrors.NewAPIError(pkgerrors.ErrUnauthorized, "missing user role"))
			return
		}

		userRole := role.(domain.Role)
		for _, requiredRole := range roles {
			if userRole == requiredRole {
				c.Next()
				return
			}
		}

		m.abort(c, pkgerrors.NewAPIError(pkgerrors.ErrForbidden, "insufficient permissions"))
	}
}
