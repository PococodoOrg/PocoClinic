package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
)

// CreatePatientCommand represents the command to create a new patient
type CreatePatientCommand struct {
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	MiddleName  string         `json:"middleName"`
	DateOfBirth domain.Date    `json:"dateOfBirth"`
	Gender      domain.Gender  `json:"gender"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phoneNumber"`
	Height      float64        `json:"height,omitempty"`
	Weight      float64        `json:"weight,omitempty"`
	Address     domain.Address `json:"address"`
}

// CreatePatientHandler handles the creation of a new patient
type CreatePatientHandler interface {
	Handle(ctx context.Context, cmd CreatePatientCommand) (*domain.Patient, error)
}

type createPatientHandler struct {
	patientRepository domain.CreatePatientRepository
	settings          domain.SettingsRepository
}

func NewCreatePatientHandler(repo domain.CreatePatientRepository, settings domain.SettingsRepository) CreatePatientHandler {
	return &createPatientHandler{
		patientRepository: repo,
		settings:          settings,
	}
}

func (h *createPatientHandler) Handle(ctx context.Context, cmd CreatePatientCommand) (*domain.Patient, error) {
	reqs, err := h.settings.GetPatientFieldRequirements(ctx)
	if err != nil {
		return nil, fmt.Errorf("load patient field requirements: %w", err)
	}

	dob := ""
	if !cmd.DateOfBirth.Time().IsZero() {
		dob = cmd.DateOfBirth.Time().Format("2006-01-02")
	}
	var height *float64
	if cmd.Height > 0 {
		height = &cmd.Height
	}
	var weight *float64
	if cmd.Weight > 0 {
		weight = &cmd.Weight
	}

	if err := domain.ValidatePatientInput(domain.PatientInput{
		FirstName:   cmd.FirstName,
		LastName:    cmd.LastName,
		MiddleName:  cmd.MiddleName,
		DateOfBirth: dob,
		Gender:      string(cmd.Gender),
		Email:       cmd.Email,
		PhoneNumber: cmd.PhoneNumber,
		Address:     cmd.Address,
		Height:      height,
		Weight:      weight,
	}, reqs); err != nil {
		return nil, err
	}

	if dob == "" {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Date of birth is required")
	}
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid date of birth format")
	}
	if !domain.ValidGender(string(cmd.Gender)) {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid gender")
	}

	patient := domain.NewPatient(cmd.FirstName, cmd.LastName, parsedDOB, cmd.Gender)
	patient.MiddleName = cmd.MiddleName
	patient.Email = cmd.Email
	patient.PhoneNumber = cmd.PhoneNumber
	patient.Height = cmd.Height
	patient.Weight = cmd.Weight
	patient.Address = cmd.Address

	if err := h.patientRepository.Create(ctx, patient); err != nil {
		return nil, err
	}

	return patient, nil
}
