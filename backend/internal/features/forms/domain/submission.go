package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FormSubmission stores answers for a patient.
type FormSubmission struct {
	ID           uuid.UUID              `json:"id"`
	EntryID      uuid.UUID              `json:"entryId"`
	Version      int                    `json:"version"`
	IsCurrent    bool                   `json:"isCurrent"`
	TemplateID   uuid.UUID              `json:"templateId"`
	TemplateName string                 `json:"templateName,omitempty"`
	FormType     FormType               `json:"formType,omitempty"`
	PatientID    uuid.UUID              `json:"patientId"`
	SubmittedBy  uuid.UUID              `json:"submittedBy"`
	UpdatedBy    uuid.UUID              `json:"updatedBy"`
	Answers      map[string]interface{} `json:"answers"`
	FieldSnapshot []FormField           `json:"fieldSnapshot,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// NewFormSubmission creates the first version of a form entry.
func NewFormSubmission(
	entryID uuid.UUID,
	version int,
	template *FormTemplate,
	patientID, actorID uuid.UUID,
	answers map[string]interface{},
) *FormSubmission {
	now := time.Now()
	if entryID == uuid.Nil {
		entryID = uuid.New()
	}
	if version < 1 {
		version = 1
	}

	return &FormSubmission{
		ID:            uuid.New(),
		EntryID:       entryID,
		Version:       version,
		IsCurrent:     true,
		TemplateID:    template.ID,
		TemplateName:  template.Name,
		FormType:      template.FormType,
		PatientID:     patientID,
		SubmittedBy:   actorID,
		UpdatedBy:     actorID,
		Answers:       answers,
		FieldSnapshot: cloneFields(template.Fields),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func cloneFields(fields []FormField) []FormField {
	cloned := make([]FormField, len(fields))
	copy(cloned, fields)
	return cloned
}

// ValidateAnswers checks submitted values against a template.
func ValidateAnswers(template *FormTemplate, answers map[string]interface{}) error {
	if answers == nil {
		answers = map[string]interface{}{}
	}

	for _, field := range template.Fields {
		value, exists := answers[field.ID]
		if !exists || value == nil || value == "" {
			if field.Required {
				return fmt.Errorf("field %q is required", field.Label)
			}
			continue
		}

		switch field.Type {
		case FieldTypeText, FieldTypeTextarea, FieldTypeSelect:
			if _, ok := value.(string); !ok {
				return fmt.Errorf("field %q must be text", field.Label)
			}
			if field.Type == FieldTypeSelect {
				selected := value.(string)
				valid := false
				for _, option := range field.Options {
					if option == selected {
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf("field %q has an invalid selection", field.Label)
				}
			}
		case FieldTypeNumber:
			switch value.(type) {
			case float64, int, int32, int64:
			default:
				return fmt.Errorf("field %q must be a number", field.Label)
			}
		case FieldTypeCheckbox:
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("field %q must be true or false", field.Label)
			}
		}
	}

	return nil
}
