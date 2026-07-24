package domain

import (
	"fmt"
	"strings"
	"time"

	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
)

const patientFieldRequirementsKey = "patient_field_requirements"

// PatientFieldRequirements defines which patient chart fields must be filled in.
type PatientFieldRequirements struct {
	FirstName         bool `json:"firstName"`
	LastName          bool `json:"lastName"`
	MiddleName        bool `json:"middleName"`
	DateOfBirth       bool `json:"dateOfBirth"`
	Gender            bool `json:"gender"`
	Email             bool `json:"email"`
	PhoneNumber       bool `json:"phoneNumber"`
	AddressStreet     bool `json:"addressStreet"`
	AddressCity       bool `json:"addressCity"`
	AddressState      bool `json:"addressState"`
	AddressPostalCode bool `json:"addressPostalCode"`
	Height            bool `json:"height"`
	Weight            bool `json:"weight"`
}

// DefaultPatientFieldRequirements is used when no clinic setting exists.
func DefaultPatientFieldRequirements() PatientFieldRequirements {
	return PatientFieldRequirements{
		FirstName:   true,
		LastName:    true,
		DateOfBirth: true,
		Gender:      true,
	}
}

// PatientFieldRequirementsSettingKey returns the clinic_settings row key.
func PatientFieldRequirementsSettingKey() string {
	return patientFieldRequirementsKey
}

// NormalizePatientFieldRequirements merges stored values with defaults and enforces core identity fields.
func NormalizePatientFieldRequirements(req PatientFieldRequirements) PatientFieldRequirements {
	defaults := DefaultPatientFieldRequirements()
	if !req.FirstName && !req.LastName && !req.DateOfBirth && !req.Gender &&
		!req.MiddleName && !req.Email && !req.PhoneNumber &&
		!req.AddressStreet && !req.AddressCity && !req.AddressState && !req.AddressPostalCode &&
		!req.Height && !req.Weight {
		return defaults
	}
	req.FirstName = req.FirstName || defaults.FirstName
	req.LastName = req.LastName || defaults.LastName
	req.DateOfBirth = req.DateOfBirth || defaults.DateOfBirth
	return req
}

// PatientInput captures create/update payload fields for validation.
type PatientInput struct {
	FirstName   string
	LastName    string
	MiddleName  string
	DateOfBirth string
	Gender      string
	Email       string
	PhoneNumber string
	Address     Address
	Height      *float64
	Weight      *float64
}

func requireField(required bool, value, label string) error {
	if !required {
		return nil
	}
	if strings.TrimSpace(value) == "" {
		return pkgerrors.NewAPIError(pkgerrors.ErrValidation, fmt.Sprintf("%s is required", label))
	}
	return nil
}

func requireFloat(required bool, value *float64, label string) error {
	if !required {
		return nil
	}
	if value == nil || *value <= 0 {
		return pkgerrors.NewAPIError(pkgerrors.ErrValidation, fmt.Sprintf("%s is required", label))
	}
	return nil
}

// ValidatePatientInput enforces configured required fields.
func ValidatePatientInput(input PatientInput, reqs PatientFieldRequirements) error {
	reqs = NormalizePatientFieldRequirements(reqs)

	checks := []struct {
		required bool
		value    string
		label    string
	}{
		{reqs.FirstName, input.FirstName, "First name"},
		{reqs.LastName, input.LastName, "Last name"},
		{reqs.MiddleName, input.MiddleName, "Middle name"},
		{reqs.Email, input.Email, "Email"},
		{reqs.PhoneNumber, input.PhoneNumber, "Phone number"},
		{reqs.AddressStreet, input.Address.Street, "Street address"},
		{reqs.AddressCity, input.Address.City, "City"},
		{reqs.AddressState, input.Address.State, "State"},
		{reqs.AddressPostalCode, input.Address.PostalCode, "ZIP code"},
	}
	for _, check := range checks {
		if err := requireField(check.required, check.value, check.label); err != nil {
			return err
		}
	}

	if reqs.DateOfBirth {
		if strings.TrimSpace(input.DateOfBirth) == "" {
			return pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Date of birth is required")
		}
		if _, err := time.Parse("2006-01-02", input.DateOfBirth); err != nil {
			return pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid date of birth format")
		}
	}

	if reqs.Gender {
		if strings.TrimSpace(input.Gender) == "" {
			return pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Gender is required")
		}
		if !ValidGender(input.Gender) {
			return pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid gender")
		}
	}

	if err := requireFloat(reqs.Height, input.Height, "Height"); err != nil {
		return err
	}
	if err := requireFloat(reqs.Weight, input.Weight, "Weight"); err != nil {
		return err
	}

	return nil
}
