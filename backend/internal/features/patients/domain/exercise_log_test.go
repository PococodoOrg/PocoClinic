package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidExercisePlanStatus(t *testing.T) {
	if !ValidExercisePlanStatus(ExercisePlanStatusActive) ||
		!ValidExercisePlanStatus(ExercisePlanStatusCompleted) ||
		!ValidExercisePlanStatus(ExercisePlanStatusArchived) {
		t.Fatal("expected known statuses to be valid")
	}
	if ValidExercisePlanStatus("paused") {
		t.Fatal("expected paused to be invalid")
	}
}

func TestNewExercisePlanAndEntry(t *testing.T) {
	patientID := uuid.New()
	authorID := uuid.New()
	plan := NewExercisePlan(patientID, authorID, "Knee", "goals")
	if plan.ID == uuid.Nil || plan.Status != ExercisePlanStatusActive {
		t.Fatalf("unexpected plan: %+v", plan)
	}

	diff := 5
	entry := NewExerciseLogEntry(
		plan.ID,
		patientID,
		authorID,
		time.Now().UTC(),
		"Quad sets",
		3,
		10,
		nil,
		"band",
		&diff,
		"good form",
	)
	if entry.PlanID != plan.ID || entry.Sets != 3 || entry.Difficulty == nil || *entry.Difficulty != 5 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}
