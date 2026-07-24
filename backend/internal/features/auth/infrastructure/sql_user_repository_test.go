package infrastructure

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
	"github.com/google/uuid"
)

func openTestDB(t *testing.T) (*database.DB, context.Context) {
	t.Helper()
	ctx := context.Background()
	db, err := database.Connect(ctx, filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db, ctx
}

func TestSQLUserRepository_CreateGetListSearch(t *testing.T) {
	db, ctx := openTestDB(t)
	repo := NewSQLUserRepository(db)

	plaintextKey, keyCred, err := domain.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pinCred, err := domain.NewCredential("1234")
	if err != nil {
		t.Fatalf("pin credential: %v", err)
	}

	user := domain.NewUser("nurse@clinic.test", "Nina Nurse", domain.RoleNurse)
	user.SetKeyCredential(keyCred, plaintextKey)
	user.SetPINCredential(pinCred)
	user.MustChangePIN = true

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("create (16 placeholders): %v", err)
	}

	got, err := repo.GetByEmail(ctx, "nurse@clinic.test")
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if got.Name != "Nina Nurse" || got.Role != domain.RoleNurse {
		t.Fatalf("unexpected user: %+v", got)
	}
	if !got.MustChangePIN || !got.IsActive {
		t.Fatalf("bool columns not round-tripped: mustChange=%v active=%v", got.MustChangePIN, got.IsActive)
	}

	users, total, err := repo.ListPaginated(ctx, 1, 10, "nina")
	if err != nil {
		t.Fatalf("list search: %v", err)
	}
	if total != 1 || len(users) != 1 {
		t.Fatalf("search expected 1 user, total=%d len=%d", total, len(users))
	}

	got.Name = "Nina Updated"
	got.FailedAttempts = 2
	got.UpdatedAt = time.Now().UTC()
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update ($10+ placeholders): %v", err)
	}

	updated, err := repo.GetByID(ctx, got.ID.String())
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if updated.Name != "Nina Updated" || updated.FailedAttempts != 2 {
		t.Fatalf("update not persisted: %+v", updated)
	}
}

func TestSQLUserRepository_SeedDefaultAdmin(t *testing.T) {
	db, ctx := openTestDB(t)
	repo := NewSQLUserRepository(db)

	key, err := SeedDefaultAdmin(ctx, repo)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	if key == "" {
		t.Fatal("expected one-time badge key from first seed")
	}

	again, err := SeedDefaultAdmin(ctx, repo)
	if err != nil {
		t.Fatalf("second seed: %v", err)
	}
	if again != "" {
		t.Fatal("expected empty key when admin already exists")
	}

	admin, err := repo.GetByEmail(ctx, DefaultAdminEmail)
	if err != nil {
		t.Fatalf("get admin: %v", err)
	}
	if admin.Role != domain.RoleAdmin {
		t.Fatalf("expected admin role, got %s", admin.Role)
	}
	if admin.ID == uuid.Nil {
		t.Fatal("admin id missing")
	}
}
