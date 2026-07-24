package queries

import (
	"context"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/commands"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	patientinfra "github.com/PococodoOrg/PocoClinic/internal/features/patients/infrastructure"
	"github.com/google/uuid"
)

func TestExerciseLogQueryHandler(t *testing.T) {
	patientRepo := patientinfra.NewMemoryRepository()
	patient := domain.NewPatient("Ann", "Query", time.Date(1990, 3, 4, 0, 0, 0, 0, time.UTC), domain.GenderFemale)
	if err := patientRepo.Create(context.Background(), patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}

	repo := patientinfra.NewMemoryExerciseLogRepository(patientRepo)
	cmd := commands.NewExerciseLogCommandHandler(repo, patientRepo)
	qry := NewExerciseLogQueryHandler(repo)
	authorID := uuid.New()

	plans, err := qry.ListPlans(context.Background(), ListExercisePlansQuery{PatientID: patient.ID.String()})
	if err != nil {
		t.Fatalf("list empty plans: %v", err)
	}
	if len(plans) != 0 {
		t.Fatalf("expected empty plans, got %d", len(plans))
	}

	plan, err := cmd.CreatePlan(context.Background(), commands.CreateExercisePlanCommand{
		PatientID: patient.ID.String(),
		CreatedBy: authorID,
		Name:      "Balance",
	})
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}

	if _, err := cmd.CreateEntry(context.Background(), commands.CreateExerciseEntryCommand{
		PatientID:    patient.ID.String(),
		PlanID:       plan.ID.String(),
		RecordedBy:   authorID,
		ExerciseName: "Single-leg stand",
		Sets:         2,
		Reps:         30,
	}); err != nil {
		t.Fatalf("create entry: %v", err)
	}

	plans, err = qry.ListPlans(context.Background(), ListExercisePlansQuery{PatientID: patient.ID.String()})
	if err != nil || len(plans) != 1 {
		t.Fatalf("list plans: err=%v len=%d", err, len(plans))
	}

	entries, err := qry.ListEntries(context.Background(), ListExerciseEntriesQuery{
		PatientID: patient.ID.String(),
		PlanID:    plan.ID.String(),
	})
	if err != nil || len(entries) != 1 {
		t.Fatalf("list entries: err=%v len=%d", err, len(entries))
	}

	_, err = qry.ListEntries(context.Background(), ListExerciseEntriesQuery{
		PatientID: patient.ID.String(),
		PlanID:    uuid.New().String(),
	})
	if err == nil {
		t.Fatal("expected not found for unknown plan")
	}
}
