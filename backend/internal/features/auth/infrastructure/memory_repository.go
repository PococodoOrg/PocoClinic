package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
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
		return fmt.Errorf("email %s already registered", user.Email)
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
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.emails[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return r.users[id], nil
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

	if _, exists := r.sessions[session.ID.String()]; !exists {
		return fmt.Errorf("session not found")
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
