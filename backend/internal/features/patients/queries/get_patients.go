package queries

import (
	"context"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
)

// GetPatientsQuery represents the query to retrieve patients
type GetPatientsQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=20"`
	Search   string `form:"search"`
}

// PaginatedPatients represents a paginated list of patients
type PaginatedPatients struct {
	Patients    []*domain.Patient `json:"patients"`
	TotalCount  int64             `json:"totalCount"`
	CurrentPage int               `json:"currentPage"`
	PageSize    int               `json:"pageSize"`
	TotalPages  int               `json:"totalPages"`
}

// GetPatientsHandler handles the retrieval of patients
type GetPatientsHandler interface {
	Handle(ctx context.Context, query GetPatientsQuery) (*PaginatedPatients, error)
}

// NewGetPatientsHandler creates a new handler for patient retrieval
type getPatientsHandler struct {
	patientRepository domain.GetPatientsRepository
}

func NewGetPatientsHandler(repo domain.GetPatientsRepository) GetPatientsHandler {
	return &getPatientsHandler{
		patientRepository: repo,
	}
}

// Handle processes the get patients query
func (h *getPatientsHandler) Handle(ctx context.Context, query GetPatientsQuery) (*PaginatedPatients, error) {
	// Ensure valid pagination parameters
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}

	// Get patients with pagination
	patients, totalCount, err := h.patientRepository.ListPaginated(ctx, query.Page, query.PageSize, query.Search)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(totalCount) / query.PageSize
	if int(totalCount)%query.PageSize > 0 {
		totalPages++
	}

	return &PaginatedPatients{
		Patients:    patients,
		TotalCount:  totalCount,
		CurrentPage: query.Page,
		PageSize:    query.PageSize,
		TotalPages:  totalPages,
	}, nil
}
