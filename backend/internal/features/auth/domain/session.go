package domain

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// Claims extends jwt.RegisteredClaims with custom claims
type Claims struct {
	jwt.RegisteredClaims
	UserID    string    `json:"uid"`
	Role      Role      `json:"role"`
	TokenType TokenType `json:"type"`
}

// Session represents an active user session
type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	RefreshToken string    `json:"-"`
	UserAgent    string    `json:"userAgent"`
	IPAddress    string    `json:"ipAddress"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TokenConfig holds JWT token configuration
type TokenConfig struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	Issuer             string
}

// NewSession creates a new session for a user
func NewSession(userID uuid.UUID, userAgent, ipAddress string, expiresAt time.Time) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GenerateTokens creates both access and refresh tokens
func (s *Session) GenerateTokens(user *User, config TokenConfig) (accessToken string, refreshToken string, err error) {
	accessToken, err = generateToken(TokenTypeAccess, user, config.AccessTokenSecret, config.AccessTokenTTL, config.Issuer)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = generateToken(TokenTypeRefresh, user, config.RefreshTokenSecret, config.RefreshTokenTTL, config.Issuer)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	s.RefreshToken = refreshToken
	return accessToken, refreshToken, nil
}

// ValidateToken verifies a JWT token's signature and claims
func ValidateToken(tokenString string, tokenType TokenType, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, fmt.Errorf("invalid token type")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateToken creates a new JWT token
func generateToken(tokenType TokenType, user *User, secret []byte, ttl time.Duration, issuer string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    issuer,
			Subject:   user.ID.String(),
		},
		UserID:    user.ID.String(),
		Role:      user.Role,
		TokenType: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Refresh extends the session's expiration time
func (s *Session) Refresh(ttl time.Duration) {
	now := time.Now()
	s.ExpiresAt = now.Add(ttl)
	s.UpdatedAt = now
}
