package commands

import (
	"context"
	"strings"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/google/uuid"
)

type CreateExercisePlanCommand struct {
	PatientID   string `json:"-"`
	CreatedBy   uuid.UUID
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateExercisePlanCommand struct {
	PatientID   string `json:"-"`
	PlanID      string `json:"-"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required"`
}

type DeleteExercisePlanCommand struct {
	PatientID string
	PlanID    string
}

type CreateExerciseEntryCommand struct {
	PatientID       string `json:"-"`
	PlanID          string `json:"-"`
	RecordedBy      uuid.UUID
	PerformedAt     *time.Time `json:"performedAt"`
	ExerciseName    string     `json:"exerciseName" binding:"required"`
	Sets            int        `json:"sets"`
	Reps            int        `json:"reps"`
	DurationSeconds *int       `json:"durationSeconds"`
	Resistance      string     `json:"resistance"`
	Difficulty      *int       `json:"difficulty"`
	Notes           string     `json:"notes"`
}

type UpdateExerciseEntryCommand struct {
	PatientID       string `json:"-"`
	PlanID          string `json:"-"`
	EntryID         string `json:"-"`
	PerformedAt     *time.Time `json:"performedAt"`
	ExerciseName    string     `json:"exerciseName" binding:"required"`
	Sets            int        `json:"sets"`
	Reps            int        `json:"reps"`
	DurationSeconds *int       `json:"durationSeconds"`
	Resistance      string     `json:"resistance"`
	Difficulty      *int       `json:"difficulty"`
	Notes           string     `json:"notes"`
}

type DeleteExerciseEntryCommand struct {
	PatientID string
	PlanID    string
	EntryID   string
}

type ExerciseLogCommandHandler struct {
	repo        domain.ExerciseLogRepository
	patientRepo domain.GetPatientRepository
}

func NewExerciseLogCommandHandler(
	repo domain.ExerciseLogRepository,
	patientRepo domain.GetPatientRepository,
) *ExerciseLogCommandHandler {
	return &ExerciseLogCommandHandler{repo: repo, patientRepo: patientRepo}
}

func (h *ExerciseLogCommandHandler) CreatePlan(ctx context.Context, cmd CreateExercisePlanCommand) (*domain.ExercisePlan, error) {
	name, description, err := validatePlanFields(cmd.Name, cmd.Description)
	if err != nil {
		return nil, err
	}
	if _, err := h.patientRepo.GetByID(ctx, cmd.PatientID); err != nil {
		return nil, err
	}
	patientID, err := uuid.Parse(cmd.PatientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}

	plan := domain.NewExercisePlan(patientID, cmd.CreatedBy, name, description)
	if err := h.repo.CreatePlan(ctx, plan); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to save exercise plan")
	}
	return plan, nil
}

func (h *ExerciseLogCommandHandler) UpdatePlan(ctx context.Context, cmd UpdateExercisePlanCommand) (*domain.ExercisePlan, error) {
	name, description, err := validatePlanFields(cmd.Name, cmd.Description)
	if err != nil {
		return nil, err
	}
	if !domain.ValidExercisePlanStatus(cmd.Status) {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid plan status")
	}

	plan, err := h.repo.GetPlanByID(ctx, cmd.PatientID, cmd.PlanID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	plan.Name = name
	plan.Description = description
	plan.Status = cmd.Status
	plan.UpdatedAt = &now
	if err := h.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to update exercise plan")
	}
	return plan, nil
}

func (h *ExerciseLogCommandHandler) DeletePlan(ctx context.Context, cmd DeleteExercisePlanCommand) error {
	if _, err := h.repo.GetPlanByID(ctx, cmd.PatientID, cmd.PlanID); err != nil {
		return err
	}
	if err := h.repo.DeletePlan(ctx, cmd.PatientID, cmd.PlanID); err != nil {
		return err
	}
	return nil
}

