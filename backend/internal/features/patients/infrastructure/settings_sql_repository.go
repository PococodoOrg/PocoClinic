package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

// SQLSettingsRepository stores clinic settings in SQLite.
type SQLSettingsRepository struct {
	pool *database.DB
}

func NewSQLSettingsRepository(pool *database.DB) *SQLSettingsRepository {
	return &SQLSettingsRepository{pool: pool}
}

func (r *SQLSettingsRepository) GetPatientFieldRequirements(ctx context.Context) (domain.PatientFieldRequirements, error) {
	var raw string
	err := r.pool.QueryRow(ctx,
		`SELECT value_json FROM clinic_settings WHERE key = $1`,
		domain.PatientFieldRequirementsSettingKey(),
	).Scan(&raw)
	if err != nil {
		return domain.DefaultPatientFieldRequirements(), nil
	}

	var reqs domain.PatientFieldRequirements
	if err := json.Unmarshal([]byte(raw), &reqs); err != nil {
		return domain.DefaultPatientFieldRequirements(), nil
	}
	return domain.NormalizePatientFieldRequirements(reqs), nil
}

func (r *SQLSettingsRepository) SavePatientFieldRequirements(ctx context.Context, requirements domain.PatientFieldRequirements) error {
	requirements = domain.NormalizePatientFieldRequirements(requirements)
	payload, err := json.Marshal(requirements)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO clinic_settings (key, value_json, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT(key) DO UPDATE SET
			value_json = excluded.value_json,
			updated_at = excluded.updated_at
	`, domain.PatientFieldRequirementsSettingKey(), string(payload), time.Now().UTC())
	return err
}
