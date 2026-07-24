package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func validField(id string) FormField {
	return FormField{ID: id, Label: "Label " + id, Type: FieldTypeText, Required: true}
}

func TestValidateFields_RejectsEmpty(t *testing.T) {
	err := ValidateFields(nil)
	if err == nil || !strings.Contains(err.Error(), "at least one field") {
		t.Fatalf("expected empty fields error, got %v", err)
	}
}

func TestValidateFields_RejectsDuplicateIDs(t *testing.T) {
	fields := []FormField{validField("q1"), validField("q1")}
	err := ValidateFields(fields)
	if err == nil || !strings.Contains(err.Error(), "duplicate field id") {
		t.Fatalf("expected duplicate id error, got %v", err)
	}
}

func TestValidateFields_RejectsSelectWithoutOptions(t *testing.T) {
	fields := []FormField{{ID: "s1", Label: "Choice", Type: FieldTypeSelect, Required: true}}
	err := ValidateFields(fields)
	if err == nil || !strings.Contains(err.Error(), "requires at least one option") {
		t.Fatalf("expected select options error, got %v", err)
	}
}

func TestValidateFields_AcceptsAllTypes(t *testing.T) {
	fields := []FormField{
		{ID: "t", Label: "Text", Type: FieldTypeText},
		{ID: "a", Label: "Area", Type: FieldTypeTextarea},
		{ID: "n", Label: "Num", Type: FieldTypeNumber},
		{ID: "s", Label: "Select", Type: FieldTypeSelect, Options: []string{"a"}},
		{ID: "c", Label: "Check", Type: FieldTypeCheckbox},
	}
	if err := ValidateFields(fields); err != nil {
		t.Fatalf("expected valid fields, got %v", err)
	}
}

func TestValidateFormType(t *testing.T) {
	if err := ValidateFormType(FormTypeSingleton); err != nil {
		t.Fatalf("singleton: %v", err)
	}
	if err := ValidateFormType(FormTypeLog); err != nil {
		t.Fatalf("log: %v", err)
	}
	if err := ValidateFormType("weekly"); err == nil {
		t.Fatal("expected invalid form type error")
	}
}

func TestNewFormTemplateDefaults(t *testing.T) {
	template := NewFormTemplate(uuid.Nil, "Intake", "", []FormField{validField("q1")})
	if template.GroupID != DefaultGroupID {
		t.Fatalf("expected default group, got %v", template.GroupID)
	}
	if template.FormType != FormTypeSingleton {
		t.Fatalf("expected singleton default, got %q", template.FormType)
	}
}
