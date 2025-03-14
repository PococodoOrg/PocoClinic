package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type patientTestSuite struct {
	defaultPatient *Patient
	defaultAddress Address
}

func setupPatientTest() patientTestSuite {
	address := Address{
		Street:     "123 Medical Drive",
		City:       "Healthcare City",
		State:      "HC",
		PostalCode: "12345",
		Country:    "Medical Land",
	}

	patient := NewPatient(
		"John",
		"Doe",
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		GenderMale,
	)
	patient.MiddleName = "Robert"
	patient.Email = "john.doe@example.com"
	patient.PhoneNumber = "555-0123"
	patient.Address = address

	return patientTestSuite{
		defaultPatient: patient,
		defaultAddress: address,
	}
}

func TestNewPatient(t *testing.T) {
	suite := setupPatientTest()
	patient := suite.defaultPatient

	assert.NotNil(t, patient.ID)
	assert.Equal(t, "John", patient.FirstName)
	assert.Equal(t, "Doe", patient.LastName)
	assert.Equal(t, "Robert", patient.MiddleName)
	assert.Equal(t, GenderMale, patient.Gender)
	assert.False(t, patient.CreatedAt.IsZero())
	assert.False(t, patient.UpdatedAt.IsZero())
}

func TestPatientFullName(t *testing.T) {
	testCases := []struct {
		name           string
		firstName      string
		middleName     string
		lastName       string
		expectedResult string
	}{
		{
			name:           "full_name_with_middle",
			firstName:      "John",
			middleName:     "Robert",
			lastName:       "Doe",
			expectedResult: "John Robert Doe",
		},
		{
			name:           "full_name_without_middle",
			firstName:      "John",
			middleName:     "",
			lastName:       "Doe",
			expectedResult: "John Doe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patient := NewPatient(tc.firstName, tc.lastName, time.Now(), GenderMale)
			patient.MiddleName = tc.middleName

			assert.Equal(t, tc.expectedResult, patient.FullName())
		})
	}
}

func TestPatientAge(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name        string
		dateOfBirth time.Time
		expectedAge int
	}{
		{
			name:        "age_full_years",
			dateOfBirth: now.AddDate(-30, 0, -1),
			expectedAge: 30,
		},
		{
			name:        "age_not_birthday_yet",
			dateOfBirth: now.AddDate(-30, 0, 1),
			expectedAge: 29,
		},
		{
			name:        "age_birthday_today",
			dateOfBirth: now.AddDate(-25, 0, 0),
			expectedAge: 25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patient := NewPatient("John", "Doe", tc.dateOfBirth, GenderMale)
			assert.Equal(t, tc.expectedAge, patient.Age())
		})
	}
}

func TestPatientUpdate(t *testing.T) {
	suite := setupPatientTest()
	patient := suite.defaultPatient

	originalUpdate := patient.UpdatedAt
	time.Sleep(time.Millisecond)

	patient.Update()

	assert.True(t, patient.UpdatedAt.After(originalUpdate))
}

func TestPatientGenderValidation(t *testing.T) {
	testCases := []struct {
		name   string
		gender Gender
		valid  bool
	}{
		{
			name:   "valid_male",
			gender: GenderMale,
			valid:  true,
		},
		{
			name:   "valid_female",
			gender: GenderFemale,
			valid:  true,
		},
		{
			name:   "valid_other",
			gender: GenderOther,
			valid:  true,
		},
		{
			name:   "valid_unknown",
			gender: GenderUnknown,
			valid:  true,
		},
		{
			name:   "invalid_gender",
			gender: Gender("invalid"),
			valid:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patient := NewPatient("John", "Doe", time.Now(), tc.gender)
			assert.Equal(t, tc.gender, patient.Gender)
		})
	}
}
