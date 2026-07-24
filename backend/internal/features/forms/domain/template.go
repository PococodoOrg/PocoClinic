package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FormField defines one input on a template.
type FormField struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Type      FieldType `json:"type"`
	Required  bool      `json:"required"`
	Options   []string  `json:"options,omitempty"`
	SortOrder int       `json:"sortOrder"`
}

// FormTemplate is a reusable form definition.
type FormTemplate struct {
	ID        uuid.UUID   `json:"id"`
	GroupID   uuid.UUID   `json:"groupId"`
	Name      string      `json:"name"`
	FormType  FormType    `json:"formType"`
	Fields    []FormField `json:"fields"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// NewFormTemplate creates a template with generated metadata.
func NewFormTemplate(groupID uuid.UUID, name string, formType FormType, fields []FormField) *FormTemplate {
	now := time.Now()
	if groupID == uuid.Nil {
		groupID = DefaultGroupID
	}
	if formType == "" {
		formType = FormTypeSingleton
	}
	return &FormTemplate{
		ID:        uuid.New(),
		GroupID:   groupID,
		Name:      name,
		FormType:  formType,
		Fields:    fields,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ValidateFormType checks the template behavior type.
func ValidateFormType(formType FormType) error {
	switch formType {
	case FormTypeSingleton, FormTypeLog:
		return nil
	default:
		return fmt.Errorf("form type must be singleton or log")
	}
}

// ValidateFields checks template field definitions.
func ValidateFields(fields []FormField) error {
	if len(fields) == 0 {
		return fmt.Errorf("at least one field is required")
	}

	seen := make(map[string]struct{}, len(fields))
	for i, field := range fields {
		if field.ID == "" {
			return fmt.Errorf("field %d is missing an id", i+1)
		}
		if field.Label == "" {
			return fmt.Errorf("field %q is missing a label", field.ID)
		}
		if _, exists := seen[field.ID]; exists {
			return fmt.Errorf("duplicate field id %q", field.ID)
		}
		seen[field.ID] = struct{}{}

		switch field.Type {
		case FieldTypeText, FieldTypeTextarea, FieldTypeNumber, FieldTypeSelect, FieldTypeCheckbox:
		default:
			return fmt.Errorf("field %q has unsupported type %q", field.ID, field.Type)
		}
		if field.Type == FieldTypeSelect && len(field.Options) == 0 {
			return fmt.Errorf("field %q requires at least one option", field.ID)
		}
	}
	return nil
}
