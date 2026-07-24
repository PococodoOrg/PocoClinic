package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dksch/pococlinic/internal/features/auth/commands"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

// SeedDefaultAdmin creates a default admin user when one does not exist yet.
func SeedDefaultAdmin(ctx context.Context, repo domain.UserRepository) (string, error) {
	_, err := repo.GetByEmail(ctx, "admin@pococlinic.local")
	if err == nil {
		return "", nil
	}

	var authErr *domain.AuthError
	if !errors.As(err, &authErr) || authErr.Code != domain.ErrUserNotFound {
		return "", fmt.Errorf("check default admin: %w", err)
	}

	handler := commands.NewCreateUserHandler(repo)
	_, key, err := handler.Handle(ctx, commands.CreateUserCommand{
		Email: "admin@pococlinic.local",
		Name:  "System Administrator",
		Role:  domain.RoleAdmin,
	})
	if err != nil {
		return "", fmt.Errorf("failed to seed default admin: %w", err)
	}

	return key, nil
}

const DefaultAdminEmail = "admin@pococlinic.local"

// WriteBootstrapCredentialsFile stores one-time setup credentials on disk (never log the contents).
func WriteBootstrapCredentialsFile(dataDir, email, badgeKey string) (string, error) {
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return "", fmt.Errorf("create bootstrap credentials directory: %w", err)
	}

	path := filepath.Join(dataDir, "bootstrap-admin-once.txt")
	content := fmt.Sprintf(`PocoClinic default administrator credentials (first boot only)

Email: %s
Temporary PIN: 0000 (required to change on first login)
Badge key: %s

Store these securely, then delete this file.
`, email, badgeKey)

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return "", fmt.Errorf("write bootstrap credentials file: %w", err)
	}

	return path, nil
}

// BootstrapCredentialsPath returns the one-time bootstrap credentials file path.
func BootstrapCredentialsPath(dataDir string) string {
	return filepath.Join(dataDir, "bootstrap-admin-once.txt")
}

// BootstrapCredentialsExist reports whether the bootstrap file is still on disk.
func BootstrapCredentialsExist(dataDir string) bool {
	_, err := os.Stat(BootstrapCredentialsPath(dataDir))
	return err == nil
}

// RemoveBootstrapCredentialsFile deletes the bootstrap credentials file if present.
func RemoveBootstrapCredentialsFile(dataDir string) error {
	if err := os.Remove(BootstrapCredentialsPath(dataDir)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove bootstrap credentials file: %w", err)
	}
	return nil
}
