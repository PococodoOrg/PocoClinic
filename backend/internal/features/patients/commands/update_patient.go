package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
)

// UpdatePatientCommand represents the command to update a patient
type UpdatePatientCommand struct {
	ID          string  `json:"-"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	MiddleName  *string `json:"middleName,omitempty"`
	DateOfBirth string  `json:"dateOfBirth"`
	Gender      string  `json:"gender"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Address     *struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		State      string `json:"state"`
		PostalCode string `json:"postalCode"`
		Country    string `json:"country"`
	} `json:"address,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
}

// UpdatePatientHandler handles the update patient command
type UpdatePatientHandler interface {
	Handle(ctx context.Context, cmd UpdatePatientCommand) (*domain.Patient, error)
}

type updatePatientHandler struct {
	repo     domain.PatientRepository
	settings domain.SettingsRepository
}

func NewUpdatePatientHandler(repo domain.PatientRepository, settings domain.SettingsRepository) UpdatePatientHandler {
	return &updatePatientHandler{repo: repo, settings: settings}
}

func (h *updatePatientHandler) Handle(ctx context.Context, cmd UpdatePatientCommand) (*domain.Patient, error) {
	patient, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.NewAPIError(errors.ErrNotFound, "Patient not found")
	}

	reqs, err := h.settings.GetPatientFieldRequirements(ctx)
	if err != nil {
		return nil, fmt.Errorf("load patient field requirements: %w", err)
	}

	middleName := patient.MiddleName
	if cmd.MiddleName != nil {
		middleName = *cmd.MiddleName
	}
	address := patient.Address
	if cmd.Address != nil {
		address = domain.Address{
			Street:     cmd.Address.Street,
			City:       cmd.Address.City,
			State:      cmd.Address.State,
			PostalCode: cmd.Address.PostalCode,
			Country:    cmd.Address.Country,
		}
	}

	if err := domain.ValidatePatientInput(domain.PatientInput{
		FirstName:   cmd.FirstName,
		LastName:    cmd.LastName,
		MiddleName:  middleName,
		DateOfBirth: cmd.DateOfBirth,
		Gender:      cmd.Gender,
		Email:       cmd.Email,
		PhoneNumber: cmd.PhoneNumber,
		Address:     address,
		Height:      cmd.Height,
		Weight:      cmd.Weight,
	}, reqs); err != nil {
		return nil, err
	}

	dob, err := time.Parse("2006-01-02", cmd.DateOfBirth)
	if err != nil {
		return nil, errors.NewAPIError(errors.ErrValidation, "Invalid date of birth format")
	}

	patient.FirstName = cmd.FirstName
	patient.LastName = cmd.LastName
	patient.MiddleName = middleName
	patient.DateOfBirth = domain.Date(dob)
	patient.Gender = domain.Gender(cmd.Gender)
	patient.Email = cmd.Email
	patient.PhoneNumber = cmd.PhoneNumber
	patient.Address = address
	if cmd.Height != nil {
		patient.Height = *cmd.Height
	}
	if cmd.Weight != nil {
		patient.Weight = *cmd.Weight
	}

	if err := h.repo.Update(ctx, patient); err != nil {
		return nil, err
	}

	return patient, nil
}
