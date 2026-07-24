package commands

import (
	"context"
	"strings"
	"time"

	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/google/uuid"
)

type UpdateNoteCommand struct {
	PatientID string
	NoteID    string
	AuthorID  uuid.UUID
	Body      string `json:"body" binding:"required"`
}

type UpdateNoteHandler struct {
	noteRepo domain.NoteRepository
}

func NewUpdateNoteHandler(noteRepo domain.NoteRepository) UpdateNoteHandler {
	return UpdateNoteHandler{noteRepo: noteRepo}
}

func (h UpdateNoteHandler) Handle(ctx context.Context, cmd UpdateNoteCommand) (*domain.Note, error) {
	body, err := validateNoteBody(cmd.Body)
	if err != nil {
		return nil, err
	}

	note, err := h.noteRepo.GetByID(ctx, cmd.PatientID, cmd.NoteID)
	if err != nil {
		return nil, err
	}
	if note.AuthorID != cmd.AuthorID {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrForbidden, "Only the author can edit this note")
	}

	now := time.Now()
	note.Body = body
	note.UpdatedAt = &now
	if err := h.noteRepo.Update(ctx, note); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to update note")
	}
	return note, nil
}

type DeleteNoteCommand struct {
	PatientID string
	NoteID    string
	ActorID   uuid.UUID
	IsAdmin   bool
}

type DeleteNoteHandler struct {
	noteRepo domain.NoteRepository
}

func NewDeleteNoteHandler(noteRepo domain.NoteRepository) DeleteNoteHandler {
	return DeleteNoteHandler{noteRepo: noteRepo}
}

func (h DeleteNoteHandler) Handle(ctx context.Context, cmd DeleteNoteCommand) (*domain.Note, error) {
	note, err := h.noteRepo.GetByID(ctx, cmd.PatientID, cmd.NoteID)
	if err != nil {
		return nil, err
	}
	if !cmd.IsAdmin && note.AuthorID != cmd.ActorID {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrForbidden, "Only the author or an administrator can delete this note")
	}
	if err := h.noteRepo.Delete(ctx, cmd.PatientID, cmd.NoteID); err != nil {
		return nil, err
	}
	return note, nil
}

func validateNoteBody(body string) (string, error) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return "", pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Note body is required")
	}
	if len(trimmed) > domain.MaxNoteBodyLength {
		return "", pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Note is too long")
	}
	return trimmed, nil
}
