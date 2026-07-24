package infrastructure

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

func TestSQLPatientRepository_CreateListSearch(t *testing.T) {
	ctx := context.Background()
	db, err := database.Connect(ctx, filepath.Join(t.TempDir(), "patients.db"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	dob := time.Date(1988, 5, 4, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Jamie", "Rivera", dob, domain.GenderFemale)
	patient.Email = "jamie@example.com"
	patient.PhoneNumber = "555-0100"
	patient.Height = 170
	patient.Weight = 65
	patient.Address = domain.Address{
		Street:     "1 Main",
		City:       "Town",
		State:      "TX",
		PostalCode: "75001",
		Country:    "US",
	}

	if err := repo.Create(ctx, patient); err != nil {
		t.Fatalf("create (17 placeholders): %v", err)
	}

	got, err := repo.GetByID(ctx, patient.ID.String())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.FirstName != "Jamie" || got.Email != "jamie@example.com" {
		t.Fatalf("unexpected patient: %+v", got)
	}

	list, total, err := repo.ListPaginated(ctx, 1, 20, domain.PatientListFilter{Search: "jamie"})
	if err != nil {
		t.Fatalf("list search (reused $1): %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("expected 1 match, total=%d len=%d", total, len(list))
	}

	got.MiddleName = "A"
	got.UpdatedAt = time.Now().UTC()
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update ($10+ placeholders): %v", err)
	}

	updated, err := repo.GetByID(ctx, got.ID.String())
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if updated.MiddleName != "A" {
		t.Fatalf("middle name not saved: %q", updated.MiddleName)
	}
}
