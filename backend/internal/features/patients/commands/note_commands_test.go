package commands

import (
	"context"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPatientRepo struct {
	patient *domain.Patient
}

func (s stubPatientRepo) GetByID(_ context.Context, id string) (*domain.Patient, error) {
	if s.patient != nil && s.patient.ID.String() == id {
		return s.patient, nil
	}
	return nil, assert.AnError
}

type stubNoteRepo struct {
	notes   map[uuid.UUID]*domain.Note
	created *domain.Note
}

func newStubNoteRepo(note *domain.Note) *stubNoteRepo {
	repo := &stubNoteRepo{notes: make(map[uuid.UUID]*domain.Note)}
	if note != nil {
		copy := *note
		repo.notes[note.ID] = &copy
	}
	return repo
}

func (s *stubNoteRepo) Create(_ context.Context, note *domain.Note) error {
	s.created = note
	copy := *note
	s.notes[note.ID] = &copy
	return nil
}

func (s *stubNoteRepo) GetByID(_ context.Context, patientID, noteID string) (*domain.Note, error) {
	parsedPatientID, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}
	parsedNoteID, err := uuid.Parse(noteID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid note ID")
	}
	note, ok := s.notes[parsedNoteID]
	if !ok || note.PatientID != parsedPatientID {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Note not found")
	}
	copy := *note
	return &copy, nil
}

func (s *stubNoteRepo) ListByPatientID(_ context.Context, _ string) ([]*domain.Note, error) {
	return nil, nil
}

func (s *stubNoteRepo) Update(_ context.Context, note *domain.Note) error {
	copy := *note
	s.notes[note.ID] = &copy
	return nil
}

func (s *stubNoteRepo) Delete(_ context.Context, patientID, noteID string) error {
	note, err := s.GetByID(context.Background(), patientID, noteID)
	if err != nil {
		return err
	}
	delete(s.notes, note.ID)
	return nil
}

func TestCreateNoteHandler_RejectsEmptyBody(t *testing.T) {
	patientID := uuid.New()
	handler := NewCreateNoteHandler(newStubNoteRepo(nil), stubPatientRepo{
		patient: &domain.Patient{ID: patientID},
	})

	_, err := handler.Handle(context.Background(), CreateNoteCommand{
		PatientID: patientID.String(),
		AuthorID:  uuid.New(),
		Body:      "   ",
	})
	require.Error(t, err)
}

func TestCreateNoteHandler_Success(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	noteRepo := newStubNoteRepo(nil)
	handler := NewCreateNoteHandler(noteRepo, stubPatientRepo{
		patient: &domain.Patient{ID: patientID},
	})

	note, err := handler.Handle(context.Background(), CreateNoteCommand{
		PatientID: patientID.String(),
		AuthorID:  authorID,
		Body:      "  Follow-up in two weeks.  ",
	})
	require.NoError(t, err)
	assert.Equal(t, patientID, note.PatientID)
	assert.Equal(t, authorID, note.AuthorID)
	assert.Equal(t, "Follow-up in two weeks.", note.Body)
	assert.WithinDuration(t, time.Now(), note.CreatedAt, 2*time.Second)
	assert.NotNil(t, noteRepo.created)
}

func TestUpdateNoteHandler_RequiresAuthor(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	noteID := uuid.New()
	noteRepo := newStubNoteRepo(&domain.Note{
		ID:        noteID,
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      "Original",
		CreatedAt: time.Now(),
	})
	handler := NewUpdateNoteHandler(noteRepo)

	_, err := handler.Handle(context.Background(), UpdateNoteCommand{
		PatientID: patientID.String(),
		NoteID:    noteID.String(),
		AuthorID:  uuid.New(),
		Body:      "Updated",
	})
	require.Error(t, err)
}

func TestUpdateNoteHandler_Success(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	noteID := uuid.New()
	noteRepo := newStubNoteRepo(&domain.Note{
		ID:        noteID,
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      "Original",
		CreatedAt: time.Now(),
	})
	handler := NewUpdateNoteHandler(noteRepo)

	note, err := handler.Handle(context.Background(), UpdateNoteCommand{
		PatientID: patientID.String(),
		NoteID:    noteID.String(),
		AuthorID:  authorID,
		Body:      " Updated text ",
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated text", note.Body)
	assert.NotNil(t, note.UpdatedAt)
}

func TestDeleteNoteHandler_RequiresAuthorOrAdmin(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	noteID := uuid.New()
	noteRepo := newStubNoteRepo(&domain.Note{
		ID:        noteID,
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      "Original",
		CreatedAt: time.Now(),
	})
	handler := NewDeleteNoteHandler(noteRepo)

	_, err := handler.Handle(context.Background(), DeleteNoteCommand{
		PatientID: patientID.String(),
		NoteID:    noteID.String(),
		ActorID:   uuid.New(),
		IsAdmin:   false,
	})
	require.Error(t, err)
}

func TestDeleteNoteHandler_AdminCanDelete(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	noteID := uuid.New()
	noteRepo := newStubNoteRepo(&domain.Note{
		ID:        noteID,
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      "Original",
		CreatedAt: time.Now(),
	})
	handler := NewDeleteNoteHandler(noteRepo)

	deleted, err := handler.Handle(context.Background(), DeleteNoteCommand{
		PatientID: patientID.String(),
		NoteID:    noteID.String(),
		ActorID:   uuid.New(),
		IsAdmin:   true,
	})
	require.NoError(t, err)
	assert.Equal(t, noteID, deleted.ID)
}

func TestDeleteNoteHandler_AuthorCanDelete(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	noteID := uuid.New()
	noteRepo := newStubNoteRepo(&domain.Note{
		ID:        noteID,
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      "Original",
		CreatedAt: time.Now(),
	})
	handler := NewDeleteNoteHandler(noteRepo)

	deleted, err := handler.Handle(context.Background(), DeleteNoteCommand{
		PatientID: patientID.String(),
		NoteID:    noteID.String(),
		ActorID:   authorID,
		IsAdmin:   false,
	})
	require.NoError(t, err)
	assert.Equal(t, noteID, deleted.ID)
	_, err = noteRepo.GetByID(context.Background(), patientID.String(), noteID.String())
	require.Error(t, err)
}
