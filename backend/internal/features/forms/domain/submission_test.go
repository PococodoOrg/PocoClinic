package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func testTemplate() *FormTemplate {
	return &FormTemplate{
		ID:       uuid.New(),
		Name:     "Vitals",
		FormType: FormTypeSingleton,
		Fields: []FormField{
			{ID: "note", Label: "Notes", Type: FieldTypeText, Required: true},
			{ID: "score", Label: "Score", Type: FieldTypeNumber, Required: false},
			{ID: "ok", Label: "OK", Type: FieldTypeCheckbox, Required: true},
			{ID: "side", Label: "Side", Type: FieldTypeSelect, Required: true, Options: []string{"left", "right"}},
		},
	}
}

func TestValidateAnswers_RequiredMissing(t *testing.T) {
	err := ValidateAnswers(testTemplate(), map[string]interface{}{})
	if err == nil || !strings.Contains(err.Error(), "Notes") {
		t.Fatalf("expected required error, got %v", err)
	}
}

func TestValidateAnswers_TypeMismatch(t *testing.T) {
	answers := map[string]interface{}{
		"note": "fine",
		"score": "not-a-number",
		"ok": true,
		"side": "left",
	}
	err := ValidateAnswers(testTemplate(), answers)
	if err == nil || !strings.Contains(err.Error(), "must be a number") {
		t.Fatalf("expected number type error, got %v", err)
	}
}

func TestValidateAnswers_InvalidSelect(t *testing.T) {
	answers := map[string]interface{}{
		"note": "fine",
		"ok": true,
		"side": "both",
	}
	err := ValidateAnswers(testTemplate(), answers)
	if err == nil || !strings.Contains(err.Error(), "invalid selection") {
		t.Fatalf("expected invalid selection error, got %v", err)
	}
}

func TestValidateAnswers_Valid(t *testing.T) {
	answers := map[string]interface{}{
		"note": "Patient stable",
		"score": float64(7),
		"ok": true,
		"side": "left",
	}
	if err := ValidateAnswers(testTemplate(), answers); err != nil {
		t.Fatalf("expected valid answers, got %v", err)
	}
}

func TestNewFormSubmissionMetadata(t *testing.T) {
	template := testTemplate()
	patientID := uuid.New()
	actorID := uuid.New()
	submission := NewFormSubmission(uuid.Nil, 0, template, patientID, actorID, map[string]interface{}{"note": "x"})
	if submission.Version != 1 || !submission.IsCurrent {
		t.Fatalf("unexpected version/current: v=%d current=%v", submission.Version, submission.IsCurrent)
	}
	if submission.EntryID == uuid.Nil {
		t.Fatal("expected generated entry id")
	}
	if len(submission.FieldSnapshot) != len(template.Fields) {
		t.Fatalf("expected field snapshot, got %d fields", len(submission.FieldSnapshot))
	}
}
