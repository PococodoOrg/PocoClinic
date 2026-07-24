package commands

import (
	"context"

	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/google/uuid"
)

type CreateNoteCommand struct {
	PatientID string `json:"-"`
	AuthorID  uuid.UUID
	Body      string `json:"body" binding:"required"`
}

type CreateNoteHandler struct {
	noteRepo    domain.NoteRepository
	patientRepo domain.GetPatientRepository
}

func NewCreateNoteHandler(noteRepo domain.NoteRepository, patientRepo domain.GetPatientRepository) CreateNoteHandler {
	return CreateNoteHandler{noteRepo: noteRepo, patientRepo: patientRepo}
}

func (h CreateNoteHandler) Handle(ctx context.Context, cmd CreateNoteCommand) (*domain.Note, error) {
	body, err := validateNoteBody(cmd.Body)
	if err != nil {
		return nil, err
	}

	if _, err := h.patientRepo.GetByID(ctx, cmd.PatientID); err != nil {
		return nil, err
	}

	patientID, err := uuid.Parse(cmd.PatientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}

	note := domain.NewNote(patientID, cmd.AuthorID, body)
	if err := h.noteRepo.Create(ctx, note); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to save note")
	}

	return note, nil
}
