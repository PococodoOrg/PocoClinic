package commands

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

type mockRefreshUserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func (r *mockRefreshUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFoundError
	}
	return user, nil
}

type mockRefreshSessionRepository struct {
	sessions map[string]*domain.Session
	mu       sync.RWMutex
}

func (r *mockRefreshSessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, session := range r.sessions {
		if session.RefreshToken == token {
			return session, nil
		}
	}
	return nil, domain.ErrInvalidTokenError
}

func (r *mockRefreshSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID.String()] = session
	return nil
}

func setupRefreshFixture(t *testing.T) (*domain.User, string, *mockRefreshSessionRepository, *mockRefreshUserRepository) {
	t.Helper()

	user := domain.NewUser("refresh@clinic.test", "Refresh User", domain.RoleStaff)
	key, keyCred, err := domain.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pinCred, err := domain.NewCredential("1234")
	if err != nil {
		t.Fatalf("pin: %v", err)
	}
	user.SetKeyCredential(keyCred, key)
	user.SetPINCredential(pinCred)
	user.IsActive = true

	userRepo := &mockRefreshUserRepository{users: map[string]*domain.User{user.ID.String(): user}}
	tokenCfg := testTokenConfig()
	sessionRepo := &mockRefreshSessionRepository{sessions: make(map[string]*domain.Session)}
	session := domain.NewSession(user.ID, "ua", "127.0.0.1", time.Now().UTC().Add(24*time.Hour))
	_, refreshToken, err := session.GenerateTokens(user, tokenCfg)
	if err != nil {
		t.Fatalf("generate tokens: %v", err)
	}
	sessionRepo.sessions[session.ID.String()] = session

	return user, refreshToken, sessionRepo, userRepo
}

func TestRefreshTokenHandler_Success(t *testing.T) {
	user, refreshToken, sessionRepo, userRepo := setupRefreshFixture(t)
	handler := NewRefreshTokenHandler(userRepo, sessionRepo, testTokenConfig(), nil)

	resp, err := handler.Handle(context.Background(), RefreshTokenCommand{
		RefreshToken: refreshToken,
		UserAgent:    "test",
		IPAddress:    "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("expected rotated tokens")
	}
	if resp.User == nil || resp.User.ID != user.ID {
		t.Fatalf("unexpected user: %+v", resp.User)
	}
}

func TestRefreshTokenHandler_EmptyToken(t *testing.T) {
	handler := NewRefreshTokenHandler(
		&mockRefreshUserRepository{users: map[string]*domain.User{}},
		&mockRefreshSessionRepository{sessions: make(map[string]*domain.Session)},
		testTokenConfig(),
		nil,
	)
	_, err := handler.Handle(context.Background(), RefreshTokenCommand{})
	if err != domain.ErrInvalidTokenError {
		t.Fatalf("expected ErrInvalidTokenError, got %v", err)
	}
}

func TestRefreshTokenHandler_ExpiredSession(t *testing.T) {
	user := domain.NewUser("expired@clinic.test", "Expired", domain.RoleStaff)
	user.IsActive = true
	userRepo := &mockRefreshUserRepository{users: map[string]*domain.User{user.ID.String(): user}}

	tokenCfg := testTokenConfig()
	sessionRepo := &mockRefreshSessionRepository{sessions: make(map[string]*domain.Session)}
	session := domain.NewSession(user.ID, "ua", "127.0.0.1", time.Now().UTC().Add(-time.Hour))
	_, refreshToken, err := session.GenerateTokens(user, tokenCfg)
	if err != nil {
		t.Fatalf("tokens: %v", err)
	}
	sessionRepo.sessions[session.ID.String()] = session

	handler := NewRefreshTokenHandler(userRepo, sessionRepo, tokenCfg, nil)
	_, err = handler.Handle(context.Background(), RefreshTokenCommand{RefreshToken: refreshToken})
	if err != domain.ErrSessionExpiredError {
		t.Fatalf("expected ErrSessionExpiredError, got %v", err)
	}
}

func TestRefreshTokenHandler_InactiveUser(t *testing.T) {
	user, refreshToken, sessionRepo, userRepo := setupRefreshFixture(t)
	user.IsActive = false
	userRepo.users[user.ID.String()] = user

	handler := NewRefreshTokenHandler(userRepo, sessionRepo, testTokenConfig(), nil)
	_, err := handler.Handle(context.Background(), RefreshTokenCommand{RefreshToken: refreshToken})
	if err != domain.ErrSessionExpiredError {
		t.Fatalf("expected ErrSessionExpiredError for inactive user, got %v", err)
	}
}

func TestRefreshTokenHandler_MissingSession(t *testing.T) {
	user := domain.NewUser("missing@clinic.test", "Missing", domain.RoleStaff)
	user.IsActive = true
	userRepo := &mockRefreshUserRepository{users: map[string]*domain.User{user.ID.String(): user}}
	tokenCfg := testTokenConfig()
	session := domain.NewSession(user.ID, "ua", "127.0.0.1", time.Now().UTC().Add(time.Hour))
	_, refreshToken, err := session.GenerateTokens(user, tokenCfg)
	if err != nil {
		t.Fatalf("tokens: %v", err)
	}

	handler := NewRefreshTokenHandler(
		userRepo,
		&mockRefreshSessionRepository{sessions: make(map[string]*domain.Session)},
		tokenCfg,
		nil,
	)
	_, err = handler.Handle(context.Background(), RefreshTokenCommand{RefreshToken: refreshToken})
	if err != domain.ErrInvalidTokenError {
		t.Fatalf("expected ErrInvalidTokenError, got %v", err)
	}
}