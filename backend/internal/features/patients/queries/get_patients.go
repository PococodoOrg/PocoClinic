package queries

import (
	"context"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
)

// GetPatientsQuery represents the query to retrieve patients
type GetPatientsQuery struct {
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"pageSize,default=20"`
	Search          string `form:"search"`
	Gender          string `form:"gender"`
	DobFrom         string `form:"dobFrom"`
	DobTo           string `form:"dobTo"`
	RegisteredSince string `form:"registeredSince"`
}

func (q GetPatientsQuery) Filter() domain.PatientListFilter {
	return domain.PatientListFilter{
		Search:          q.Search,
		Gender:          q.Gender,
		DateOfBirthFrom: q.DobFrom,
		DateOfBirthTo:   q.DobTo,
		RegisteredSince: q.RegisteredSince,
	}
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
	patients, totalCount, err := h.patientRepository.ListPaginated(ctx, query.Page, query.PageSize, query.Filter())
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

// GetPatientQuery represents the query to retrieve a single patient by ID
type GetPatientQuery struct {
	ID string `json:"id"`
}

// GetPatientHandler handles the retrieval of a single patient
type GetPatientHandler interface {
	Handle(ctx context.Context, query GetPatientQuery) (*domain.Patient, error)
}

// NewGetPatientHandler creates a new handler for retrieving a single patient
type getPatientHandler struct {
	patientRepository domain.GetPatientRepository
}

func NewGetPatientHandler(repo domain.GetPatientRepository) GetPatientHandler {
	return &getPatientHandler{
		patientRepository: repo,
	}
}

func (h *getPatientHandler) Handle(ctx context.Context, query GetPatientQuery) (*domain.Patient, error) {
	return h.patientRepository.GetByID(ctx, query.ID)
}
