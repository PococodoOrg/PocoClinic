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
	ErrAccountInactive    = "ACCOUNT_INACTIVE"
	ErrEmailTaken         = "EMAIL_TAKEN"
	ErrUserNotFound       = "USER_NOT_FOUND"
	ErrSessionNotFound    = "SESSION_NOT_FOUND"
	ErrInvalidToken       = "INVALID_TOKEN"
	ErrSessionExpired     = "SESSION_EXPIRED"
	ErrInvalidPIN         = "INVALID_PIN"
	ErrWeakPIN            = "WEAK_PIN"
	ErrLastAdmin          = "LAST_ADMIN"
	ErrCannotDeleteSelf   = "CANNOT_DELETE_SELF"
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
	ErrAccountInactiveError    = NewAuthError(ErrAccountInactive, "account is deactivated")
	ErrEmailTakenError         = func(email string) *AuthError {
		return NewAuthError(ErrEmailTaken, fmt.Sprintf("email %s is already registered", email))
	}
	ErrUserNotFoundError    = NewAuthError(ErrUserNotFound, "user not found")
	ErrSessionNotFoundError = NewAuthError(ErrSessionNotFound, "session not found")
	ErrInvalidTokenError    = NewAuthError(ErrInvalidToken, "invalid token")
	ErrSessionExpiredError  = NewAuthError(ErrSessionExpired, "session expired due to inactivity")
	ErrInvalidPINError      = NewAuthError(ErrInvalidPIN, "current PIN is incorrect")
	ErrWeakPINError         = NewAuthError(ErrWeakPIN, "new PIN must be four digits and different from the current PIN")
	ErrLastAdminError       = NewAuthError(ErrLastAdmin, "cannot remove or delete the last administrator")
	ErrCannotDeleteSelfError = NewAuthError(ErrCannotDeleteSelf, "you cannot delete your own account")
)
