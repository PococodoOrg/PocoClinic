package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type sessionTestSuite struct {
	defaultUser    *User
	defaultConfig  TokenConfig
	defaultSession *Session
}

func setupSessionTest() sessionTestSuite {
	user := &User{
		ID:   uuid.New(),
		Name: "Dr. Smith",
		Role: RoleDoctor,
	}

	config := TokenConfig{
		AccessTokenSecret:  []byte("access-secret-key-for-test"),
		RefreshTokenSecret: []byte("refresh-secret-key-for-test"),
		AccessTokenTTL:     15 * time.Minute,
		RefreshTokenTTL:    24 * time.Hour,
		Issuer:             "poco-clinic-test",
	}

	session := NewSession(
		user.ID,
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124",
		"192.168.1.1",
		time.Now().Add(24*time.Hour),
	)

	return sessionTestSuite{
		defaultUser:    user,
		defaultConfig:  config,
		defaultSession: session,
	}
}

func TestNewSessionCreation(t *testing.T) {
	suite := setupSessionTest()
	session := suite.defaultSession

	assert.NotEqual(t, uuid.Nil, session.ID)
	assert.Equal(t, suite.defaultUser.ID, session.UserID)
	assert.False(t, session.CreatedAt.IsZero())
	assert.False(t, session.UpdatedAt.IsZero())
	assert.False(t, session.ExpiresAt.IsZero())
}

func TestTokenGeneration(t *testing.T) {
	suite := setupSessionTest()

	accessToken, refreshToken, err := suite.defaultSession.GenerateTokens(
		suite.defaultUser,
		suite.defaultConfig,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.Equal(t, refreshToken, suite.defaultSession.RefreshToken)

	accessClaims, err := ValidateToken(accessToken, TokenTypeAccess, suite.defaultConfig.AccessTokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, suite.defaultUser.ID.String(), accessClaims.UserID)
	assert.Equal(t, suite.defaultUser.Role, accessClaims.Role)
	assert.Equal(t, TokenTypeAccess, accessClaims.TokenType)

	refreshClaims, err := ValidateToken(refreshToken, TokenTypeRefresh, suite.defaultConfig.RefreshTokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, suite.defaultUser.ID.String(), refreshClaims.UserID)
	assert.Equal(t, suite.defaultUser.Role, refreshClaims.Role)
	assert.Equal(t, TokenTypeRefresh, refreshClaims.TokenType)
}

func TestSessionExpiration(t *testing.T) {
	suite := setupSessionTest()
	session := suite.defaultSession

	assert.False(t, session.IsExpired())

	session.ExpiresAt = time.Now().Add(-time.Hour)
	assert.True(t, session.IsExpired())
}

func TestSessionRefresh(t *testing.T) {
	suite := setupSessionTest()
	session := suite.defaultSession
	originalExpiry := session.ExpiresAt
	originalUpdate := session.UpdatedAt

	t.Logf("Original expiry: %v", originalExpiry)
	t.Logf("Original update: %v", originalUpdate)

	time.Sleep(time.Millisecond)
	session.Refresh(48 * time.Hour)

	t.Logf("New expiry: %v", session.ExpiresAt)
	t.Logf("New update: %v", session.UpdatedAt)

	if !session.ExpiresAt.After(originalExpiry) {
		t.Errorf("Expected new expiry %v to be after original expiry %v", session.ExpiresAt, originalExpiry)
	}
	if !session.UpdatedAt.After(originalUpdate) {
		t.Errorf("Expected new update time %v to be after original update time %v", session.UpdatedAt, originalUpdate)
	}
}

func TestTokenValidationFailures(t *testing.T) {
	suite := setupSessionTest()

	testCases := []struct {
		name          string
		tokenType     TokenType
		secret        []byte
		expectedError string
	}{
		{
			name:          "wrong_token_type",
			tokenType:     TokenTypeRefresh,
			secret:        suite.defaultConfig.AccessTokenSecret,
			expectedError: "invalid token type",
		},
		{
			name:          "wrong_secret",
			tokenType:     TokenTypeAccess,
			secret:        []byte("wrong-secret"),
			expectedError: "failed to parse token",
		},
	}

	accessToken, _, err := suite.defaultSession.GenerateTokens(suite.defaultUser, suite.defaultConfig)
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateToken(accessToken, tc.tokenType, tc.secret)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}
