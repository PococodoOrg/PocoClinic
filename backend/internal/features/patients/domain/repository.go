package domain

import (
	"context"
)

// PatientRepository defines the interface for patient persistence
type PatientRepository interface {
	Create(ctx context.Context, patient *Patient) error
	Update(ctx context.Context, patient *Patient) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Patient, error)
	List(ctx context.Context) ([]*Patient, error)
	ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*Patient, int64, error)
}

// CreatePatientRepository defines the minimal interface for patient creation
type CreatePatientRepository interface {
	Create(ctx context.Context, patient *Patient) error
}

// GetPatientsRepository defines the minimal interface for patient retrieval
type GetPatientsRepository interface {
	ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*Patient, int64, error)
}
