package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/google/uuid"
)

func TestMemoryExerciseLogRepository_PlanAndEntryLifecycle(t *testing.T) {
	patientRepo := NewMemoryRepository()
	patient := domain.NewPatient("Mem", "Repo", time.Date(1988, 2, 2, 0, 0, 0, 0, time.UTC), domain.GenderOther)
	if err := patientRepo.Create(context.Background(), patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}

	repo := NewMemoryExerciseLogRepository(patientRepo)
	authorID := uuid.New()
	plan := domain.NewExercisePlan(patient.ID, authorID, "Ankle", "PT")
	plan.CreatedByName = "Therapist"

	if err := repo.CreatePlan(context.Background(), plan); err != nil {
		t.Fatalf("create plan: %v", err)
	}

	got, err := repo.GetPlanByID(context.Background(), patient.ID.String(), plan.ID.String())
	if err != nil || got.Name != "Ankle" {
		t.Fatalf("get plan: err=%v plan=%+v", err, got)
	}

	plans, err := repo.ListPlansByPatientID(context.Background(), patient.ID.String())
	if err != nil || len(plans) != 1 {
		t.Fatalf("list plans: err=%v len=%d", err, len(plans))
	}

	now := time.Now().UTC()
	plan.Name = "Ankle v2"
	plan.Status = domain.ExercisePlanStatusCompleted
	plan.UpdatedAt = &now
	if err := repo.UpdatePlan(context.Background(), plan); err != nil {
		t.Fatalf("update plan: %v", err)
	}

	diff := 3
	entry := domain.NewExerciseLogEntry(
		plan.ID, patient.ID, authorID, time.Now().UTC(),
		"Calf raise", 3, 15, nil, "bodyweight", &diff, "steady",
	)
	entry.RecordedByName = "Therapist"
	if err := repo.CreateEntry(context.Background(), entry); err != nil {
		t.Fatalf("create entry: %v", err)
	}

	entries, err := repo.ListEntriesByPlanID(context.Background(), patient.ID.String(), plan.ID.String())
	if err != nil || len(entries) != 1 {
		t.Fatalf("list entries: err=%v len=%d", err, len(entries))
	}

	entry.Reps = 20
	entry.UpdatedAt = &now
	if err := repo.UpdateEntry(context.Background(), entry); err != nil {
		t.Fatalf("update entry: %v", err)
	}
	updated, err := repo.GetEntryByID(context.Background(), patient.ID.String(), plan.ID.String(), entry.ID.String())
	if err != nil || updated.Reps != 20 {
		t.Fatalf("get updated entry: err=%v entry=%+v", err, updated)
	}

	if err := repo.DeleteEntry(context.Background(), patient.ID.String(), plan.ID.String(), entry.ID.String()); err != nil {
		t.Fatalf("delete entry: %v", err)
	}
	if _, err := repo.GetEntryByID(context.Background(), patient.ID.String(), plan.ID.String(), entry.ID.String()); err == nil {
		t.Fatal("expected entry not found after delete")
	}

	// Recreate entry then delete plan (cascades in memory)
	entry2 := domain.NewExerciseLogEntry(
		plan.ID, patient.ID, authorID, time.Now().UTC(),
		"Calf raise", 2, 10, nil, "", nil, "",
	)
	if err := repo.CreateEntry(context.Background(), entry2); err != nil {
		t.Fatalf("create entry2: %v", err)
	}
	if err := repo.DeletePlan(context.Background(), patient.ID.String(), plan.ID.String()); err != nil {
		t.Fatalf("delete plan: %v", err)
	}
	entries, err = repo.ListEntriesByPlanID(context.Background(), patient.ID.String(), plan.ID.String())
	if err == nil {
		t.Fatalf("expected missing plan error, got %d entries", len(entries))
	}
}

func TestMemoryExerciseLogRepository_ListEntriesOrderedByPerformedAt(t *testing.T) {
	patientRepo := NewMemoryRepository()
	patient := domain.NewPatient("Ord", "Er", time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC), domain.GenderMale)
	if err := patientRepo.Create(context.Background(), patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}
	repo := NewMemoryExerciseLogRepository(patientRepo)
	authorID := uuid.New()
	plan := domain.NewExercisePlan(patient.ID, authorID, "Order", "")
	if err := repo.CreatePlan(context.Background(), plan); err != nil {
		t.Fatalf("create plan: %v", err)
	}

	older := domain.NewExerciseLogEntry(plan.ID, patient.ID, authorID, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), "A", 1, 1, nil, "", nil, "")
	newer := domain.NewExerciseLogEntry(plan.ID, patient.ID, authorID, time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), "A", 1, 2, nil, "", nil, "")
	if err := repo.CreateEntry(context.Background(), older); err != nil {
		t.Fatalf("create older: %v", err)
	}
	if err := repo.CreateEntry(context.Background(), newer); err != nil {
		t.Fatalf("create newer: %v", err)
	}

	entries, err := repo.ListEntriesByPlanID(context.Background(), patient.ID.String(), plan.ID.String())
	if err != nil || len(entries) != 2 {
		t.Fatalf("list: err=%v len=%d", err, len(entries))
	}
	if !entries[0].PerformedAt.After(entries[1].PerformedAt) {
		t.Fatalf("expected newest first, got %v then %v", entries[0].PerformedAt, entries[1].PerformedAt)
	}
}
