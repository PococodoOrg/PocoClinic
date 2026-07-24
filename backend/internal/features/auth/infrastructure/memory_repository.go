package infrastructure

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
)

// MemoryUserRepository is a simple in-memory implementation of the user repository
type MemoryUserRepository struct {
	users  map[string]*domain.User // key: user ID
	emails map[string]string       // key: email, value: user ID
	mu     sync.RWMutex
}

// NewMemoryUserRepository creates a new in-memory user repository
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:  make(map[string]*domain.User),
		emails: make(map[string]string),
	}
}

// Create adds a new user
func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.emails[user.Email]; exists {
		return domain.ErrEmailTakenError(user.Email)
	}

	r.users[user.ID.String()] = user
	r.emails[user.Email] = user.ID.String()
	return nil
}

// Update modifies an existing user
func (r *MemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID.String()]; !exists {
		return fmt.Errorf("user not found")
	}

	r.users[user.ID.String()] = user
	return nil
}

// Delete removes a user
func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)
	delete(r.emails, user.Email)
	return nil
}

// GetByID retrieves a user by ID
func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFoundError
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.emails[email]
	if !exists {
		return nil, domain.ErrUserNotFoundError
	}

	return r.users[id], nil
}

// FindByKey retrieves a user whose badge key matches the provided value.
func (r *MemoryUserRepository) FindByKey(ctx context.Context, key string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lookup := domain.KeyLookup(key)
	for _, user := range r.users {
		if user.KeyLookup != "" && user.KeyLookup == lookup {
			if user.KeyCredential != nil && user.KeyCredential.Validate(key) {
				return user, nil
			}
			return nil, domain.ErrInvalidCredentialsError
		}
	}

	return nil, domain.ErrInvalidCredentialsError
}

// ListPaginated returns a paginated list of users with optional search.
func (r *MemoryUserRepository) ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*domain.User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		if matchesUserSearch(user, search) {
			all = append(all, user)
		}
	}

	sort.Slice(all, func(i, j int) bool {
		return strings.ToLower(all[i].Name) < strings.ToLower(all[j].Name)
	})

	totalCount := int64(len(all))
	start := (page - 1) * pageSize
	if start >= len(all) {
		return []*domain.User{}, totalCount, nil
	}

	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}

	return all[start:end], totalCount, nil
}

// CountByRole returns how many users have the given role.
func (r *MemoryUserRepository) CountByRole(ctx context.Context, role domain.Role) (int64, error) {
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

func matchesUserSearch(user *domain.User, search string) bool {
	if search == "" {
		return true
	}
	search = strings.ToLower(strings.TrimSpace(search))
	return strings.Contains(strings.ToLower(user.Name), search) ||
		strings.Contains(strings.ToLower(user.Email), search) ||
		strings.Contains(strings.ToLower(string(user.Role)), search)
}

// MemorySessionRepository is a simple in-memory implementation of the session repository
type MemorySessionRepository struct {
	sessions map[string]*domain.Session // key: session ID
	tokens   map[string]string          // key: refresh token, value: session ID
	mu       sync.RWMutex
}

// NewMemorySessionRepository creates a new in-memory session repository
func NewMemorySessionRepository() *MemorySessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[string]*domain.Session),
		tokens:   make(map[string]string),
	}
}

// Create adds a new session
func (r *MemorySessionRepository) Create(ctx context.Context, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID.String()] = session
	if session.RefreshToken != "" {
		r.tokens[session.RefreshToken] = session.ID.String()
	}
	return nil
}

// Update modifies an existing session
func (r *MemorySessionRepository) Update(ctx context.Context, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.sessions[session.ID.String()]
	if !exists {
		return fmt.Errorf("session not found")
	}

	if existing.RefreshToken != "" && existing.RefreshToken != session.RefreshToken {
		delete(r.tokens, existing.RefreshToken)
	}
	if session.RefreshToken != "" {
		r.tokens[session.RefreshToken] = session.ID.String()
	}

	r.sessions[session.ID.String()] = session
	return nil
}

// Delete removes a session
func (r *MemorySessionRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[id]
	if !exists {
		return fmt.Errorf("session not found")
	}

	delete(r.sessions, id)
	if session.RefreshToken != "" {
		delete(r.tokens, session.RefreshToken)
	}
	return nil
}

// DeleteByUserID removes all sessions for a user.
func (r *MemorySessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.DeleteByUserIDExcept(ctx, userID, "")
}

// DeleteByUserIDExcept removes all sessions for a user except the given session ID.
func (r *MemorySessionRepository) DeleteByUserIDExcept(ctx context.Context, userID, exceptSessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, session := range r.sessions {
		if session.UserID.String() != userID {
			continue
		}
		if exceptSessionID != "" && id == exceptSessionID {
			continue
		}
		delete(r.sessions, id)
		if session.RefreshToken != "" {
			delete(r.tokens, session.RefreshToken)
		}
	}
	return nil
}

// GetByID retrieves a session by ID
func (r *MemorySessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}

// GetByRefreshToken retrieves a session by refresh token
func (r *MemorySessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.tokens[token]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	return r.sessions[id], nil
}

// DeleteExpired removes all expired sessions
func (r *MemorySessionRepository) DeleteExpired(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, session := range r.sessions {
		if session.IsExpired() {
			delete(r.sessions, id)
			if session.RefreshToken != "" {
				delete(r.tokens, session.RefreshToken)
			}
		}
	}
	return nil
}
