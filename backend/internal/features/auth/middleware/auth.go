package middleware

import (
	"net/http"
	"strings"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides authentication and authorization middleware
type AuthMiddleware struct {
	tokenConfig domain.TokenConfig
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(config domain.TokenConfig) *AuthMiddleware {
	return &AuthMiddleware{
		tokenConfig: config,
	}
}

// RequireAuth validates the access token and adds user claims to the context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		// Check Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		// Validate token
		claims, err := domain.ValidateToken(parts[1], domain.TokenTypeAccess, m.tokenConfig.AccessTokenSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Add claims to context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// RequireRole ensures the user has one of the required roles
func (m *AuthMiddleware) RequireRole(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user role"})
			return
		}

		userRole := role.(domain.Role)
		for _, requiredRole := range roles {
			if userRole == requiredRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
