// Package domain provides the core domain models and business logic for authentication.
// It implements a secure two-factor authentication system using a 64-bit key and PIN,
// as specified in ADR-0002. The package handles user management, credential validation,
// and account security measures such as rate limiting and account locking.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user's role in the system
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleDoctor  Role = "doctor"
	RoleNurse   Role = "nurse"
	RoleStaff   Role = "staff"
	RolePatient Role = "patient"
)

// User represents a user in the system
type User struct {
	ID             uuid.UUID   `json:"id"`
	Email          string      `json:"email"`
	Name           string      `json:"name"`
	Role           Role        `json:"role"`
	KeyCredential  *Credential `json:"-"`
	PINCredential  *Credential `json:"-"`
	FailedAttempts int         `json:"-"`
	LockedUntil    *time.Time  `json:"-"`
	LastLogin      *time.Time  `json:"lastLogin,omitempty"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

// NewUser creates a new user with a generated ID and timestamps
func NewUser(email, name string, role Role) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetKeyCredential sets the user's key credential
func (u *User) SetKeyCredential(cred *Credential) {
	u.KeyCredential = cred
	u.UpdatedAt = time.Now()
}

// SetPINCredential sets the user's PIN credential
func (u *User) SetPINCredential(cred *Credential) {
	u.PINCredential = cred
	u.UpdatedAt = time.Now()
}

// ValidateCredentials validates both key and PIN
func (u *User) ValidateCredentials(key, pin string) bool {
	if u.KeyCredential == nil || u.PINCredential == nil {
		return false
	}
	return u.KeyCredential.Validate(key) && u.PINCredential.Validate(pin)
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// RecordFailedAttempt records a failed authentication attempt
func (u *User) RecordFailedAttempt() {
	u.FailedAttempts++
	if u.FailedAttempts >= 5 {
		lockUntil := time.Now().Add(15 * time.Minute)
		u.LockedUntil = &lockUntil
	}
	u.UpdatedAt = time.Now()
}

// ResetFailedAttempts resets the failed attempts counter
func (u *User) ResetFailedAttempts() {
	u.FailedAttempts = 0
	u.LockedUntil = nil
	u.UpdatedAt = time.Now()
}

// RecordLogin records a successful login
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.ResetFailedAttempts()
	u.UpdatedAt = now
}
