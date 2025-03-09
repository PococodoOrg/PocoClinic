package queries

import (
	"context"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

// GetUserQuery represents the query to get a user
type GetUserQuery struct {
	ID string `json:"id" binding:"required"`
}

// GetUserHandler handles user retrieval
type GetUserHandler interface {
	Handle(ctx context.Context, query GetUserQuery) (*domain.User, error)
}

// getUserHandler implements GetUserHandler
type getUserHandler struct {
	userRepository domain.UserRepository
}

// NewGetUserHandler creates a new handler for user retrieval
func NewGetUserHandler(repo domain.UserRepository) GetUserHandler {
	return &getUserHandler{
		userRepository: repo,
	}
}

// Handle processes the get user query
func (h *getUserHandler) Handle(ctx context.Context, query GetUserQuery) (*domain.User, error) {
	return h.userRepository.GetByID(ctx, query.ID)
}
