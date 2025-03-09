package commands

import (
	"context"
	"fmt"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

// CreateUserCommand represents the command to create a new user
type CreateUserCommand struct {
	Email string      `json:"email" binding:"required,email"`
	Name  string      `json:"name" binding:"required"`
	Role  domain.Role `json:"role" binding:"required"`
}

// CreateUserHandler handles user creation
type CreateUserHandler interface {
	Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, string, error)
}

// createUserHandler implements CreateUserHandler
type createUserHandler struct {
	userRepository domain.CreateUserRepository
}

// NewCreateUserHandler creates a new handler for user creation
func NewCreateUserHandler(repo domain.CreateUserRepository) CreateUserHandler {
	return &createUserHandler{
		userRepository: repo,
	}
}

// Handle processes the create user command
func (h *createUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, string, error) {
	user := domain.NewUser(cmd.Email, cmd.Name, cmd.Role)

	// Generate the initial key and credentials
	key, keyCred, err := domain.GenerateKey()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate key: %w", err)
	}
	user.SetKeyCredential(keyCred)

	// Generate default PIN credentials
	pinCred, err := domain.GeneratePINCredential("0000")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PIN: %w", err)
	}
	user.SetPINCredential(pinCred)

	// Save the user
	err = h.userRepository.Create(ctx, user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	return user, key, nil
}
