package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

type SQLExerciseLogRepository struct {
	pool *database.DB
}

func NewSQLExerciseLogRepository(pool *database.DB) *SQLExerciseLogRepository {
	return &SQLExerciseLogRepository{pool: pool}
}

func (r *SQLExerciseLogRepository) CreatePlan(ctx context.Context, plan *domain.ExercisePlan) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO exercise_plans (id, patient_id, name, description, status, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, plan.ID, plan.PatientID, plan.Name, plan.Description, plan.Status, plan.CreatedBy, plan.CreatedAt)
	if err != nil {
		return err
	}
	_ = r.pool.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, plan.CreatedBy).Scan(&plan.CreatedByName)
	return nil
}

func (r *SQLExerciseLogRepository) GetPlanByID(ctx context.Context, patientID, planID string) (*domain.ExercisePlan, error) {
	plan := &domain.ExercisePlan{}
	var updatedAt *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT p.id, p.patient_id, p.name, p.description, p.status, p.created_by, u.name, p.created_at, p.updated_at
		FROM exercise_plans p
		JOIN users u ON u.id = p.created_by
		WHERE p.patient_id = $1 AND p.id = $2
	`, patientID, planID).Scan(
		&plan.ID, &plan.PatientID, &plan.Name, &plan.Description, &plan.Status,
		&plan.CreatedBy, &plan.CreatedByName, database.Time(&plan.CreatedAt), database.NullTime(&updatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise plan not found")
		}
		return nil, err
	}
	plan.UpdatedAt = updatedAt
	return plan, nil
}

func (r *SQLExerciseLogRepository) ListPlansByPatientID(ctx context.Context, patientID string) ([]*domain.ExercisePlan, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.patient_id, p.name, p.description, p.status, p.created_by, u.name, p.created_at, p.updated_at
		FROM exercise_plans p
		JOIN users u ON u.id = p.created_by
		WHERE p.patient_id = $1
		ORDER BY
			CASE p.status WHEN 'active' THEN 0 WHEN 'completed' THEN 1 ELSE 2 END,
			p.created_at DESC
	`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]*domain.ExercisePlan, 0)
	for rows.Next() {
		plan := &domain.ExercisePlan{}
		var updatedAt *time.Time
		if err := rows.Scan(
			&plan.ID, &plan.PatientID, &plan.Name, &plan.Description, &plan.Status,
			&plan.CreatedBy, &plan.CreatedByName, database.Time(&plan.CreatedAt), database.NullTime(&updatedAt),
		); err != nil {
			return nil, err
		}
		plan.UpdatedAt = updatedAt
		plans = append(plans, plan)
	}
	return plans, rows.Err()
}

func (r *SQLExerciseLogRepository) UpdatePlan(ctx context.Context, plan *domain.ExercisePlan) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE exercise_plans
		SET name = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5 AND patient_id = $6
	`, plan.Name, plan.Description, plan.Status, plan.UpdatedAt, plan.ID, plan.PatientID)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, plan.CreatedBy).Scan(&plan.CreatedByName)
}

func (r *SQLExerciseLogRepository) DeletePlan(ctx context.Context, patientID, planID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM exercise_plans WHERE patient_id = $1 AND id = $2
	`, patientID, planID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise plan not found")
	}
	return nil
}

func (r *SQLExerciseLogRepository) CreateEntry(ctx context.Context, entry *domain.ExerciseLogEntry) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO exercise_log_entries (
			id, plan_id, patient_id, performed_at, exercise_name, sets, reps,
			duration_seconds, resistance, difficulty, notes, recorded_by, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`,
		entry.ID, entry.PlanID, entry.PatientID, entry.PerformedAt, entry.ExerciseName,
		entry.Sets, entry.Reps, entry.DurationSeconds, entry.Resistance, entry.Difficulty,
		entry.Notes, entry.RecordedBy, entry.CreatedAt,
	)
	if err != nil {
		return err
	}
	_ = r.pool.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, entry.RecordedBy).Scan(&entry.RecordedByName)
	return nil
}

func (r *SQLExerciseLogRepository) GetEntryByID(ctx context.Context, patientID, planID, entryID string) (*domain.ExerciseLogEntry, error) {
	entry := &domain.ExerciseLogEntry{}
	var updatedAt *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT e.id, e.plan_id, e.patient_id, e.performed_at, e.exercise_name, e.sets, e.reps,
			e.duration_seconds, e.resistance, e.difficulty, e.notes, e.recorded_by, u.name, e.created_at, e.updated_at
		FROM exercise_log_entries e
		JOIN users u ON u.id = e.recorded_by
		WHERE e.patient_id = $1 AND e.plan_id = $2 AND e.id = $3
	`, patientID, planID, entryID).Scan(
		&entry.ID, &entry.PlanID, &entry.PatientID, database.Time(&entry.PerformedAt), &entry.ExerciseName,
		&entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Resistance, &entry.Difficulty,
		&entry.Notes, &entry.RecordedBy, &entry.RecordedByName, database.Time(&entry.CreatedAt), database.NullTime(&updatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise log entry not found")
		}
		return nil, err
	}
	entry.UpdatedAt = updatedAt
	return entry, nil
}

func (r *SQLExerciseLogRepository) ListEntriesByPlanID(ctx context.Context, patientID, planID string) ([]*domain.ExerciseLogEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT e.id, e.plan_id, e.patient_id, e.performed_at, e.exercise_name, e.sets, e.reps,
			e.duration_seconds, e.resistance, e.difficulty, e.notes, e.recorded_by, u.name, e.created_at, e.updated_at
		FROM exercise_log_entries e
		JOIN users u ON u.id = e.recorded_by
		WHERE e.patient_id = $1 AND e.plan_id = $2
		ORDER BY e.performed_at DESC, e.created_at DESC
	`, patientID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]*domain.ExerciseLogEntry, 0)
	for rows.Next() {
		entry := &domain.ExerciseLogEntry{}
		var updatedAt *time.Time
		if err := rows.Scan(
			&entry.ID, &entry.PlanID, &entry.PatientID, database.Time(&entry.PerformedAt), &entry.ExerciseName,
			&entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Resistance, &entry.Difficulty,
			&entry.Notes, &entry.RecordedBy, &entry.RecordedByName, database.Time(&entry.CreatedAt), database.NullTime(&updatedAt),
		); err != nil {
			return nil, err
		}
		entry.UpdatedAt = updatedAt
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (r *SQLExerciseLogRepository) UpdateEntry(ctx context.Context, entry *domain.ExerciseLogEntry) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE exercise_log_entries
		SET performed_at = $1, exercise_name = $2, sets = $3, reps = $4,
			duration_seconds = $5, resistance = $6, difficulty = $7, notes = $8, updated_at = $9
		WHERE id = $10 AND plan_id = $11 AND patient_id = $12
	`,
		entry.PerformedAt, entry.ExerciseName, entry.Sets, entry.Reps,
		entry.DurationSeconds, entry.Resistance, entry.Difficulty, entry.Notes, entry.UpdatedAt,
		entry.ID, entry.PlanID, entry.PatientID,
	)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, entry.RecordedBy).Scan(&entry.RecordedByName)
}

func (r *SQLExerciseLogRepository) DeleteEntry(ctx context.Context, patientID, planID, entryID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM exercise_log_entries
		WHERE patient_id = $1 AND plan_id = $2 AND id = $3
	`, patientID, planID, entryID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Exercise log entry not found")
	}
	return nil
}

var _ domain.ExerciseLogRepository = (*SQLExerciseLogRepository)(nil)
