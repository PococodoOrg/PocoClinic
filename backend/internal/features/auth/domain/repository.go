package domain

import (
	"context"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Session, error)
	GetByRefreshToken(ctx context.Context, token string) (*Session, error)
	DeleteExpired(ctx context.Context) error
}

// CreateUserRepository defines the minimal interface for user creation
type CreateUserRepository interface {
	Create(ctx context.Context, user *User) error
}

// ValidateUserRepository defines the minimal interface for user validation
type ValidateUserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// CreateSessionRepository defines the minimal interface for session creation
type CreateSessionRepository interface {
	Create(ctx context.Context, session *Session) error
}

// RefreshSessionRepository defines the minimal interface for session refresh
type RefreshSessionRepository interface {
	GetByRefreshToken(ctx context.Context, token string) (*Session, error)
	Update(ctx context.Context, session *Session) error
}
