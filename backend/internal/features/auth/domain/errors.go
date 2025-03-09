package domain

import "fmt"

// AuthError represents a domain-specific authentication error
type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// Common auth error codes
const (
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrAccountLocked      = "ACCOUNT_LOCKED"
	ErrEmailTaken         = "EMAIL_TAKEN"
	ErrUserNotFound       = "USER_NOT_FOUND"
	ErrSessionNotFound    = "SESSION_NOT_FOUND"
	ErrInvalidToken       = "INVALID_TOKEN"
)

// NewAuthError creates a new auth error
func NewAuthError(code string, message string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}

// Common auth errors
var (
	ErrInvalidCredentialsError = NewAuthError(ErrInvalidCredentials, "invalid credentials")
	ErrAccountLockedError      = NewAuthError(ErrAccountLocked, "account is locked")
	ErrEmailTakenError         = func(email string) *AuthError {
		return NewAuthError(ErrEmailTaken, fmt.Sprintf("email %s is already registered", email))
	}
	ErrUserNotFoundError    = NewAuthError(ErrUserNotFound, "user not found")
	ErrSessionNotFoundError = NewAuthError(ErrSessionNotFound, "session not found")
	ErrInvalidTokenError    = NewAuthError(ErrInvalidToken, "invalid token")
)
