package commands

import (
	"context"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/pkg/errors"
)

// UpdatePatientCommand represents the command to update a patient
type UpdatePatientCommand struct {
	ID          string  `json:"-"`
	FirstName   string  `json:"firstName" binding:"required"`
	LastName    string  `json:"lastName" binding:"required"`
	MiddleName  *string `json:"middleName,omitempty"`
	DateOfBirth string  `json:"dateOfBirth" binding:"required"`
	Gender      string  `json:"gender" binding:"required,oneof=male female other unknown"`
	Email       string  `json:"email" binding:"required,email"`
	PhoneNumber string  `json:"phoneNumber" binding:"required"`
	Address     *struct {
		Street     string `json:"street" binding:"required"`
		City       string `json:"city" binding:"required"`
		State      string `json:"state" binding:"required"`
		PostalCode string `json:"postalCode" binding:"required"`
		Country    string `json:"country" binding:"required"`
	} `json:"address,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
}

// UpdatePatientHandler handles the update patient command
type UpdatePatientHandler interface {
	Handle(ctx context.Context, cmd UpdatePatientCommand) (*domain.Patient, error)
}

type updatePatientHandler struct {
	repo domain.PatientRepository
}

// NewUpdatePatientHandler creates a new update patient handler
func NewUpdatePatientHandler(repo domain.PatientRepository) UpdatePatientHandler {
	return &updatePatientHandler{repo: repo}
}

// Handle processes the update patient command
func (h *updatePatientHandler) Handle(ctx context.Context, cmd UpdatePatientCommand) (*domain.Patient, error) {
	// Get the existing patient
	patient, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.NewAPIError(errors.ErrNotFound, "Patient not found")
	}

	// Parse the date of birth
	dob, err := time.Parse("2006-01-02", cmd.DateOfBirth)
	if err != nil {
		return nil, errors.NewAPIError(errors.ErrValidation, "Invalid date of birth format")
	}

	// Update the patient fields
	patient.FirstName = cmd.FirstName
	patient.LastName = cmd.LastName
	if cmd.MiddleName != nil {
		patient.MiddleName = *cmd.MiddleName
	}
	patient.DateOfBirth = domain.Date(dob)
	patient.Gender = domain.Gender(cmd.Gender)
	patient.Email = cmd.Email
	patient.PhoneNumber = cmd.PhoneNumber
	if cmd.Height != nil {
		patient.Height = *cmd.Height
	}
	if cmd.Weight != nil {
		patient.Weight = *cmd.Weight
	}

	if cmd.Address != nil {
		patient.Address = domain.Address{
			Street:     cmd.Address.Street,
			City:       cmd.Address.City,
			State:      cmd.Address.State,
			PostalCode: cmd.Address.PostalCode,
			Country:    cmd.Address.Country,
		}
	}

	// Save the updated patient
	if err := h.repo.Update(ctx, patient); err != nil {
		return nil, err
	}

	return patient, nil
}
