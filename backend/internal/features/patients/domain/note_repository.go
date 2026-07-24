package domain

import "context"

// NoteRepository persists patient chart notes.
type NoteRepository interface {
	Create(ctx context.Context, note *Note) error
	GetByID(ctx context.Context, patientID, noteID string) (*Note, error)
	ListByPatientID(ctx context.Context, patientID string) ([]*Note, error)
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, patientID, noteID string) error
}
