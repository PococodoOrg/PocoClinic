package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	name := "Test User"
	role := RoleDoctor

	user := NewUser(email, name, role)

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
	if user.Name != name {
		t.Errorf("Expected name %s, got %s", name, user.Name)
	}
	if user.Role != role {
		t.Errorf("Expected role %s, got %s", role, user.Role)
	}
	if user.ID == uuid.Nil {
		t.Error("Expected non-nil UUID")
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("Expected non-zero UpdatedAt")
	}
}

func TestCredentials(t *testing.T) {
	user := NewUser("test@example.com", "Test User", RoleDoctor)

	// Test key credential
	key, keyCred, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	user.SetKeyCredential(keyCred)
	if user.KeyCredential != keyCred {
		t.Error("Expected key credential to be set")
	}

	// Test PIN credential
	pin := "1234"
	pinCred, err := GeneratePINCredential(pin)
	if err != nil {
		t.Fatalf("Failed to generate PIN credential: %v", err)
	}
	user.SetPINCredential(pinCred)
	if user.PINCredential != pinCred {
		t.Error("Expected PIN credential to be set")
	}

	// Test credential validation
	if !user.ValidateCredentials(key, pin) {
		t.Error("Expected credential validation to succeed")
	}
}

func TestAccountLocking(t *testing.T) {
	user := NewUser("test@example.com", "Test User", RoleDoctor)

	// Test initial state
	if user.IsLocked() {
		t.Error("Expected new user to be unlocked")
	}

	// Test failed attempts
	for i := 0; i < 4; i++ {
		user.RecordFailedAttempt()
		if user.IsLocked() {
			t.Errorf("Expected user to be unlocked after %d attempts", i+1)
		}
	}

	// Test locking after 5 attempts
	user.RecordFailedAttempt()
	if !user.IsLocked() {
		t.Error("Expected user to be locked after 5 attempts")
	}

	// Test reset
	user.ResetFailedAttempts()
	if user.IsLocked() {
		t.Error("Expected user to be unlocked after reset")
	}
	if user.FailedAttempts != 0 {
		t.Error("Expected failed attempts to be reset to 0")
	}
}

func TestRecordLogin(t *testing.T) {
	user := NewUser("test@example.com", "Test User", RoleDoctor)

	// Record some failed attempts
	user.RecordFailedAttempt()
	user.RecordFailedAttempt()

	// Record login
	originalUpdatedAt := user.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference
	user.RecordLogin()

	if user.LastLogin == nil {
		t.Error("Expected LastLogin to be set")
	}
	if user.FailedAttempts != 0 {
		t.Error("Expected failed attempts to be reset")
	}
	if user.IsLocked() {
		t.Error("Expected user to be unlocked")
	}
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}
