package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	EventLoginSuccess     = "auth.login.success"
	EventLoginFailure     = "auth.login.failure"
	EventLoginLocked      = "auth.login.locked"
	EventTokenRefresh     = "auth.token.refresh"
	EventLogout           = "auth.logout"
	EventBadgeReissued    = "auth.badge.reissued"
	EventPINChanged       = "auth.pin.changed"
	EventUserCreated      = "user.created"
	EventUserUpdated      = "user.updated"
	EventUserDeleted      = "user.deleted"
	EventUserUnlocked     = "user.unlocked"
	EventPatientCreated  = "patient.created"
	EventPatientUpdated  = "patient.updated"
	EventPatientDeleted  = "patient.deleted"
	EventPatientViewed   = "patient.viewed"
	EventPatientNoteCreated = "patient.note.created"
	EventPatientNoteUpdated = "patient.note.updated"
	EventPatientNoteDeleted = "patient.note.deleted"
	EventPatientDocumentUploaded = "patient.document.uploaded"
	EventPatientDocumentViewed   = "patient.document.viewed"
	EventPatientDocumentDeleted  = "patient.document.deleted"
	EventPatientExercisePlanCreated  = "patient.exercise_plan.created"
	EventPatientExercisePlanUpdated  = "patient.exercise_plan.updated"
	EventPatientExercisePlanDeleted  = "patient.exercise_plan.deleted"
	EventPatientExerciseEntryCreated = "patient.exercise_entry.created"
	EventPatientExerciseEntryUpdated = "patient.exercise_entry.updated"
	EventPatientExerciseEntryDeleted = "patient.exercise_entry.deleted"
	EventFormSubmissionCreated   = "form.submission.created"
	EventBackupCreated    = "system.backup.created"
	EventBackupRestored   = "system.backup.restored"
	EventHealthCheckRun   = "system.healthcheck.run"
	EventComplianceCheckRun = "system.compliance.run"
	EventAuditPurged      = "system.audit.purged"
)

// Event represents a HIPAA-relevant audit log entry.
type Event struct {
	ID           uuid.UUID
	EventType    string
	UserID       *uuid.UUID
	ResourceType string
	ResourceID   string
	IPAddress    string
	UserAgent    string
	Details      map[string]string
	Success      bool
	CreatedAt    time.Time
}

// Repository persists audit events.
type Repository interface {
	Create(ctx context.Context, event *Event) error
	ListForUser(ctx context.Context, userID string, page, pageSize int) ([]*Event, int64, error)
	ListRecent(ctx context.Context, filter ListFilter) ([]*Event, int64, error)
}

// ListFilter narrows a system-wide audit query.
type ListFilter struct {
	Page      int
	PageSize  int
	EventType string
	UserID    string
	SinceDays int
}

// Logger records audit events.
type Logger interface {
	Log(ctx context.Context, event Event)
}
