package infrastructure

import (
	"context"
	"sync"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/google/uuid"
)

type MemoryExerciseLogRepository struct {
	mu       sync.RWMutex
	plans    map[uuid.UUID]*domain.ExercisePlan
	entries  map[uuid.UUID]*domain.ExerciseLogEntry
	patients domain.PatientRepository
}

func NewMemoryExerciseLogRepository(patients domain.PatientRepository) *MemoryExerciseLogRepository {
	return &MemoryExerciseLogRepository{
		plans:    make(map[uuid.UUID]*domain.ExercisePlan),
		entries:  make(map[uuid.UUID]*domain.ExerciseLogEntry),
		patients: patients,
	}
}

func (r *MemoryExerciseLogRepository) ensurePatient(ctx context.Context, patientID string) error {
	if _, err := r.patients.GetByID(ctx, patientID); err != nil {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Patient not found")
	}
	return nil
}

func (r *MemoryExerciseLogRepository) CreatePlan(ctx context.Context, plan *domain.ExercisePlan) error {
	if err := r.ensurePatient(ctx, plan.PatientID.String()); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *plan
	if copy.CreatedByName == "" {
		copy.CreatedByName = "Staff"
	}
	r.plans[plan.ID] = &copy
	plan.CreatedByName = copy.CreatedByName
	return nil
}

func (r *MemoryExerciseLogRepository) GetPlanByID(ctx context.Context, patientID, planID string) (*domain.ExercisePlan, error) {
	if err := r.ensurePatient(ctx, patientID); err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}
	id, err := uuid.Parse(planID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid plan ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	plan, ok := r.plans[id]
	if !ok || plan.PatientID != pid {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise plan not found")
	}
	copy := *plan
	return &copy, nil
}

func (r *MemoryExerciseLogRepository) ListPlansByPatientID(ctx context.Context, patientID string) ([]*domain.ExercisePlan, error) {
	if err := r.ensurePatient(ctx, patientID); err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.ExercisePlan, 0)
	for _, plan := range r.plans {
		if plan.PatientID != pid {
			continue
		}
		copy := *plan
		out = append(out, &copy)
	}
	return out, nil
}

func (r *MemoryExerciseLogRepository) UpdatePlan(ctx context.Context, plan *domain.ExercisePlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.plans[plan.ID]
	if !ok {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise plan not found")
	}
	copy := *plan
	if copy.CreatedByName == "" {
		copy.CreatedByName = existing.CreatedByName
	}
	r.plans[plan.ID] = &copy
	plan.CreatedByName = copy.CreatedByName
	return nil
}

func (r *MemoryExerciseLogRepository) DeletePlan(ctx context.Context, patientID, planID string) error {
	if _, err := r.GetPlanByID(ctx, patientID, planID); err != nil {
		return err
	}
	id, _ := uuid.Parse(planID)
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.plans, id)
	for entryID, entry := range r.entries {
		if entry.PlanID == id {
			delete(r.entries, entryID)
		}
	}
	return nil
}

func (r *MemoryExerciseLogRepository) CreateEntry(ctx context.Context, entry *domain.ExerciseLogEntry) error {
	if _, err := r.GetPlanByID(ctx, entry.PatientID.String(), entry.PlanID.String()); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *entry
	if copy.RecordedByName == "" {
		copy.RecordedByName = "Staff"
	}
	r.entries[entry.ID] = &copy
	entry.RecordedByName = copy.RecordedByName
	return nil
}

func (r *MemoryExerciseLogRepository) GetEntryByID(ctx context.Context, patientID, planID, entryID string) (*domain.ExerciseLogEntry, error) {
	if _, err := r.GetPlanByID(ctx, patientID, planID); err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid plan ID")
	}
	id, err := uuid.Parse(entryID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid entry ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.entries[id]
	if !ok || entry.PatientID != pid || entry.PlanID != planUUID {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise log entry not found")
	}
	copy := *entry
	return &copy, nil
}

func (r *MemoryExerciseLogRepository) ListEntriesByPlanID(ctx context.Context, patientID, planID string) ([]*domain.ExerciseLogEntry, error) {
	if _, err := r.GetPlanByID(ctx, patientID, planID); err != nil {
		return nil, err
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid plan ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.ExerciseLogEntry, 0)
	for _, entry := range r.entries {
		if entry.PlanID != planUUID {
			continue
		}
		copy := *entry
		out = append(out, &copy)
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].PerformedAt.After(out[i].PerformedAt) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

func (r *MemoryExerciseLogRepository) UpdateEntry(ctx context.Context, entry *domain.ExerciseLogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.entries[entry.ID]
	if !ok {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise log entry not found")
	}
	copy := *entry
	if copy.RecordedByName == "" {
		copy.RecordedByName = existing.RecordedByName
	}
	r.entries[entry.ID] = &copy
	entry.RecordedByName = copy.RecordedByName
	return nil
}

func (r *MemoryExerciseLogRepository) DeleteEntry(ctx context.Context, patientID, planID, entryID string) error {
	if _, err := r.GetEntryByID(ctx, patientID, planID, entryID); err != nil {
		return err
	}
	id, _ := uuid.Parse(entryID)
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, id)
	return nil
}

var _ domain.ExerciseLogRepository = (*MemoryExerciseLogRepository)(nil)
