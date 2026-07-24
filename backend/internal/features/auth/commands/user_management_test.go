package commands

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

type mockUserManagementRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func (r *mockUserManagementRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.users {
		if existing.Email == user.Email {
			return fmt.Errorf("duplicate email")
		}
	}
	r.users[user.ID.String()] = user
	return nil
}

func (r *mockUserManagementRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID.String()] = user
	return nil
}

func (r *mockUserManagementRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFoundError
}

func (r *mockUserManagementRepository) FindByKey(ctx context.Context, key string) (*domain.User, error) {
	return nil, domain.ErrUserNotFoundError
}

func (r *mockUserManagementRepository) ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*domain.User, int64, error) {
	return nil, 0, nil
}

func (r *mockUserManagementRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFoundError
	}
	return user, nil
}

func (r *mockUserManagementRepository) CountByRole(ctx context.Context, role domain.Role) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, user := range r.users {
		if user.Role == role {
			count++
		}
	}
	return count, nil
}

func (r *mockUserManagementRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}

type mockDeleteSessionRepository struct{}

func (r *mockDeleteSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	return nil
}
func (r *mockDeleteSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	return nil
}
func (r *mockDeleteSessionRepository) Delete(ctx context.Context, id string) error { return nil }
func (r *mockDeleteSessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	return nil, domain.ErrInvalidTokenError
}
func (r *mockDeleteSessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	return nil, domain.ErrInvalidTokenError
}
func (r *mockDeleteSessionRepository) DeleteExpired(ctx context.Context) error { return nil }
func (r *mockDeleteSessionRepository) DeleteByUserIDExcept(ctx context.Context, userID, exceptSessionID string) error {
	return nil
}

func (r *mockDeleteSessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return nil
}

func TestCreateUserHandler_RejectsInvalidRole(t *testing.T) {
	handler := NewCreateUserHandler(&mockUserManagementRepository{users: map[string]*domain.User{}})
	_, _, err := handler.Handle(context.Background(), CreateUserCommand{
		Email: "bad@clinic.test",
		Name:  "Bad Role",
		Role:  domain.Role("superuser"),
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCreateUserHandler_CreatesStaffWithDefaultPIN(t *testing.T) {
	repo := &mockUserManagementRepository{users: map[string]*domain.User{}}
	handler := NewCreateUserHandler(repo)

	user, key, err := handler.Handle(context.Background(), CreateUserCommand{
		Email: "nurse@clinic.test",
		Name:  "Clinic Nurse",
		Role:  domain.RoleNurse,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if key == "" {
		t.Fatal("expected badge key")
	}
	if !user.MustChangePIN {
		t.Fatal("expected MustChangePIN=true for new users")
	}
	if user.Role != domain.RoleNurse {
		t.Fatalf("unexpected role %q", user.Role)
	}
	if _, ok := repo.users[user.ID.String()]; !ok {
		t.Fatal("user not persisted")
	}
}

func TestDeleteUserHandler_CannotDeleteSelf(t *testing.T) {
	repo := &mockUserManagementRepository{users: map[string]*domain.User{}}
	handler := NewDeleteUserHandler(repo, &mockDeleteSessionRepository{}, nil)

	user := domain.NewUser("self@clinic.test", "Self", domain.RoleAdmin)
	repo.users[user.ID.String()] = user

	err := handler.Handle(context.Background(), DeleteUserCommand{
		UserID:  user.ID.String(),
		ActorID: user.ID.String(),
	})
	if err != domain.ErrCannotDeleteSelfError {
		t.Fatalf("expected cannot delete self, got %v", err)
	}
}

func TestDeleteUserHandler_CannotDeleteLastAdmin(t *testing.T) {
	repo := &mockUserManagementRepository{users: map[string]*domain.User{}}
	handler := NewDeleteUserHandler(repo, &mockDeleteSessionRepository{}, nil)

	admin := domain.NewUser("admin@clinic.test", "Admin", domain.RoleAdmin)
	actor := domain.NewUser("actor@clinic.test", "Actor", domain.RoleNurse)
	repo.users[admin.ID.String()] = admin
	repo.users[actor.ID.String()] = actor

	err := handler.Handle(context.Background(), DeleteUserCommand{
		UserID:  admin.ID.String(),
		ActorID: actor.ID.String(),
	})
	if err != domain.ErrLastAdminError {
		t.Fatalf("expected last admin error, got %v", err)
	}
}
