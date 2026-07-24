package commands

import (
	"context"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
)

type UpdateFieldRequirementsCommand struct {
	Requirements domain.PatientFieldRequirements `json:"requirements" binding:"required"`
}

type UpdateFieldRequirementsHandler interface {
	Handle(ctx context.Context, cmd UpdateFieldRequirementsCommand) (domain.PatientFieldRequirements, error)
}

type updateFieldRequirementsHandler struct {
	settings domain.SettingsRepository
}

func NewUpdateFieldRequirementsHandler(settings domain.SettingsRepository) UpdateFieldRequirementsHandler {
	return &updateFieldRequirementsHandler{settings: settings}
}

func (h *updateFieldRequirementsHandler) Handle(ctx context.Context, cmd UpdateFieldRequirementsCommand) (domain.PatientFieldRequirements, error) {
	if !cmd.Requirements.FirstName || !cmd.Requirements.LastName || !cmd.Requirements.DateOfBirth {
		return domain.PatientFieldRequirements{}, pkgerrors.NewAPIError(
			pkgerrors.ErrValidation,
			"First name, last name, and date of birth must remain required",
		)
	}

	normalized := domain.NormalizePatientFieldRequirements(cmd.Requirements)
	if err := h.settings.SavePatientFieldRequirements(ctx, normalized); err != nil {
		return domain.PatientFieldRequirements{}, err
	}
	return normalized, nil
}
