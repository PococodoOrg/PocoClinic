package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

// SQLUserRepository persists users in SQLite.
type SQLUserRepository struct {
	pool *database.DB
}

// NewSQLUserRepository creates a SQL-backed user repository.
func NewSQLUserRepository(pool *database.DB) *SQLUserRepository {
	return &SQLUserRepository{pool: pool}
}

// Create adds a new user.
func (r *SQLUserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.KeyCredential == nil || user.PINCredential == nil {
		return fmt.Errorf("user credentials are required")
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (
			id, email, name, role, key_hash, key_salt, pin_hash, pin_salt, key_lookup,
			failed_attempts, locked_until, last_login, must_change_pin, is_active, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`,
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		user.KeyCredential.Hash,
		user.KeyCredential.Salt,
		user.PINCredential.Hash,
		user.PINCredential.Salt,
		nullIfEmpty(user.KeyLookup),
		user.FailedAttempts,
		user.LockedUntil,
		user.LastLogin,
		user.MustChangePIN,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrEmailTakenError(user.Email)
		}
		return err
	}
	return nil
}

// Update modifies an existing user.
func (r *SQLUserRepository) Update(ctx context.Context, user *domain.User) error {
	if user.KeyCredential == nil || user.PINCredential == nil {
		return fmt.Errorf("user credentials are required")
	}

	tag, err := r.pool.Exec(ctx, `
		UPDATE users SET
			email = $2,
			name = $3,
			role = $4,
			key_hash = $5,
			key_salt = $6,
			pin_hash = $7,
			pin_salt = $8,
			key_lookup = $9,
			failed_attempts = $10,
			locked_until = $11,
			last_login = $12,
			must_change_pin = $13,
			is_active = $14,
			updated_at = $15
		WHERE id = $1
	`,
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		user.KeyCredential.Hash,
		user.KeyCredential.Salt,
		user.PINCredential.Hash,
		user.PINCredential.Salt,
		nullIfEmpty(user.KeyLookup),
		user.FailedAttempts,
		user.LockedUntil,
		user.LastLogin,
		user.MustChangePIN,
		user.IsActive,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFoundError
	}
	return nil
}

// Delete removes a user.
func (r *SQLUserRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFoundError
	}
	return nil
}

// GetByID retrieves a user by ID.
func (r *SQLUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, userSelectQuery+` WHERE id = $1`, id)
	return scanUser(row)
}

// GetByEmail retrieves a user by email.
func (r *SQLUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, userSelectQuery+` WHERE email = $1`, email)
	user, err := scanUser(row)
	if errors.Is(err, database.ErrNoRows) {
		return nil, domain.ErrUserNotFoundError
	}
	return user, err
}

// FindByKey retrieves a user whose badge key matches the provided value.
func (r *SQLUserRepository) FindByKey(ctx context.Context, key string) (*domain.User, error) {
	lookup := domain.KeyLookup(key)
	row := r.pool.QueryRow(ctx, userSelectQuery+` WHERE key_lookup = $1`, lookup)
	user, err := scanUser(row)
	if err == nil {
		if user.KeyCredential != nil && user.KeyCredential.Validate(key) {
			return user, nil
		}
		return nil, domain.ErrInvalidCredentialsError
	}
	if !errors.Is(err, database.ErrNoRows) {
		return nil, err
	}

	return nil, domain.ErrInvalidCredentialsError
}

// ListPaginated returns a paginated list of users with optional search.
func (r *SQLUserRepository) ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*domain.User, int64, error) {
	search = strings.TrimSpace(search)
	searchPattern := "%" + strings.ToLower(search) + "%"

	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM users`
	countArgs := []any{}
	if search != "" {
		countQuery += `
			WHERE lower(name) LIKE $1
			   OR lower(email) LIKE $1
			   OR lower(role) LIKE $1
		`
		countArgs = append(countArgs, searchPattern)
	}
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	listQuery := userPublicSelectQuery
	listArgs := []any{}
	if search != "" {
		listQuery += `
			WHERE lower(name) LIKE $1
			   OR lower(email) LIKE $1
			   OR lower(role) LIKE $1
		`
		listArgs = append(listArgs, searchPattern)
	}
	offset := (page - 1) * pageSize
	listQuery += fmt.Sprintf(` ORDER BY name LIMIT $%d OFFSET $%d`, len(listArgs)+1, len(listArgs)+2)
	listArgs = append(listArgs, pageSize, offset)

	rows, err := r.pool.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user, err := scanUserPublic(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, totalCount, rows.Err()
}

// CountByRole returns how many users have the given role.
func (r *SQLUserRepository) CountByRole(ctx context.Context, role domain.Role) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE role = $1`, role).Scan(&count)
	return count, err
}

const userPublicSelectQuery = `
	SELECT id, email, name, role, last_login, must_change_pin, is_active, locked_until, created_at, updated_at
	FROM users
`

func scanUserPublic(row database.Row) (*domain.User, error) {
	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		database.NullTime(&user.LastLogin),
		&user.MustChangePIN,
		&user.IsActive,
		database.NullTime(&user.LockedUntil),
		database.Time(&user.CreatedAt),
		database.Time(&user.UpdatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, domain.ErrUserNotFoundError
		}
		return nil, err
	}
	return user, nil
}

const userSelectQuery = `
	SELECT id, email, name, role, key_hash, key_salt, pin_hash, pin_salt, key_lookup,
	       failed_attempts, locked_until, last_login, must_change_pin, is_active, created_at, updated_at
	FROM users
`

func scanUser(row database.Row) (*domain.User, error) {
	user := &domain.User{}
	keyHash := make([]byte, 0)
	keySalt := make([]byte, 0)
	pinHash := make([]byte, 0)
	pinSalt := make([]byte, 0)
	var keyLookup *string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&keyHash,
		&keySalt,
		&pinHash,
		&pinSalt,
		&keyLookup,
		&user.FailedAttempts,
		database.NullTime(&user.LockedUntil),
		database.NullTime(&user.LastLogin),
		&user.MustChangePIN,
		&user.IsActive,
		database.Time(&user.CreatedAt),
		database.Time(&user.UpdatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, domain.ErrUserNotFoundError
		}
		return nil, err
	}

	user.KeyCredential = &domain.Credential{Hash: keyHash, Salt: keySalt}
	user.PINCredential = &domain.Credential{Hash: pinHash, Salt: pinSalt}
	if keyLookup != nil {
		user.KeyLookup = *keyLookup
	}
	return user, nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key") || strings.Contains(message, "unique constraint")
}
