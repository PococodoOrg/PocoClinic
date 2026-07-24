package domain

import "context"

// ExerciseLogRepository persists PT exercise plans and session entries.
type ExerciseLogRepository interface {
	CreatePlan(ctx context.Context, plan *ExercisePlan) error
	GetPlanByID(ctx context.Context, patientID, planID string) (*ExercisePlan, error)
	ListPlansByPatientID(ctx context.Context, patientID string) ([]*ExercisePlan, error)
	UpdatePlan(ctx context.Context, plan *ExercisePlan) error
	DeletePlan(ctx context.Context, patientID, planID string) error

	CreateEntry(ctx context.Context, entry *ExerciseLogEntry) error
	GetEntryByID(ctx context.Context, patientID, planID, entryID string) (*ExerciseLogEntry, error)
	ListEntriesByPlanID(ctx context.Context, patientID, planID string) ([]*ExerciseLogEntry, error)
	UpdateEntry(ctx context.Context, entry *ExerciseLogEntry) error
	DeleteEntry(ctx context.Context, patientID, planID, entryID string) error
}
