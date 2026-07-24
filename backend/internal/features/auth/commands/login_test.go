package commands

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

type mockLoginUserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func newMockLoginUserRepository(user *domain.User) *mockLoginUserRepository {
	repo := &mockLoginUserRepository{users: make(map[string]*domain.User)}
	repo.users[user.ID.String()] = user
	return repo
}

func (r *mockLoginUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFoundError
}

func (r *mockLoginUserRepository) FindByKey(ctx context.Context, key string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.KeyCredential != nil && user.KeyCredential.Validate(key) {
			return user, nil
		}
	}
	return nil, domain.ErrInvalidCredentialsError
}

func (r *mockLoginUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID.String()]; !exists {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID.String()] = user
	return nil
}

type mockLoginSessionRepository struct{}

func (r *mockLoginSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	return nil
}

func testTokenConfig() domain.TokenConfig {
	return domain.TokenConfig{
		AccessTokenSecret:    []byte("test-access-secret"),
		RefreshTokenSecret:   []byte("test-refresh-secret"),
		AccessTokenTTL:       15 * time.Minute,
		RefreshTokenTTL:      24 * time.Hour,
		SessionInactivityTTL: 15 * time.Minute,
		Issuer:               "pococlinic-test",
	}
}

func setupLoginUser(t *testing.T) (*domain.User, string) {
	t.Helper()

	user := domain.NewUser("staff@example.com", "Staff User", domain.RoleStaff)
	key, keyCred, err := domain.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pinCred, err := domain.GeneratePINCredential("1234")
	if err != nil {
		t.Fatalf("generate pin: %v", err)
	}
	user.SetKeyCredential(keyCred, key)
	user.SetPINCredential(pinCred)
	return user, key
}

func TestLoginHandler_AdminLoginRejectsNonAdminRole(t *testing.T) {
	ctx := context.Background()
	user, key := setupLoginUser(t)
	userRepo := newMockLoginUserRepository(user)
	handler := NewLoginHandler(userRepo, &mockLoginSessionRepository{}, testTokenConfig(), nil)

	_, err := handler.HandleAdmin(ctx, AdminLoginCommand{
		Email:     user.Email,
		Key:       key,
		PIN:       "1234",
		IPAddress: "127.0.0.1",
	})
	if !errors.Is(err, domain.ErrInvalidCredentialsError) {
		t.Fatalf("expected invalid credentials for non-admin, got %v", err)
	}
}

func TestLoginHandler_AdminLoginAcceptsAdminRole(t *testing.T) {
	ctx := context.Background()
	user, key := setupLoginUser(t)
	user.Role = domain.RoleAdmin
	userRepo := newMockLoginUserRepository(user)
	handler := NewLoginHandler(userRepo, &mockLoginSessionRepository{}, testTokenConfig(), nil)

	resp, err := handler.HandleAdmin(ctx, AdminLoginCommand{
		Email:     user.Email,
		Key:       key,
		PIN:       "1234",
		IPAddress: "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("expected successful admin login, got %v", err)
	}
	if resp.User.Role != domain.RoleAdmin {
		t.Fatalf("expected admin role, got %s", resp.User.Role)
	}
}

func TestLoginHandler_LocksAccountAfterFiveFailures(t *testing.T) {
	ctx := context.Background()
	user, key := setupLoginUser(t)
	userRepo := newMockLoginUserRepository(user)
	handler := NewLoginHandler(userRepo, &mockLoginSessionRepository{}, testTokenConfig(), nil)

	for attempt := 1; attempt <= 5; attempt++ {
		_, err := handler.HandleStaff(ctx, StaffLoginCommand{
			Key:       key,
			PIN:       "9999",
			IPAddress: "127.0.0.1",
		})
		if !errors.Is(err, domain.ErrInvalidCredentialsError) {
			t.Fatalf("attempt %d: expected invalid credentials, got %v", attempt, err)
		}
	}

	_, err := handler.HandleStaff(ctx, StaffLoginCommand{
		Key:       key,
		PIN:       "9999",
		IPAddress: "127.0.0.1",
	})
	if !errors.Is(err, domain.ErrAccountLockedError) {
		t.Fatalf("expected account locked, got %v", err)
	}

	_, err = handler.HandleStaff(ctx, StaffLoginCommand{
		Key:       key,
		PIN:       "1234",
		IPAddress: "127.0.0.1",
	})
	if !errors.Is(err, domain.ErrAccountLockedError) {
		t.Fatalf("expected account to remain locked with correct pin, got %v", err)
	}
}

func TestLoginHandler_SuccessResetsFailedAttempts(t *testing.T) {
	ctx := context.Background()
	user, key := setupLoginUser(t)
	userRepo := newMockLoginUserRepository(user)
	handler := NewLoginHandler(userRepo, &mockLoginSessionRepository{}, testTokenConfig(), nil)

	for attempt := 1; attempt <= 4; attempt++ {
		if _, err := handler.HandleStaff(ctx, StaffLoginCommand{
			Key: key, PIN: "9999", IPAddress: "127.0.0.1",
		}); !errors.Is(err, domain.ErrInvalidCredentialsError) {
			t.Fatalf("attempt %d: expected invalid credentials, got %v", attempt, err)
		}
	}

	if _, err := handler.HandleStaff(ctx, StaffLoginCommand{
		Key: key, PIN: "1234", IPAddress: "127.0.0.1",
	}); err != nil {
		t.Fatalf("expected successful login, got %v", err)
	}

	updated := userRepo.users[user.ID.String()]
	if updated.FailedAttempts != 0 {
		t.Fatalf("expected failed attempts reset, got %d", updated.FailedAttempts)
	}
	if updated.IsLocked() {
		t.Fatal("expected account unlocked after successful login")
	}
}
