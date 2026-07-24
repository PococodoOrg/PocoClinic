package domain

import (
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPatientFieldRequirements(t *testing.T) {
	reqs := DefaultPatientFieldRequirements()
	assert.True(t, reqs.FirstName)
	assert.True(t, reqs.LastName)
	assert.True(t, reqs.DateOfBirth)
	assert.True(t, reqs.Gender)
	assert.False(t, reqs.Email)
	assert.False(t, reqs.PhoneNumber)
}

func TestNormalizePatientFieldRequirements_emptyUsesDefaults(t *testing.T) {
	reqs := NormalizePatientFieldRequirements(PatientFieldRequirements{})
	assert.Equal(t, DefaultPatientFieldRequirements(), reqs)
}

func TestNormalizePatientFieldRequirements_enforcesCoreFields(t *testing.T) {
	reqs := NormalizePatientFieldRequirements(PatientFieldRequirements{
		FirstName:   false,
		LastName:    false,
		DateOfBirth: false,
		Email:       true,
	})
	assert.True(t, reqs.FirstName)
	assert.True(t, reqs.LastName)
	assert.True(t, reqs.DateOfBirth)
	assert.True(t, reqs.Email)
}

func TestValidatePatientInput_requiredEmail(t *testing.T) {
	err := ValidatePatientInput(PatientInput{
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		Gender:      "female",
	}, PatientFieldRequirements{
		FirstName:   true,
		LastName:    true,
		DateOfBirth: true,
		Gender:      true,
		Email:       true,
	})
	require.Error(t, err)
	apiErr, ok := err.(*errors.APIError)
	require.True(t, ok)
	assert.Equal(t, errors.ErrValidation, apiErr.Code)
	assert.Contains(t, apiErr.Message, "Email")
}

func TestValidatePatientInput_optionalContactFields(t *testing.T) {
	err := ValidatePatientInput(PatientInput{
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		Gender:      "female",
	}, DefaultPatientFieldRequirements())
	assert.NoError(t, err)
}
