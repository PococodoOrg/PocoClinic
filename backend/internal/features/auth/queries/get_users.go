package queries

import (
	"context"

	"github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/pagination"
)

// GetUsersQuery represents the query to retrieve users
type GetUsersQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=20"`
	Search   string `form:"search"`
}

// PaginatedUsers represents a paginated list of users
type PaginatedUsers struct {
	Users       []*domain.User `json:"users"`
	TotalCount  int64          `json:"totalCount"`
	CurrentPage int            `json:"currentPage"`
	PageSize    int            `json:"pageSize"`
	TotalPages  int            `json:"totalPages"`
}

// GetUsersHandler handles the retrieval of users
type GetUsersHandler interface {
	Handle(ctx context.Context, query GetUsersQuery) (*PaginatedUsers, error)
}

type getUsersHandler struct {
	userRepository domain.UserRepository
}

// NewGetUsersHandler creates a new handler for user retrieval
func NewGetUsersHandler(repo domain.UserRepository) GetUsersHandler {
	return &getUsersHandler{
		userRepository: repo,
	}
}

// Handle processes the get users query
func (h *getUsersHandler) Handle(ctx context.Context, query GetUsersQuery) (*PaginatedUsers, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	query.PageSize = pagination.Clamp(query.PageSize, 20)

	users, totalCount, err := h.userRepository.ListPaginated(ctx, query.Page, query.PageSize, query.Search)
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount) / query.PageSize
	if int(totalCount)%query.PageSize > 0 {
		totalPages++
	}

	return &PaginatedUsers{
		Users:       users,
		TotalCount:  totalCount,
		CurrentPage: query.Page,
		PageSize:    query.PageSize,
		TotalPages:  totalPages,
	}, nil
}
