package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaxNoteBodyLength = 10_000
)

// Note is a clinical chart note attached to a patient record.
type Note struct {
	ID         uuid.UUID `json:"id"`
	PatientID  uuid.UUID `json:"patientId"`
	AuthorID   uuid.UUID `json:"authorId"`
	AuthorName string    `json:"authorName"`
	Body       string     `json:"body"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

// NewNote creates a chart note with generated ID and timestamp.
func NewNote(patientID, authorID uuid.UUID, body string) *Note {
	now := time.Now()
	return &Note{
		ID:        uuid.New(),
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      body,
		CreatedAt: now,
	}
}
