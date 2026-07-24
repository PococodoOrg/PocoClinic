package domain

import (
	"time"

	"github.com/google/uuid"
)

// FormType defines how instances of a template behave for each patient.
type FormType string

const (
	FormTypeSingleton FormType = "singleton"
	FormTypeLog       FormType = "log"
)

// FormGroup organizes templates in the patient chart UI.
type FormGroup struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// DefaultGroupID is the seeded general-purpose group.
var DefaultGroupID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

// NewFormGroup creates a group with timestamps.
func NewFormGroup(name string, sortOrder int) *FormGroup {
	now := time.Now()
	return &FormGroup{
		ID:        uuid.New(),
		Name:      name,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
