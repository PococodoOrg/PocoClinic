package domain

import "context"

// SettingsRepository stores clinic-wide patient configuration.
type SettingsRepository interface {
	GetPatientFieldRequirements(ctx context.Context) (PatientFieldRequirements, error)
	SavePatientFieldRequirements(ctx context.Context, requirements PatientFieldRequirements) error
}
