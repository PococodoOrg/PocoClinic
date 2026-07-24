package commands

import (
	"context"
	"testing"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/features/patients/infrastructure"
	"github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateFieldRequirementsHandler_rejectsDisablingCoreFields(t *testing.T) {
	handler := NewUpdateFieldRequirementsHandler(infrastructure.NewMemorySettingsRepository())

	_, err := handler.Handle(context.Background(), UpdateFieldRequirementsCommand{
		Requirements: domain.PatientFieldRequirements{
			FirstName:   false,
			LastName:    false,
			DateOfBirth: false,
		},
	})
	require.Error(t, err)
	apiErr, ok := err.(*errors.APIError)
	require.True(t, ok)
	assert.Equal(t, errors.ErrValidation, apiErr.Code)
}

func TestUpdateFieldRequirementsHandler_persistsSettings(t *testing.T) {
	repo := infrastructure.NewMemorySettingsRepository()
	handler := NewUpdateFieldRequirementsHandler(repo)

	reqs, err := handler.Handle(context.Background(), UpdateFieldRequirementsCommand{
		Requirements: domain.PatientFieldRequirements{
			FirstName:   true,
			LastName:    true,
			DateOfBirth: true,
			Gender:      true,
			Email:       true,
			PhoneNumber: true,
		},
	})
	require.NoError(t, err)
	assert.True(t, reqs.Email)
	assert.True(t, reqs.PhoneNumber)

	stored, err := repo.GetPatientFieldRequirements(context.Background())
	require.NoError(t, err)
	assert.True(t, stored.Email)
	assert.True(t, stored.PhoneNumber)
}
