package infrastructure

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

func TestSQLSessionRepository_CreateGetByRefresh(t *testing.T) {
	ctx := context.Background()
	db, err := database.Connect(ctx, filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	userRepo := NewSQLUserRepository(db)
	key, keyCred, err := domain.GenerateKey()
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	pinCred, err := domain.NewCredential("1234")
	if err != nil {
		t.Fatalf("pin: %v", err)
	}
	user := domain.NewUser("session@clinic.test", "Session User", domain.RoleNurse)
	user.SetKeyCredential(keyCred, key)
	user.SetPINCredential(pinCred)
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	sessionRepo := NewSQLSessionRepository(db)
	session := domain.NewSession(user.ID, "ua", "127.0.0.1", time.Now().UTC().Add(24*time.Hour))
	session.RefreshToken = "hashed-refresh-token"
	if err := sessionRepo.Create(ctx, session); err != nil {
		t.Fatalf("create session: %v", err)
	}

	got, err := sessionRepo.GetByRefreshToken(ctx, "hashed-refresh-token")
	if err != nil {
		t.Fatalf("get by refresh: %v", err)
	}
	if got.UserID != user.ID {
		t.Fatalf("user id mismatch: %v", got.UserID)
	}
	if got.ExpiresAt.Before(time.Now().UTC()) {
		t.Fatalf("expires_at not scanned: %v", got.ExpiresAt)
	}
}
