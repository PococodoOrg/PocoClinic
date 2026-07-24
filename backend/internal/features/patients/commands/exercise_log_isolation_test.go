package commands

import (
	"context"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	patientinfra "github.com/dksch/pococlinic/internal/features/patients/infrastructure"
	"github.com/dksch/pococlinic/internal/features/patients/queries"
	"github.com/google/uuid"
)

func TestExerciseLog_CannotAccessAnotherPatientsPlan(t *testing.T) {
	patientRepo := patientinfra.NewMemoryRepository()
	a := domain.NewPatient("A", "Patient", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), domain.GenderFemale)
	b := domain.NewPatient("B", "Patient", time.Date(1981, 1, 1, 0, 0, 0, 0, time.UTC), domain.GenderMale)
	if err := patientRepo.Create(context.Background(), a); err != nil {
		t.Fatalf("create a: %v", err)
	}
	if err := patientRepo.Create(context.Background(), b); err != nil {
		t.Fatalf("create b: %v", err)
	}

	repo := patientinfra.NewMemoryExerciseLogRepository(patientRepo)
	cmds := NewExerciseLogCommandHandler(repo, patientRepo)
	qry := queries.NewExerciseLogQueryHandler(repo)
	author := uuid.New()

	plan, err := cmds.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID: a.ID.String(),
		CreatedBy: author,
		Name:      "Only for A",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}

	_, err = cmds.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:    b.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   author,
		ExerciseName: "Squat",
		Sets:         1,
		Reps:         1,
	})
	if err == nil {
		t.Fatal("patient B must not log against patient A's plan")
	}

	_, err = qry.ListEntries(context.Background(), queries.ListExerciseEntriesQuery{
		PatientID: b.ID.String(),
		PlanID:    plan.ID.String(),
	})
	if err == nil {
		t.Fatal("patient B must not list patient A's plan entries")
	}

	_, err = cmds.UpdatePlan(context.Background(), UpdateExercisePlanCommand{
		PatientID: b.ID.String(),
		PlanID:    plan.ID.String(),
		Name:      "Hijacked",
		Status:    domain.ExercisePlanStatusArchived,
	})
	if err == nil {
		t.Fatal("patient B must not update patient A's plan")
	}

	if err := cmds.DeletePlan(context.Background(), DeleteExercisePlanCommand{
		PatientID: b.ID.String(),
		PlanID:    plan.ID.String(),
	}); err == nil {
		t.Fatal("patient B must not delete patient A's plan")
	}

	// Original patient can still use the plan.
	if _, err := cmds.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:    a.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   author,
		ExerciseName: "Squat",
		Sets:         2,
		Reps:         8,
	}); err != nil {
		t.Fatalf("patient A should still write entries: %v", err)
	}
}

func TestExerciseLog_DeletePlanRemovesEntries(t *testing.T) {
	patient, handler, authorID := newExerciseTestEnv(t)
	plan, err := handler.CreatePlan(context.Background(), CreateExercisePlanCommand{
		PatientID: patient.ID.String(),
		CreatedBy: authorID,
		Name:      "Temp",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}
	entry, err := handler.CreateEntry(context.Background(), CreateExerciseEntryCommand{
		PatientID:    patient.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   authorID,
		ExerciseName: "Bridge",
		Sets:         2,
		Reps:         10,
	})
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if err := handler.DeletePlan(context.Background(), DeleteExercisePlanCommand{
		PatientID: patient.ID.String(),
		PlanID:    plan.ID.String(),
	}); err != nil {
		t.Fatalf("delete plan: %v", err)
	}
	if err := handler.DeleteEntry(context.Background(), DeleteExerciseEntryCommand{
		PatientID: patient.ID.String(),
		PlanID:    plan.ID.String(),
		EntryID:   entry.ID.String(),
	}); err == nil {
		t.Fatal("expected entry gone after plan delete")
	}
}
