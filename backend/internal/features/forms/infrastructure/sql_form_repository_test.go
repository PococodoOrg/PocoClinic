package infrastructure

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	authdomain "github.com/dksch/pococlinic/internal/features/auth/domain"
	authinfra "github.com/dksch/pococlinic/internal/features/auth/infrastructure"
	"github.com/dksch/pococlinic/internal/features/forms/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
	"github.com/google/uuid"
)

func TestSQLFormRepository_SubmissionHighPlaceholders(t *testing.T) {
	ctx := context.Background()
	db, err := database.Connect(ctx, filepath.Join(t.TempDir(), "forms.db"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	userRepo := authinfra.NewSQLUserRepository(db)
	key, keyCred, err := authdomain.GenerateKey()
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	pinCred, err := authdomain.NewCredential("1234")
	if err != nil {
		t.Fatalf("pin: %v", err)
	}
	user := authdomain.NewUser("doc@clinic.test", "Doc", authdomain.RoleDoctor)
	user.SetKeyCredential(keyCred, key)
	user.SetPINCredential(pinCred)
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	_, err = db.Exec(ctx, `
		INSERT INTO patients (id, first_name, last_name, date_of_birth, gender, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, patientID, "A", "B", "1990-01-01", "female", time.Now().UTC(), time.Now().UTC())
	if err != nil {
		t.Fatalf("patient: %v", err)
	}

	formRepo := NewSQLRepository(db)
	if err := formRepo.EnsureDefaultGroup(ctx); err != nil {
		t.Fatalf("default group: %v", err)
	}

	template := &domain.FormTemplate{
		ID:        uuid.New(),
		Name:      "Intake",
		FormType:  domain.FormTypeSingleton,
		GroupID:   domain.DefaultGroupID,
		Fields:    []domain.FormField{{ID: "q1", Label: "Hello", Type: domain.FieldTypeText}},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := formRepo.CreateTemplate(ctx, template); err != nil {
		t.Fatalf("template: %v", err)
	}

	now := time.Now().UTC()
	submission := &domain.FormSubmission{
		ID:            uuid.New(),
		EntryID:       uuid.New(),
		Version:       1,
		IsCurrent:     true,
		TemplateID:    template.ID,
		TemplateName:  template.Name,
		FormType:      template.FormType,
		PatientID:     patientID,
		SubmittedBy:   user.ID,
		UpdatedBy:     user.ID,
		Answers:       map[string]interface{}{"q1": "hi"},
		FieldSnapshot: template.Fields,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := formRepo.CreateSubmission(ctx, submission); err != nil {
		t.Fatalf("create submission ($10+ placeholders): %v", err)
	}

	got, err := formRepo.GetCurrentSubmissionByEntryID(ctx, submission.EntryID.String())
	if err != nil {
		t.Fatalf("get submission: %v", err)
	}
	if got == nil || !got.IsCurrent {
		t.Fatalf("expected current submission, got %#v", got)
	}
	if got.Answers["q1"] != "hi" {
		t.Fatalf("answers: %#v", got.Answers)
	}
}
