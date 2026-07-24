package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

// SQLRepository persists patients in SQLite.
type SQLRepository struct {
	pool *database.DB
}

// NewSQLRepository creates a SQL-backed patient repository.
func NewSQLRepository(pool *database.DB) *SQLRepository {
	return &SQLRepository{pool: pool}
}

// Create adds a new patient.
func (r *SQLRepository) Create(ctx context.Context, patient *domain.Patient) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO patients (
			id, first_name, last_name, middle_name, date_of_birth, gender,
			email, phone_number, height, weight,
			address_street, address_city, address_state, address_postal_code, address_country,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`,
		patient.ID,
		patient.FirstName,
		patient.LastName,
		nullString(patient.MiddleName),
		patient.DateOfBirth.Time(),
		patient.Gender,
		nullString(patient.Email),
		nullString(patient.PhoneNumber),
		nullFloat(patient.Height),
		nullFloat(patient.Weight),
		nullString(patient.Address.Street),
		nullString(patient.Address.City),
		nullString(patient.Address.State),
		nullString(patient.Address.PostalCode),
		nullString(patient.Address.Country),
		patient.CreatedAt,
		patient.UpdatedAt,
	)
	return err
}

// Update modifies an existing patient.
func (r *SQLRepository) Update(ctx context.Context, patient *domain.Patient) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE patients SET
			first_name = $2,
			last_name = $3,
			middle_name = $4,
			date_of_birth = $5,
			gender = $6,
			email = $7,
			phone_number = $8,
			height = $9,
			weight = $10,
			address_street = $11,
			address_city = $12,
			address_state = $13,
			address_postal_code = $14,
			address_country = $15,
			updated_at = $16
		WHERE id = $1
	`,
		patient.ID,
		patient.FirstName,
		patient.LastName,
		nullString(patient.MiddleName),
		patient.DateOfBirth.Time(),
		patient.Gender,
		nullString(patient.Email),
		nullString(patient.PhoneNumber),
		nullFloat(patient.Height),
		nullFloat(patient.Weight),
		nullString(patient.Address.Street),
		nullString(patient.Address.City),
		nullString(patient.Address.State),
		nullString(patient.Address.PostalCode),
		nullString(patient.Address.Country),
		patient.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Patient not found")
	}
	return nil
}

// Delete removes a patient.
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM patients WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("patient with ID %s not found", id)
	}
	return nil
}

// GetByID retrieves a patient by ID.
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	row := r.pool.QueryRow(ctx, patientSelectQuery+` WHERE id = $1`, id)
	patient, err := scanPatient(row)
	if errors.Is(err, database.ErrNoRows) {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Patient not found")
	}
	return patient, err
}

// List returns all patients.
func (r *SQLRepository) List(ctx context.Context) ([]*domain.Patient, error) {
	rows, err := r.pool.Query(ctx, patientSelectQuery+` ORDER BY last_name, first_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := make([]*domain.Patient, 0)
	for rows.Next() {
		patient, err := scanPatient(rows)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, rows.Err()
}

// ListPaginated returns a paginated list of patients with optional filters.
func (r *SQLRepository) ListPaginated(ctx context.Context, page, pageSize int, filter domain.PatientListFilter) ([]*domain.Patient, int64, error) {
	whereClause, whereArgs := buildPatientListWhere(filter)

	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM patients` + whereClause
	if err := r.pool.QueryRow(ctx, countQuery, whereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	listQuery := patientSelectQuery + whereClause
	listArgs := append([]any{}, whereArgs...)
	offset := (page - 1) * pageSize
	listQuery += fmt.Sprintf(` ORDER BY last_name, first_name LIMIT $%d OFFSET $%d`, len(listArgs)+1, len(listArgs)+2)
	listArgs = append(listArgs, pageSize, offset)

	rows, err := r.pool.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	patients := make([]*domain.Patient, 0)
	for rows.Next() {
		patient, err := scanPatient(rows)
		if err != nil {
			return nil, 0, err
		}
		patients = append(patients, patient)
	}
	return patients, totalCount, rows.Err()
}

const patientSelectQuery = `
	SELECT id, first_name, last_name, middle_name, date_of_birth, gender,
	       email, phone_number, height, weight,
	       address_street, address_city, address_state, address_postal_code, address_country,
	       created_at, updated_at
	FROM patients
`

func scanPatient(row database.Row) (*domain.Patient, error) {
	patient := &domain.Patient{}
	var middleName, email, phoneNumber *string
	var height, weight *float64
	var street, city, state, postalCode, country *string
	var dob time.Time

	err := row.Scan(
		&patient.ID,
		&patient.FirstName,
		&patient.LastName,
		&middleName,
		database.Time(&dob),
		&patient.Gender,
		&email,
		&phoneNumber,
		&height,
		&weight,
		&street,
		&city,
		&state,
		&postalCode,
		&country,
		database.Time(&patient.CreatedAt),
		database.Time(&patient.UpdatedAt),
	)
	if err != nil {
		return nil, err
	}

	patient.DateOfBirth = domain.Date(dob)
	if middleName != nil {
		patient.MiddleName = *middleName
	}
	if email != nil {
		patient.Email = *email
	}
	if phoneNumber != nil {
		patient.PhoneNumber = *phoneNumber
	}
	if height != nil {
		patient.Height = *height
	}
	if weight != nil {
		patient.Weight = *weight
	}
	if street != nil {
		patient.Address.Street = *street
	}
	if city != nil {
		patient.Address.City = *city
	}
	if state != nil {
		patient.Address.State = *state
	}
	if postalCode != nil {
		patient.Address.PostalCode = *postalCode
	}
	if country != nil {
		patient.Address.Country = *country
	}

	return patient, nil
}

func nullString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func nullFloat(value float64) *float64 {
	if value == 0 {
		return nil
	}
	return &value
}

var _ domain.PatientRepository = (*SQLRepository)(nil)
