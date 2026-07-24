package infrastructure

import (
	"context"
	"sync"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
)

// MemorySettingsRepository stores clinic settings in memory for dev without SQLite.
type MemorySettingsRepository struct {
	mu   sync.RWMutex
	reqs domain.PatientFieldRequirements
}

func NewMemorySettingsRepository() *MemorySettingsRepository {
	defaults := domain.DefaultPatientFieldRequirements()
	return &MemorySettingsRepository{reqs: defaults}
}

func (r *MemorySettingsRepository) GetPatientFieldRequirements(ctx context.Context) (domain.PatientFieldRequirements, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.reqs, nil
}

func (r *MemorySettingsRepository) SavePatientFieldRequirements(ctx context.Context, requirements domain.PatientFieldRequirements) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reqs = domain.NormalizePatientFieldRequirements(requirements)
	return nil
}
