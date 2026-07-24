package queries

import (
	"context"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
)

type ListExercisePlansQuery struct {
	PatientID string
}

type ListExerciseEntriesQuery struct {
	PatientID string
	PlanID    string
}

type ExerciseLogQueryHandler struct {
	repo domain.ExerciseLogRepository
}

func NewExerciseLogQueryHandler(repo domain.ExerciseLogRepository) *ExerciseLogQueryHandler {
	return &ExerciseLogQueryHandler{repo: repo}
}

func (h *ExerciseLogQueryHandler) ListPlans(ctx context.Context, q ListExercisePlansQuery) ([]*domain.ExercisePlan, error) {
	plans, err := h.repo.ListPlansByPatientID(ctx, q.PatientID)
	if err != nil {
		return nil, err
	}
	if plans == nil {
		return []*domain.ExercisePlan{}, nil
	}
	return plans, nil
}

func (h *ExerciseLogQueryHandler) ListEntries(ctx context.Context, q ListExerciseEntriesQuery) ([]*domain.ExerciseLogEntry, error) {
	if _, err := h.repo.GetPlanByID(ctx, q.PatientID, q.PlanID); err != nil {
		return nil, err
	}
	entries, err := h.repo.ListEntriesByPlanID(ctx, q.PatientID, q.PlanID)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		return []*domain.ExerciseLogEntry{}, nil
	}
	return entries, nil
}
