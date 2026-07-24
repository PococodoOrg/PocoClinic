package commands

import (
	"context"
	"strings"
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/features/forms/domain"
	forminfra "github.com/PococodoOrg/PocoClinic/internal/features/forms/infrastructure"
	"github.com/google/uuid"
)

func validTemplateFields() []TemplateFieldCommand {
	return []TemplateFieldCommand{
		{ID: "q1", Label: "Question", Type: "text", Required: true},
	}
}

func TestCreateTemplateHandler_RejectsInvalidFields(t *testing.T) {
	repo := forminfra.NewMemoryRepository()
	handler := NewCreateTemplateHandler(repo)
	_, err := handler.Handle(context.Background(), CreateTemplateCommand{
		Name:     "Bad",
		FormType: string(domain.FormTypeSingleton),
		Fields:   []TemplateFieldCommand{{ID: "", Label: "Missing ID", Type: "text"}},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCreateTemplateHandler_CreatesTemplate(t *testing.T) {
	repo := forminfra.NewMemoryRepository()
	handler := NewCreateTemplateHandler(repo)
	template, err := handler.Handle(context.Background(), CreateTemplateCommand{
		Name:     "Intake",
		FormType: string(domain.FormTypeSingleton),
		Fields:   validTemplateFields(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if template.ID == uuid.Nil {
		t.Fatal("expected template id")
	}
	if template.GroupID != domain.DefaultGroupID {
		t.Fatalf("expected default group, got %v", template.GroupID)
	}
}

func TestSaveFormHandler_SingletonVersioning(t *testing.T) {
	ctx := context.Background()
	repo := forminfra.NewMemoryRepository()
	patientID := uuid.New()
	actorID := uuid.New()

	template := domain.NewFormTemplate(domain.DefaultGroupID, "Singleton", domain.FormTypeSingleton, []domain.FormField{
		{ID: "q1", Label: "Answer", Type: domain.FieldTypeText, Required: true},
	})
	if err := repo.CreateTemplate(ctx, template); err != nil {
		t.Fatalf("create template: %v", err)
	}

	patientExists := func(ctx context.Context, id string) error { return nil }
	handler := NewSaveFormHandler(repo, patientExists)

	first, err := handler.Handle(ctx, SaveFormCommand{
		PatientID:   patientID.String(),
		TemplateID:  template.ID.String(),
		SubmittedBy: actorID.String(),
		Answers:     map[string]interface{}{"q1": "first"},
	})
	if err != nil {
		t.Fatalf("first save: %v", err)
	}
	if first.Version != 1 {
		t.Fatalf("expected version 1, got %d", first.Version)
	}

	second, err := handler.Handle(ctx, SaveFormCommand{
		PatientID:   patientID.String(),
		TemplateID:  template.ID.String(),
		SubmittedBy: actorID.String(),
		Answers:     map[string]interface{}{"q1": "second"},
	})
	if err != nil {
		t.Fatalf("second save: %v", err)
	}
	if second.Version != 2 {
		t.Fatalf("expected version 2, got %d", second.Version)
	}
	if second.EntryID != first.EntryID {
		t.Fatalf("expected same entry id across singleton resubmit")
	}
}

func TestSaveFormHandler_EntryMismatchRejected(t *testing.T) {
	ctx := context.Background()
	repo := forminfra.NewMemoryRepository()
	patientID := uuid.New()
	otherPatient := uuid.New()
	actorID := uuid.New()

	template := domain.NewFormTemplate(domain.DefaultGroupID, "Log", domain.FormTypeLog, []domain.FormField{
		{ID: "q1", Label: "Answer", Type: domain.FieldTypeText, Required: true},
	})
	if err := repo.CreateTemplate(ctx, template); err != nil {
		t.Fatalf("create template: %v", err)
	}

	patientExists := func(ctx context.Context, id string) error { return nil }
	handler := NewSaveFormHandler(repo, patientExists)

	initial, err := handler.Handle(ctx, SaveFormCommand{
		PatientID:   patientID.String(),
		TemplateID:  template.ID.String(),
		SubmittedBy: actorID.String(),
		Answers:     map[string]interface{}{"q1": "first"},
	})
	if err != nil {
		t.Fatalf("initial save: %v", err)
	}

	_, err = handler.Handle(ctx, SaveFormCommand{
		PatientID:   otherPatient.String(),
		TemplateID:  template.ID.String(),
		EntryID:     initial.EntryID.String(),
		SubmittedBy: actorID.String(),
		Answers:     map[string]interface{}{"q1": "wrong patient"},
	})
	if err == nil || !strings.Contains(err.Error(), "entry does not match") {
		t.Fatalf("expected entry mismatch error, got %v", err)
	}
}
