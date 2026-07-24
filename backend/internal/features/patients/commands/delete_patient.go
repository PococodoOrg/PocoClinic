package commands

import (
	"context"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
)

// DeletePatientHandler handles patient deletion
type DeletePatientHandler interface {
	Handle(ctx context.Context, id string) error
}

type deletePatientHandler struct {
	repo domain.PatientRepository
}

// NewDeletePatientHandler creates a new delete patient handler
func NewDeletePatientHandler(repo domain.PatientRepository) DeletePatientHandler {
	return &deletePatientHandler{repo: repo}
}

// Handle processes the delete patient command
func (h *deletePatientHandler) Handle(ctx context.Context, id string) error {
	patient, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if patient == nil {
		return errors.NewAPIError(errors.ErrNotFound, "Patient not found")
	}

	return h.repo.Delete(ctx, id)
}
