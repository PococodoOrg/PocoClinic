package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
)

// MemoryRepository is a simple in-memory implementation of the PatientRepository interface
type MemoryRepository struct {
	patients map[string]*domain.Patient
	mu       sync.RWMutex
}

// NewMemoryRepository creates a new in-memory patient repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		patients: make(map[string]*domain.Patient),
	}
}

// Create adds a new patient to the repository
func (r *MemoryRepository) Create(ctx context.Context, patient *domain.Patient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.patients[patient.ID.String()]; exists {
		return fmt.Errorf("patient with ID %s already exists", patient.ID)
	}

	r.patients[patient.ID.String()] = patient
	return nil
}

// Update modifies an existing patient in the repository
func (r *MemoryRepository) Update(ctx context.Context, patient *domain.Patient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.patients[patient.ID.String()]; !exists {
		return fmt.Errorf("patient with ID %s not found", patient.ID)
	}

	r.patients[patient.ID.String()] = patient
	return nil
}

// Delete removes a patient from the repository
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.patients[id]; !exists {
		return fmt.Errorf("patient with ID %s not found", id)
	}

	delete(r.patients, id)
	return nil
}

// GetByID retrieves a patient by their ID
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	patient, exists := r.patients[id]
	if !exists {
		return nil, fmt.Errorf("patient with ID %s not found", id)
	}

	return patient, nil
}

// List returns all patients in the repository
func (r *MemoryRepository) List(ctx context.Context) ([]*domain.Patient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	patients := make([]*domain.Patient, 0, len(r.patients))
	for _, patient := range r.patients {
		patients = append(patients, patient)
	}

	return patients, nil
}

// ListPaginated returns a paginated list of patients with optional search
func (r *MemoryRepository) ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*domain.Patient, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Filter patients by search term if provided
	var filteredPatients []*domain.Patient
	for _, patient := range r.patients {
		if search == "" || strings.Contains(strings.ToLower(patient.FullName()), strings.ToLower(search)) {
			filteredPatients = append(filteredPatients, patient)
		}
	}

	// Calculate pagination
	totalCount := int64(len(filteredPatients))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(filteredPatients) {
		return []*domain.Patient{}, totalCount, nil
	}
	if end > len(filteredPatients) {
		end = len(filteredPatients)
	}

	return filteredPatients[start:end], totalCount, nil
}