func (h *ExerciseLogCommandHandler) CreateEntry(ctx context.Context, cmd CreateExerciseEntryCommand) (*domain.ExerciseLogEntry, error) {
	fields, err := validateEntryFields(cmd.ExerciseName, cmd.Sets, cmd.Reps, cmd.DurationSeconds, cmd.Resistance, cmd.Difficulty, cmd.Notes)
	if err != nil {
		return nil, err
	}

	plan, err := h.repo.GetPlanByID(ctx, cmd.PatientID, cmd.PlanID)
	if err != nil {
		return nil, err
	}

	performedAt := time.Now()
	if cmd.PerformedAt != nil {
		performedAt = cmd.PerformedAt.UTC()
	}

	entry := domain.NewExerciseLogEntry(
		plan.ID,
		plan.PatientID,
		cmd.RecordedBy,
		performedAt,
		fields.exerciseName,
		fields.sets,
		fields.reps,
		fields.durationSeconds,
		fields.resistance,
		fields.difficulty,
		fields.notes,
	)
	if err := h.repo.CreateEntry(ctx, entry); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to save exercise log entry")
	}
	return entry, nil
}

func (h *ExerciseLogCommandHandler) UpdateEntry(ctx context.Context, cmd UpdateExerciseEntryCommand) (*domain.ExerciseLogEntry, error) {
	fields, err := validateEntryFields(cmd.ExerciseName, cmd.Sets, cmd.Reps, cmd.DurationSeconds, cmd.Resistance, cmd.Difficulty, cmd.Notes)
	if err != nil {
		return nil, err
	}

	entry, err := h.repo.GetEntryByID(ctx, cmd.PatientID, cmd.PlanID, cmd.EntryID)
	if err != nil {
		return nil, err
	}

	if cmd.PerformedAt != nil {
		entry.PerformedAt = cmd.PerformedAt.UTC()
	}
	now := time.Now()
	entry.ExerciseName = fields.exerciseName
	entry.Sets = fields.sets
	entry.Reps = fields.reps
	entry.DurationSeconds = fields.durationSeconds
	entry.Resistance = fields.resistance
	entry.Difficulty = fields.difficulty
	entry.Notes = fields.notes
	entry.UpdatedAt = &now

	if err := h.repo.UpdateEntry(ctx, entry); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to update exercise log entry")
	}
	return entry, nil
}

func (h *ExerciseLogCommandHandler) DeleteEntry(ctx context.Context, cmd DeleteExerciseEntryCommand) error {
	if _, err := h.repo.GetEntryByID(ctx, cmd.PatientID, cmd.PlanID, cmd.EntryID); err != nil {
		return err
	}
	return h.repo.DeleteEntry(ctx, cmd.PatientID, cmd.PlanID, cmd.EntryID)
}

func validatePlanFields(name, description string) (string, string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", "", pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Plan name is required")
	}
	if len(trimmedName) > domain.MaxExercisePlanNameLength {
		return "", "", pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Plan name is too long")
	}
	trimmedDesc := strings.TrimSpace(description)
	if len(trimmedDesc) > domain.MaxExercisePlanDescriptionLength {
		return "", "", pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Plan description is too long")
	}
	return trimmedName, trimmedDesc, nil
}

type entryFields struct {
	exerciseName    string
	sets            int
	reps            int
	durationSeconds *int
	resistance      string
	difficulty      *int
	notes           string
}

func validateEntryFields(
	exerciseName string,
	sets, reps int,
	durationSeconds *int,
	resistance string,
	difficulty *int,
	notes string,
) (entryFields, error) {
	name := strings.TrimSpace(exerciseName)
	if name == "" {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Exercise name is required")
	}
	if len(name) > domain.MaxExerciseNameLength {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Exercise name is too long")
	}
	if sets < 0 || reps < 0 {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Sets and reps cannot be negative")
	}
	if durationSeconds != nil && *durationSeconds < 0 {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Duration cannot be negative")
	}
	if difficulty != nil && (*difficulty < 1 || *difficulty > 10) {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Difficulty must be between 1 and 10")
	}
	res := strings.TrimSpace(resistance)
	if len(res) > domain.MaxExerciseResistanceLength {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Resistance note is too long")
	}
	trimmedNotes := strings.TrimSpace(notes)
	if len(trimmedNotes) > domain.MaxExerciseEntryNotesLength {
		return entryFields{}, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Notes are too long")
	}
	return entryFields{
		exerciseName:    name,
		sets:            sets,
		reps:            reps,
		durationSeconds: durationSeconds,
		resistance:      res,
		difficulty:      difficulty,
		notes:           trimmedNotes,
	}, nil
}
