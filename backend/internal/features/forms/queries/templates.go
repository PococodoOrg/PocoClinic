package queries

import (
	"context"
	"fmt"

	"github.com/dksch/pococlinic/internal/features/forms/domain"
)

type GetTemplateQuery struct {
	ID string
}

type GetTemplateHandler interface {
	Handle(ctx context.Context, query GetTemplateQuery) (*domain.FormTemplate, error)
}

type getTemplateHandler struct {
	repository domain.Repository
}

func NewGetTemplateHandler(repo domain.Repository) GetTemplateHandler {
	return &getTemplateHandler{repository: repo}
}

func (h *getTemplateHandler) Handle(ctx context.Context, query GetTemplateQuery) (*domain.FormTemplate, error) {
	return h.repository.GetTemplateByID(ctx, query.ID)
}

type ListTemplatesHandler interface {
	Handle(ctx context.Context) ([]*domain.FormTemplate, error)
}

type listTemplatesHandler struct {
	repository domain.Repository
}

func NewListTemplatesHandler(repo domain.Repository) ListTemplatesHandler {
	return &listTemplatesHandler{repository: repo}
}

func (h *listTemplatesHandler) Handle(ctx context.Context) ([]*domain.FormTemplate, error) {
	return h.repository.ListTemplates(ctx)
}

type ListGroupsHandler interface {
	Handle(ctx context.Context) ([]*domain.FormGroup, error)
}

type listGroupsHandler struct {
	repository domain.Repository
}

func NewListGroupsHandler(repo domain.Repository) ListGroupsHandler {
	return &listGroupsHandler{repository: repo}
}

func (h *listGroupsHandler) Handle(ctx context.Context) ([]*domain.FormGroup, error) {
	if err := h.repository.EnsureDefaultGroup(ctx); err != nil {
		return nil, err
	}
	return h.repository.ListGroups(ctx)
}

type GetGroupQuery struct {
	ID string
}

type GetGroupHandler interface {
	Handle(ctx context.Context, query GetGroupQuery) (*domain.FormGroup, error)
}

type getGroupHandler struct {
	repository domain.Repository
}

func NewGetGroupHandler(repo domain.Repository) GetGroupHandler {
	return &getGroupHandler{repository: repo}
}

func (h *getGroupHandler) Handle(ctx context.Context, query GetGroupQuery) (*domain.FormGroup, error) {
	return h.repository.GetGroupByID(ctx, query.ID)
}

type ListSubmissionsQuery struct {
	PatientID string
}

type ListSubmissionsHandler interface {
	Handle(ctx context.Context, query ListSubmissionsQuery) ([]*domain.FormSubmission, error)
}

type listSubmissionsHandler struct {
	repository        domain.Repository
	patientExistsFunc func(ctx context.Context, patientID string) error
}

func NewListSubmissionsHandler(repo domain.Repository, patientExists func(ctx context.Context, patientID string) error) ListSubmissionsHandler {
	return &listSubmissionsHandler{
		repository:        repo,
		patientExistsFunc: patientExists,
	}
}

func (h *listSubmissionsHandler) Handle(ctx context.Context, query ListSubmissionsQuery) ([]*domain.FormSubmission, error) {
	if err := h.patientExistsFunc(ctx, query.PatientID); err != nil {
		return nil, err
	}
	return h.repository.ListCurrentSubmissionsByPatient(ctx, query.PatientID)
}

type ListSubmissionHistoryQuery struct {
	PatientID string
	EntryID   string
}

type ListSubmissionHistoryHandler interface {
	Handle(ctx context.Context, query ListSubmissionHistoryQuery) ([]*domain.FormSubmission, error)
}

type listSubmissionHistoryHandler struct {
	repository        domain.Repository
	patientExistsFunc func(ctx context.Context, patientID string) error
}

func NewListSubmissionHistoryHandler(repo domain.Repository, patientExists func(ctx context.Context, patientID string) error) ListSubmissionHistoryHandler {
	return &listSubmissionHistoryHandler{
		repository:        repo,
		patientExistsFunc: patientExists,
	}
}

func (h *listSubmissionHistoryHandler) Handle(ctx context.Context, query ListSubmissionHistoryQuery) ([]*domain.FormSubmission, error) {
	if err := h.patientExistsFunc(ctx, query.PatientID); err != nil {
		return nil, err
	}

	current, err := h.repository.GetCurrentSubmissionByEntryID(ctx, query.EntryID)
	if err != nil {
		return nil, err
	}
	if current.PatientID.String() != query.PatientID {
		return nil, fmt.Errorf("form entry not found")
	}

	return h.repository.ListSubmissionHistory(ctx, query.EntryID)
}

func ErrNotFound(message string) error {
	return fmt.Errorf("%s", message)
}
