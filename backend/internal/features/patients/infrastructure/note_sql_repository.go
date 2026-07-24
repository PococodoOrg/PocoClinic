package infrastructure

import (
	"context"
	"errors"
	"time"

	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

// SQLNoteRepository persists patient notes in SQLite.
type SQLNoteRepository struct {
	pool *database.DB
}

func NewSQLNoteRepository(pool *database.DB) *SQLNoteRepository {
	return &SQLNoteRepository{pool: pool}
}

func (r *SQLNoteRepository) Create(ctx context.Context, note *domain.Note) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO patient_notes (id, patient_id, author_id, body, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, note.ID, note.PatientID, note.AuthorID, note.Body, note.CreatedAt)
	if err != nil {
		return err
	}

	_ = r.pool.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, note.AuthorID).Scan(&note.AuthorName)
	return nil
}

func (r *SQLNoteRepository) ListByPatientID(ctx context.Context, patientID string) ([]*domain.Note, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT n.id, n.patient_id, n.author_id, u.name, n.body, n.created_at, n.updated_at
		FROM patient_notes n
		JOIN users u ON u.id = n.author_id
		WHERE n.patient_id = $1
		ORDER BY n.created_at DESC
	`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]*domain.Note, 0)
	for rows.Next() {
		note := &domain.Note{}
		var updatedAt *time.Time
		if err := rows.Scan(
			&note.ID,
			&note.PatientID,
			&note.AuthorID,
			&note.AuthorName,
			&note.Body,
			database.Time(&note.CreatedAt),
			database.NullTime(&updatedAt),
		); err != nil {
			return nil, err
		}
		note.UpdatedAt = updatedAt
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func (r *SQLNoteRepository) GetByID(ctx context.Context, patientID, noteID string) (*domain.Note, error) {
	note := &domain.Note{}
	var updatedAt *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT n.id, n.patient_id, n.author_id, u.name, n.body, n.created_at, n.updated_at
		FROM patient_notes n
		JOIN users u ON u.id = n.author_id
		WHERE n.patient_id = $1 AND n.id = $2
	`, patientID, noteID).Scan(
		&note.ID,
		&note.PatientID,
		&note.AuthorID,
		&note.AuthorName,
		&note.Body,
		database.Time(&note.CreatedAt),
		database.NullTime(&updatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Note not found")
		}
		return nil, err
	}
	note.UpdatedAt = updatedAt
	return note, nil
}

func (r *SQLNoteRepository) Update(ctx context.Context, note *domain.Note) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE patient_notes
		SET body = $1, updated_at = $2
		WHERE id = $3 AND patient_id = $4
	`, note.Body, note.UpdatedAt, note.ID, note.PatientID)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, note.AuthorID).Scan(&note.AuthorName)
}

func (r *SQLNoteRepository) Delete(ctx context.Context, patientID, noteID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM patient_notes WHERE patient_id = $1 AND id = $2
	`, patientID, noteID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Note not found")
	}
	return nil
}

// PatientExists checks that a patient record exists before writing notes.
func (r *SQLNoteRepository) PatientExists(ctx context.Context, patientID string) error {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM patients WHERE id = $1)`, patientID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Patient not found")
	}
	return nil
}

var _ domain.NoteRepository = (*SQLNoteRepository)(nil)
