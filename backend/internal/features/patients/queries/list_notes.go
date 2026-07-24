package queries

import (
	"context"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
)

type ListNotesQuery struct {
	PatientID string
}

type ListNotesHandler struct {
	noteRepo domain.NoteRepository
}

func NewListNotesHandler(noteRepo domain.NoteRepository) ListNotesHandler {
	return ListNotesHandler{noteRepo: noteRepo}
}

func (h ListNotesHandler) Handle(ctx context.Context, query ListNotesQuery) ([]*domain.Note, error) {
	return h.noteRepo.ListByPatientID(ctx, query.PatientID)
}
