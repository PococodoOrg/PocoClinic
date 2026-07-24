package commands

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	patientinfra "github.com/dksch/pococlinic/internal/features/patients/infrastructure"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/google/uuid"
)

func newExerciseTestEnv(t *testing.T) (*domain.Patient, *ExerciseLogCommandHandler, uuid.UUID) {
	t.Helper()
	patientRepo := patientinfra.NewMemoryRepository()
	patient := domain.NewPatient("Pat", "Exercise", time.Date(1980, 1, 2, 0, 0, 0, 0, time.UTC), domain.GenderFemale)
	if err := patientRepo.Create(context.Background(), patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}
	repo := patientinfra.NewMemoryExerciseLogRepository(patientRepo)
	handler := NewExerciseLogCommandHandler(repo, patientRepo)
	return patient, handler, uuid.New()
}

func TestExerciseLogCommandHandler_CreatePlanAndEntry(t *testing.T) {
	patient, handler, authorID := newExerciseTestEnv(t)

	plan, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID:   patient.ID.String(),
		CreatedBy:   authorID,
		Name:        "Knee strengthening",
		Description: "Home PT week 1-4",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}
	if plan.Status != domain.ExercisePlanStatusActive {
		t.Fatalf("expected active status, got %s", plan.Status)
	}

	entry, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:    patient.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   authorID,
		ExerciseName: "Quad sets",
		Sets:         3,
		Reps:         10,
	})
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if entry.Sets != 3 || entry.Reps != 10 {
		t.Fatalf("unexpected sets/reps: %d/%d", entry.Sets, entry.Reps)
	}

	updated, err := handler.UpdateEntry(context.Background(), UpdateExerciseEntryCommand{
		PatientID:    patient.ID.String(),
		PlanID:       plan.ID.String(),
		EntryID:      entry.ID.String(),
		ExerciseName: "Quad sets",
		Sets:         3,
		Reps:         12,
	})
	if err != nil {
		t.Fatalf("update entry: %v", err)
	}
	if updated.Reps != 12 {
		t.Fatalf("expected reps 12, got %d", updated.Reps)
	}
}

func TestExerciseLogCommandHandler_UpdateAndDeletePlan(t *testing.T) {
	patient, handler, authorID := newExerciseTestEnv(t)

	plan, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID: patient.ID.String(),
		CreatedBy: authorID,
		Name:      "Shoulder mobility",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}

	updated, err := handler.UpdatePlan(context.Background(), UpdateExercisePlanCommand{
		PatientID:   patient.ID.String(),
		PlanID:      plan.ID.String(),
		Name:        "Shoulder mobility v2",
		Description: "Done",
		Status:      domain.ExercisePlanStatusCompleted,
	})
	if err != nil {
		t.Fatalf("update plan: %v", err)
	}
	if updated.Status != domain.ExercisePlanStatusCompleted || updated.Name != "Shoulder mobility v2" {
		t.Fatalf("unexpected plan update: %+v", updated)
	}

	if err := handler.DeletePlan(context.Background(), DeleteExercisePlanCommand{
		PatientID: patient.ID.String(),
		PlanID:    plan.ID.String(),
	}); err != nil {
		t.Fatalf("delete plan: %v", err)
	}

	_, err = handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:    patient.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   authorID,
		ExerciseName: "External rotation",
		Sets:         2,
		Reps:         15,
	})
	if err == nil {
		t.Fatal("expected error creating entry on deleted plan")
	}
}

