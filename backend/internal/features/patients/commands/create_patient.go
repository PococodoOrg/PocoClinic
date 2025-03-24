package commands

import (
	"context"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
)

// CreatePatientCommand represents the command to create a new patient
type CreatePatientCommand struct {
	FirstName   string         `json:"firstName" binding:"required"`
	LastName    string         `json:"lastName" binding:"required"`
	MiddleName  string         `json:"middleName"`
	DateOfBirth domain.Date    `json:"dateOfBirth" binding:"required"`
	Gender      domain.Gender  `json:"gender" binding:"required"`
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

// NewCreatePatientHandler creates a new handler for patient creation
type createPatientHandler struct {
	patientRepository domain.CreatePatientRepository
}

func NewCreatePatientHandler(repo domain.CreatePatientRepository) CreatePatientHandler {
	return &createPatientHandler{
		patientRepository: repo,
	}
}

// Handle processes the create patient command
func (h *createPatientHandler) Handle(ctx context.Context, cmd CreatePatientCommand) (*domain.Patient, error) {
	patient := domain.NewPatient(cmd.FirstName, cmd.LastName, cmd.DateOfBirth.Time(), cmd.Gender)
	patient.MiddleName = cmd.MiddleName
	patient.Email = cmd.Email
	patient.PhoneNumber = cmd.PhoneNumber
	patient.Height = cmd.Height
	patient.Weight = cmd.Weight
	patient.Address = cmd.Address

	err := h.patientRepository.Create(ctx, patient)
	if err != nil {
		return nil, err
	}

	return patient, nil
}
