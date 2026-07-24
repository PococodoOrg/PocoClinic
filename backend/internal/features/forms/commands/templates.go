package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/forms/domain"
	"github.com/google/uuid"
)

type TemplateFieldCommand struct {
	ID        string   `json:"id" binding:"required"`
	Label     string   `json:"label" binding:"required"`
	Type      string   `json:"type" binding:"required"`
	Required  bool     `json:"required"`
	Options   []string `json:"options"`
	SortOrder int      `json:"sortOrder"`
}

type CreateTemplateCommand struct {
	GroupID  string                 `json:"groupId"`
	Name     string                 `json:"name" binding:"required"`
	FormType string                 `json:"formType" binding:"required"`
	Fields   []TemplateFieldCommand `json:"fields" binding:"required,min=1,dive"`
}

type CreateTemplateHandler interface {
	Handle(ctx context.Context, cmd CreateTemplateCommand) (*domain.FormTemplate, error)
}

type createTemplateHandler struct {
	repository domain.Repository
}

func NewCreateTemplateHandler(repo domain.Repository) CreateTemplateHandler {
	return &createTemplateHandler{repository: repo}
}

func (h *createTemplateHandler) Handle(ctx context.Context, cmd CreateTemplateCommand) (*domain.FormTemplate, error) {
	fields, err := toDomainFields(cmd.Fields)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidateFields(fields); err != nil {
		return nil, err
	}

	formType := domain.FormType(cmd.FormType)
	if err := domain.ValidateFormType(formType); err != nil {
		return nil, err
	}

	groupID, err := resolveGroupID(ctx, h.repository, cmd.GroupID)
	if err != nil {
		return nil, err
	}

	template := domain.NewFormTemplate(groupID, cmd.Name, formType, fields)
	if err := h.repository.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

type UpdateTemplateCommand struct {
	ID       string                 `json:"-"`
	GroupID  string                 `json:"groupId"`
	Name     string                 `json:"name" binding:"required"`
	FormType string                 `json:"formType" binding:"required"`
	Fields   []TemplateFieldCommand `json:"fields" binding:"required,min=1,dive"`
}

type UpdateTemplateHandler interface {
	Handle(ctx context.Context, cmd UpdateTemplateCommand) (*domain.FormTemplate, error)
}

type updateTemplateHandler struct {
	repository domain.Repository
}

func NewUpdateTemplateHandler(repo domain.Repository) UpdateTemplateHandler {
	return &updateTemplateHandler{repository: repo}
}

func (h *updateTemplateHandler) Handle(ctx context.Context, cmd UpdateTemplateCommand) (*domain.FormTemplate, error) {
	template, err := h.repository.GetTemplateByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	fields, err := toDomainFields(cmd.Fields)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidateFields(fields); err != nil {
		return nil, err
	}

	formType := domain.FormType(cmd.FormType)
	if err := domain.ValidateFormType(formType); err != nil {
		return nil, err
	}

	groupID, err := resolveGroupID(ctx, h.repository, cmd.GroupID)
	if err != nil {
		return nil, err
	}

	template.GroupID = groupID
	template.Name = cmd.Name
	template.FormType = formType
	template.Fields = fields
	template.UpdatedAt = time.Now()

	if err := h.repository.UpdateTemplate(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

type DeleteTemplateCommand struct {
	ID string
}

type DeleteTemplateHandler interface {
	Handle(ctx context.Context, cmd DeleteTemplateCommand) error
}

type deleteTemplateHandler struct {
	repository domain.Repository
}

func NewDeleteTemplateHandler(repo domain.Repository) DeleteTemplateHandler {
	return &deleteTemplateHandler{repository: repo}
}

func (h *deleteTemplateHandler) Handle(ctx context.Context, cmd DeleteTemplateCommand) error {
	return h.repository.DeleteTemplate(ctx, cmd.ID)
}

type SaveFormCommand struct {
	PatientID   string                 `json:"-"`
	TemplateID  string                 `json:"templateId" binding:"required"`
	EntryID     string                 `json:"entryId"`
	SubmittedBy string                 `json:"-"`
	Answers     map[string]interface{} `json:"answers" binding:"required"`
}

type SaveFormHandler interface {
	Handle(ctx context.Context, cmd SaveFormCommand) (*domain.FormSubmission, error)
}

type saveFormHandler struct {
	formRepository    domain.Repository
	patientExistsFunc func(ctx context.Context, patientID string) error
}

func NewSaveFormHandler(repo domain.Repository, patientExists func(ctx context.Context, patientID string) error) SaveFormHandler {
	return &saveFormHandler{
		formRepository:    repo,
		patientExistsFunc: patientExists,
	}
}

func (h *saveFormHandler) Handle(ctx context.Context, cmd SaveFormCommand) (*domain.FormSubmission, error) {
	if err := h.patientExistsFunc(ctx, cmd.PatientID); err != nil {
		return nil, err
	}

	template, err := h.formRepository.GetTemplateByID(ctx, cmd.TemplateID)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidateAnswers(template, cmd.Answers); err != nil {
		return nil, err
	}

	patientID, err := parseUUID(cmd.PatientID, "patient id")
	if err != nil {
		return nil, err
	}
	actorID, err := parseUUID(cmd.SubmittedBy, "submitted by")
	if err != nil {
		return nil, err
	}

	entryID := uuid.Nil
	version := 1

	if cmd.EntryID != "" {
		entryID, err = parseUUID(cmd.EntryID, "entry id")
		if err != nil {
			return nil, err
		}
		current, err := h.formRepository.GetCurrentSubmissionByEntryID(ctx, cmd.EntryID)
		if err != nil {
			return nil, err
		}
		if current.PatientID != patientID || current.TemplateID != template.ID {
			return nil, fmt.Errorf("entry does not match patient or template")
		}
		if err := h.formRepository.MarkEntryNotCurrent(ctx, cmd.EntryID); err != nil {
			return nil, err
		}
		version = current.Version + 1
	} else if template.FormType == domain.FormTypeSingleton {
		current, err := h.formRepository.GetCurrentSingletonSubmission(ctx, cmd.PatientID, cmd.TemplateID)
		if err != nil {
			return nil, err
		}
		if current != nil {
			if err := h.formRepository.MarkEntryNotCurrent(ctx, current.EntryID.String()); err != nil {
				return nil, err
			}
			entryID = current.EntryID
			version = current.Version + 1
		}
	}

	submission := domain.NewFormSubmission(entryID, version, template, patientID, actorID, cmd.Answers)
	if err := h.formRepository.CreateSubmission(ctx, submission); err != nil {
		return nil, err
	}
	return submission, nil
}

func resolveGroupID(ctx context.Context, repo domain.Repository, groupID string) (uuid.UUID, error) {
	if err := repo.EnsureDefaultGroup(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("form groups are not initialized: %w", err)
	}

	if groupID == "" {
		return domain.DefaultGroupID, nil
	}

	parsed, err := parseUUID(groupID, "group id")
	if err != nil {
		return uuid.Nil, err
	}
	if _, err := repo.GetGroupByID(ctx, groupID); err != nil {
		return uuid.Nil, err
	}
	return parsed, nil
}

func toDomainFields(fields []TemplateFieldCommand) ([]domain.FormField, error) {
	result := make([]domain.FormField, 0, len(fields))
	for _, field := range fields {
		result = append(result, domain.FormField{
			ID:        field.ID,
			Label:     field.Label,
			Type:      domain.FieldType(field.Type),
			Required:  field.Required,
			Options:   field.Options,
			SortOrder: field.SortOrder,
		})
	}
	return result, nil
}

func parseUUID(value, label string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", label)
	}
	return parsed, nil
}