func TestExerciseLogCommandHandler_EntryWithOptionalFields(t *testing.T) {
	patient, handler, authorID := newExerciseTestEnv(t)
	plan, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID: patient.ID.String(),
		CreatedBy: authorID,
		Name:      "Timed holds",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}

	duration := 45
	difficulty := 6
	performed := time.Date(2026, 3, 1, 15, 30, 0, 0, time.UTC)
	entry, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:       patient.ID.String(),
		PlanID:          plan.ID.String(),
		RecordedBy:      authorID,
		PerformedAt:     &performed,
		ExerciseName:    "Wall sit",
		Sets:            3,
		Reps:            1,
		DurationSeconds: &duration,
		Resistance:      "bodyweight",
		Difficulty:      &difficulty,
		Notes:           "Stopped at burn",
	})
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if entry.DurationSeconds == nil || *entry.DurationSeconds != 45 {
		t.Fatalf("expected duration 45, got %+v", entry.DurationSeconds)
	}
	if !entry.PerformedAt.Equal(performed) {
		t.Fatalf("expected performedAt %v, got %v", performed, entry.PerformedAt)
	}

	neg := -1
	_, err = handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:       patient.ID.String(),
		PlanID:          plan.ID.String(),
		RecordedBy:      authorID,
		ExerciseName:    "Wall sit",
		Sets:            1,
		Reps:            1,
		DurationSeconds: &neg,
	})
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestExerciseLogCommandHandler_DeleteEntry(t *testing.T) {
	patient, handler, authorID := newExerciseTestEnv(t)
	plan, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID: patient.ID.String(),
		CreatedBy: authorID,
		Name:      "Core",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}
	entry, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:    patient.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   authorID,
		ExerciseName: "Plank",
		Sets:         1,
		Reps:         1,
	})
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if err := handler.DeleteEntry(context.Background(), DeleteExerciseEntryCommand{
		PatientID: patient.ID.String(),
		PlanID:    plan.ID.String(),
		EntryID:   entry.ID.String(),
	}); err != nil {
		t.Fatalf("delete entry: %v", err)
	}
}

func TestExerciseLogCommandHandler_Validation(t *testing.T) {
	patient, handler, authorID := newExerciseTestEnv(t)
	plan, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID: patient.ID.String(),
		CreatedBy: authorID,
		Name:      "Valid",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}

	cases := []struct {
		name string
		run  func() error
	}{
		{
			name: "empty plan name",
			run: func() error {
				_, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
					PatientID: patient.ID.String(),
					CreatedBy: authorID,
					Name:      "   ",
				})
				return err
			},
		},
		{
			name: "invalid plan status",
			run: func() error {
				_, err := handler.UpdatePlan(context.Background(), UpdateExercisePlanCommand{
					PatientID: patient.ID.String(),
					PlanID:    plan.ID.String(),
					Name:      "Valid",
					Status:    "paused",
				})
				return err
			},
		},
		{
			name: "empty exercise name",
			run: func() error {
				_, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
					PatientID:    patient.ID.String(),
					PlanID:       plan.ID.String(),
					RecordedBy:   authorID,
					ExerciseName: " ",
					Sets:         1,
					Reps:         1,
				})
				return err
			},
		},
		{
			name: "negative sets",
			run: func() error {
				_, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
					PatientID:    patient.ID.String(),
					PlanID:       plan.ID.String(),
					RecordedBy:   authorID,
					ExerciseName: "Squat",
					Sets:         -1,
					Reps:         5,
				})
				return err
			},
		},
		{
			name: "difficulty out of range",
			run: func() error {
				diff := 11
				_, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
					PatientID:    patient.ID.String(),
					PlanID:       plan.ID.String(),
					RecordedBy:   authorID,
					ExerciseName: "Squat",
					Sets:         1,
					Reps:         5,
					Difficulty:   &diff,
				})
				return err
			},
		},
		{
			name: "missing patient",
			run: func() error {
				_, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
					PatientID: uuid.New().String(),
					CreatedBy: authorID,
					Name:      "Orphan",
				})
				return err
			},
		},
		{
			name: "plan name too long",
			run: func() error {
				_, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
					PatientID: patient.ID.String(),
					CreatedBy: authorID,
					Name:      strings.Repeat("a", domain.MaxExercisePlanNameLength+1),
				})
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if err == nil {
				t.Fatal("expected validation/not-found error")
			}
			if tc.name == "missing patient" {
				return
			}
			var apiErr *pkgerrors.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T: %v", err, err)
			}
		})
	}
}
