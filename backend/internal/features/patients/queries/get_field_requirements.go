package queries

import (
	"context"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
)

type GetFieldRequirementsHandler interface {
	Handle(ctx context.Context) (domain.PatientFieldRequirements, error)
}

type getFieldRequirementsHandler struct {
	settings domain.SettingsRepository
}

func NewGetFieldRequirementsHandler(settings domain.SettingsRepository) GetFieldRequirementsHandler {
	return &getFieldRequirementsHandler{settings: settings}
}

func (h *getFieldRequirementsHandler) Handle(ctx context.Context) (domain.PatientFieldRequirements, error) {
	reqs, err := h.settings.GetPatientFieldRequirements(ctx)
	if err != nil {
		return domain.PatientFieldRequirements{}, err
	}
	return domain.NormalizePatientFieldRequirements(reqs), nil
}
