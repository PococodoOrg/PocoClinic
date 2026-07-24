package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/dksch/pococlinic/internal/features/forms/domain"
)

type ListTemplateSubmissionsQuery struct {
	TemplateID string
	Page       int
	PageSize   int
	From       *string
	To         *string
}

type ListTemplateSubmissionsHandler interface {
	Handle(ctx context.Context, query ListTemplateSubmissionsQuery) (*domain.SubmissionReportPage, error)
}

type listTemplateSubmissionsHandler struct {
	repository domain.Repository
}

func NewListTemplateSubmissionsHandler(repo domain.Repository) ListTemplateSubmissionsHandler {
	return &listTemplateSubmissionsHandler{repository: repo}
}

func (h *listTemplateSubmissionsHandler) Handle(ctx context.Context, query ListTemplateSubmissionsQuery) (*domain.SubmissionReportPage, error) {
	if _, err := h.repository.GetTemplateByID(ctx, query.TemplateID); err != nil {
		return nil, err
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 25
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter := domain.TemplateSubmissionFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if query.From != nil && *query.From != "" {
		from, err := parseReportTime(*query.From)
		if err != nil {
			return nil, fmt.Errorf("invalid from date")
		}
		filter.From = &from
	}
	if query.To != nil && *query.To != "" {
		to, err := parseReportTime(*query.To)
		if err != nil {
			return nil, fmt.Errorf("invalid to date")
		}
		filter.To = &to
	}

	return h.repository.ListCurrentSubmissionsByTemplate(ctx, query.TemplateID, filter)
}

func parseReportTime(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid date")
}
