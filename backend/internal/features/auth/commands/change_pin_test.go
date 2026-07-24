package commands

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

type mockChangePINRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func newMockChangePINRepository(user *domain.User) *mockChangePINRepository {
	repo := &mockChangePINRepository{users: make(map[string]*domain.User)}
	repo.users[user.ID.String()] = user
	return repo
}

func (r *mockChangePINRepository) Create(ctx context.Context, user *domain.User) error {
	return fmt.Errorf("not implemented")
}

func (r *mockChangePINRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID.String()] = user
	return nil
}

func (r *mockChangePINRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

func (r *mockChangePINRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFoundError
	}
	return user, nil
}

func (r *mockChangePINRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *mockChangePINRepository) FindByKey(ctx context.Context, key string) (*domain.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *mockChangePINRepository) ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*domain.User, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (r *mockChangePINRepository) CountByRole(ctx context.Context, role domain.Role) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

type mockChangePINSessionRepository struct{}

func (r *mockChangePINSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	return nil
}

func (r *mockChangePINSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	return nil
}

func (r *mockChangePINSessionRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *mockChangePINSessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *mockChangePINSessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *mockChangePINSessionRepository) DeleteExpired(ctx context.Context) error {
	return nil
}

func (r *mockChangePINSessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return nil
}

func (r *mockChangePINSessionRepository) DeleteByUserIDExcept(ctx context.Context, userID, exceptSessionID string) error {
	return nil
}

func TestChangePINHandler(t *testing.T) {
	user := domain.NewUser("staff@example.com", "Staff User", domain.RoleStaff)
	key, keyCred, err := domain.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	user.SetKeyCredential(keyCred, key)
	pinCred, err := domain.GeneratePINCredential("1234")
	if err != nil {
		t.Fatalf("generate pin: %v", err)
	}
	user.SetPINCredential(pinCred)

	repo := newMockChangePINRepository(user)
	handler := NewChangePINHandler(repo, &mockChangePINSessionRepository{}, nil)

	t.Run("changes pin with valid current pin", func(t *testing.T) {
		updated, err := handler.Handle(context.Background(), ChangePINCommand{
			UserID:     user.ID.String(),
			CurrentPIN: "1234",
			NewPIN:     "5678",
		})
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if updated.MustChangePIN {
			t.Fatal("expected mustChangePIN to be cleared")
		}
		if !updated.ValidatePIN("5678") {
			t.Fatal("expected new PIN to validate")
		}
	})

	t.Run("rejects incorrect current pin", func(t *testing.T) {
		_, err := handler.Handle(context.Background(), ChangePINCommand{
			UserID:     user.ID.String(),
			CurrentPIN: "0000",
			NewPIN:     "9999",
		})
		if err != domain.ErrInvalidPINError {
			t.Fatalf("expected invalid pin error, got %v", err)
		}
	})

	t.Run("rejects same pin", func(t *testing.T) {
		_, err := handler.Handle(context.Background(), ChangePINCommand{
			UserID:     user.ID.String(),
			CurrentPIN: "5678",
			NewPIN:     "5678",
		})
		if err != domain.ErrWeakPINError {
			t.Fatalf("expected weak pin error, got %v", err)
		}
	})

	t.Run("locks account after five failed current pin attempts", func(t *testing.T) {
		lockUser := domain.NewUser("lock@example.com", "Lock User", domain.RoleStaff)
		key, keyCred, err := domain.GenerateKey()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		lockUser.SetKeyCredential(keyCred, key)
		pinCred, err := domain.GeneratePINCredential("1234")
		if err != nil {
			t.Fatalf("generate pin: %v", err)
		}
		lockUser.SetPINCredential(pinCred)
		lockRepo := newMockChangePINRepository(lockUser)
		lockHandler := NewChangePINHandler(lockRepo, &mockChangePINSessionRepository{}, nil)

		for attempt := 1; attempt <= 4; attempt++ {
			_, err := lockHandler.Handle(context.Background(), ChangePINCommand{
				UserID:     lockUser.ID.String(),
				CurrentPIN: "9999",
				NewPIN:     "5678",
			})
			if err != domain.ErrInvalidPINError {
				t.Fatalf("attempt %d: expected invalid pin, got %v", attempt, err)
			}
		}

		_, err = lockHandler.Handle(context.Background(), ChangePINCommand{
			UserID:     lockUser.ID.String(),
			CurrentPIN: "9999",
			NewPIN:     "5678",
		})
		if err != domain.ErrAccountLockedError {
			t.Fatalf("expected account locked on 5th attempt, got %v", err)
		}

		_, err = lockHandler.Handle(context.Background(), ChangePINCommand{
			UserID:     lockUser.ID.String(),
			CurrentPIN: "9999",
			NewPIN:     "5678",
		})
		if err != domain.ErrAccountLockedError {
			t.Fatalf("expected account locked, got %v", err)
		}
	})
}
