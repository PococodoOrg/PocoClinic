package infrastructure

import (
	"context"
	"sync"

	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/google/uuid"
)

type noteRecord struct {
	note       *domain.Note
	authorName string
}

// MemoryNoteRepository stores patient notes in memory for development.
type MemoryNoteRepository struct {
	mu        sync.RWMutex
	notes     map[uuid.UUID]noteRecord
	patients  domain.PatientRepository
	userNames map[uuid.UUID]string
}

func NewMemoryNoteRepository(patients domain.PatientRepository) *MemoryNoteRepository {
	return &MemoryNoteRepository{
		notes:     make(map[uuid.UUID]noteRecord),
		patients:  patients,
		userNames: make(map[uuid.UUID]string),
	}
}

func (r *MemoryNoteRepository) SetAuthorName(authorID uuid.UUID, name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userNames[authorID] = name
}

func (r *MemoryNoteRepository) Create(ctx context.Context, note *domain.Note) error {
	if err := r.PatientExists(ctx, note.PatientID.String()); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	authorName := r.userNames[note.AuthorID]
	if authorName == "" {
		authorName = "Staff"
	}

	r.notes[note.ID] = noteRecord{note: note, authorName: authorName}
	return nil
}

func (r *MemoryNoteRepository) ListByPatientID(ctx context.Context, patientID string) ([]*domain.Note, error) {
	if err := r.PatientExists(ctx, patientID); err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	notes := make([]*domain.Note, 0)
	for _, record := range r.notes {
		if record.note.PatientID != parsedID {
			continue
		}
		copy := *record.note
		copy.AuthorName = record.authorName
		notes = append(notes, &copy)
	}

	sortNotesDesc(notes)
	return notes, nil
}

func (r *MemoryNoteRepository) GetByID(ctx context.Context, patientID, noteID string) (*domain.Note, error) {
	if err := r.PatientExists(ctx, patientID); err != nil {
		return nil, err
	}

	parsedPatientID, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}
	parsedNoteID, err := uuid.Parse(noteID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid note ID")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.notes[parsedNoteID]
	if !ok || record.note.PatientID != parsedPatientID {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Note not found")
	}
	copy := *record.note
	copy.AuthorName = record.authorName
	return &copy, nil
}

func (r *MemoryNoteRepository) Update(ctx context.Context, note *domain.Note) error {
	if err := r.PatientExists(ctx, note.PatientID.String()); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.notes[note.ID]
	if !ok {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Note not found")
	}
	copy := *note
	record.note = &copy
	record.authorName = note.AuthorName
	r.notes[note.ID] = record
	return nil
}

func (r *MemoryNoteRepository) Delete(ctx context.Context, patientID, noteID string) error {
	if _, err := r.GetByID(ctx, patientID, noteID); err != nil {
		return err
	}

	parsedNoteID, err := uuid.Parse(noteID)
	if err != nil {
		return pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid note ID")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.notes, parsedNoteID)
	return nil
}

func (r *MemoryNoteRepository) PatientExists(ctx context.Context, patientID string) error {
	if _, err := r.patients.GetByID(ctx, patientID); err != nil {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Patient not found")
	}
	return nil
}

func sortNotesDesc(notes []*domain.Note) {
	for i := 0; i < len(notes); i++ {
		for j := i + 1; j < len(notes); j++ {
			if notes[j].CreatedAt.After(notes[i].CreatedAt) {
				notes[i], notes[j] = notes[j], notes[i]
			}
		}
	}
}

var _ domain.NoteRepository = (*MemoryNoteRepository)(nil)
