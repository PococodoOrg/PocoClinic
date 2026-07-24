package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaxExercisePlanNameLength        = 200
	MaxExercisePlanDescriptionLength = 2_000
	MaxExerciseNameLength            = 200
	MaxExerciseResistanceLength      = 100
	MaxExerciseEntryNotesLength      = 2_000

	ExercisePlanStatusActive    = "active"
	ExercisePlanStatusCompleted = "completed"
	ExercisePlanStatusArchived  = "archived"
)

// ExercisePlan is a physical-therapy plan attached to a patient chart.
type ExercisePlan struct {
	ID          uuid.UUID  `json:"id"`
	PatientID   uuid.UUID  `json:"patientId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedBy   uuid.UUID  `json:"createdBy"`
	CreatedByName string   `json:"createdByName"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// ExerciseLogEntry is one logged session (sets/reps) under a plan.
type ExerciseLogEntry struct {
	ID              uuid.UUID  `json:"id"`
	PlanID          uuid.UUID  `json:"planId"`
	PatientID       uuid.UUID  `json:"patientId"`
	PerformedAt     time.Time  `json:"performedAt"`
	ExerciseName    string     `json:"exerciseName"`
	Sets            int        `json:"sets"`
	Reps            int        `json:"reps"`
	DurationSeconds *int       `json:"durationSeconds,omitempty"`
	Resistance      string     `json:"resistance"`
	Difficulty      *int       `json:"difficulty,omitempty"`
	Notes           string     `json:"notes"`
	RecordedBy      uuid.UUID  `json:"recordedBy"`
	RecordedByName  string     `json:"recordedByName"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}

func NewExercisePlan(patientID, createdBy uuid.UUID, name, description string) *ExercisePlan {
	now := time.Now()
	return &ExercisePlan{
		ID:          uuid.New(),
		PatientID:   patientID,
		Name:        name,
		Description: description,
		Status:      ExercisePlanStatusActive,
		CreatedBy:   createdBy,
		CreatedAt:   now,
	}
}

func NewExerciseLogEntry(
	planID, patientID, recordedBy uuid.UUID,
	performedAt time.Time,
	exerciseName string,
	sets, reps int,
	durationSeconds *int,
	resistance string,
	difficulty *int,
	notes string,
) *ExerciseLogEntry {
	now := time.Now()
	return &ExerciseLogEntry{
		ID:              uuid.New(),
		PlanID:          planID,
		PatientID:       patientID,
		PerformedAt:     performedAt,
		ExerciseName:    exerciseName,
		Sets:            sets,
		Reps:            reps,
		DurationSeconds: durationSeconds,
		Resistance:      resistance,
		Difficulty:      difficulty,
		Notes:           notes,
		RecordedBy:      recordedBy,
		CreatedAt:       now,
	}
}

func ValidExercisePlanStatus(status string) bool {
	switch status {
	case ExercisePlanStatusActive, ExercisePlanStatusCompleted, ExercisePlanStatusArchived:
		return true
	default:
		return false
	}
}
