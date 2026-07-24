package infrastructure

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/dksch/pococlinic/internal/features/forms/domain"
)

type MemoryRepository struct {
	groups      map[string]*domain.FormGroup
	templates   map[string]*domain.FormTemplate
	submissions map[string]*domain.FormSubmission
	mu          sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		groups:      make(map[string]*domain.FormGroup),
		templates:   make(map[string]*domain.FormTemplate),
		submissions: make(map[string]*domain.FormSubmission),
	}
	repo.groups[domain.DefaultGroupID.String()] = &domain.FormGroup{
		ID:        domain.DefaultGroupID,
		Name:      "General",
		SortOrder: 0,
	}
	return repo
}

func (r *MemoryRepository) EnsureDefaultGroup(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.groups[domain.DefaultGroupID.String()]; !ok {
		r.groups[domain.DefaultGroupID.String()] = &domain.FormGroup{
			ID:        domain.DefaultGroupID,
			Name:      "General",
			SortOrder: 0,
		}
	}
	return nil
}

func (r *MemoryRepository) CreateGroup(ctx context.Context, group *domain.FormGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[group.ID.String()] = group
	return nil
}

func (r *MemoryRepository) UpdateGroup(ctx context.Context, group *domain.FormGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.groups[group.ID.String()]; !ok {
		return fmt.Errorf("form group not found")
	}
	r.groups[group.ID.String()] = group
	return nil
}

func (r *MemoryRepository) DeleteGroup(ctx context.Context, id string) error {
	if id == domain.DefaultGroupID.String() {
		return fmt.Errorf("cannot delete the default form group")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.groups[id]; !ok {
		return fmt.Errorf("form group not found")
	}
	for _, template := range r.templates {
		if template.GroupID.String() == id {
			template.GroupID = domain.DefaultGroupID
		}
	}
	delete(r.groups, id)
	return nil
}

func (r *MemoryRepository) GetGroupByID(ctx context.Context, id string) (*domain.FormGroup, error) {
	r.mu.RLock()
	defer r.mu.Unlock()
	group, ok := r.groups[id]
	if !ok {
		return nil, fmt.Errorf("form group not found")
	}
	return group, nil
}

func (r *MemoryRepository) ListGroups(ctx context.Context) ([]*domain.FormGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.FormGroup, 0, len(r.groups))
	for _, group := range r.groups {
		result = append(result, group)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SortOrder == result[j].SortOrder {
			return result[i].Name < result[j].Name
		}
		return result[i].SortOrder < result[j].SortOrder
	})
	return result, nil
}

func (r *MemoryRepository) CreateTemplate(ctx context.Context, template *domain.FormTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[template.ID.String()] = template
	return nil
}

func (r *MemoryRepository) UpdateTemplate(ctx context.Context, template *domain.FormTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[template.ID.String()]; !ok {
		return fmt.Errorf("form template not found")
	}
	r.templates[template.ID.String()] = template
	return nil
}

func (r *MemoryRepository) DeleteTemplate(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[id]; !ok {
		return fmt.Errorf("form template not found")
	}
	delete(r.templates, id)
	for submissionID, submission := range r.submissions {
		if submission.TemplateID.String() == id {
			delete(r.submissions, submissionID)
		}
	}
	return nil
}

func (r *MemoryRepository) GetTemplateByID(ctx context.Context, id string) (*domain.FormTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	template, ok := r.templates[id]
	if !ok {
		return nil, fmt.Errorf("form template not found")
	}
	return template, nil
}

func (r *MemoryRepository) ListTemplates(ctx context.Context) ([]*domain.FormTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.FormTemplate, 0, len(r.templates))
	for _, template := range r.templates {
		result = append(result, template)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func (r *MemoryRepository) CreateSubmission(ctx context.Context, submission *domain.FormSubmission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[submission.TemplateID.String()]; !ok {
		return fmt.Errorf("form template not found")
	}
	r.submissions[submission.ID.String()] = submission
	return nil
}

func (r *MemoryRepository) MarkEntryNotCurrent(ctx context.Context, entryID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, submission := range r.submissions {
		if submission.EntryID.String() == entryID && submission.IsCurrent {
			submission.IsCurrent = false
		}
	}
	return nil
}

func (r *MemoryRepository) GetCurrentSubmissionByEntryID(ctx context.Context, entryID string) (*domain.FormSubmission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, submission := range r.submissions {
		if submission.EntryID.String() == entryID && submission.IsCurrent {
			r.enrichSubmissionLocked(submission)
			return submission, nil
		}
	}
	return nil, fmt.Errorf("form entry not found")
}

func (r *MemoryRepository) GetCurrentSingletonSubmission(ctx context.Context, patientID, templateID string) (*domain.FormSubmission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, submission := range r.submissions {
		if submission.PatientID.String() == patientID &&
			submission.TemplateID.String() == templateID &&
			submission.IsCurrent {
			template, ok := r.templates[templateID]
			if ok && template.FormType == domain.FormTypeSingleton {
				r.enrichSubmissionLocked(submission)
				return submission, nil
			}
		}
	}
	return nil, nil
}

func (r *MemoryRepository) ListCurrentSubmissionsByPatient(ctx context.Context, patientID string) ([]*domain.FormSubmission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.FormSubmission, 0)
	for _, submission := range r.submissions {
		if submission.PatientID.String() == patientID && submission.IsCurrent {
			r.enrichSubmissionLocked(submission)
			result = append(result, submission)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})
	return result, nil
}

func (r *MemoryRepository) ListCurrentSubmissionsByTemplate(
	ctx context.Context,
	templateID string,
	filter domain.TemplateSubmissionFilter,
) (*domain.SubmissionReportPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matching := make([]*domain.FormSubmissionReport, 0)
	for _, submission := range r.submissions {
		if submission.TemplateID.String() != templateID || !submission.IsCurrent {
			continue
		}
		if filter.From != nil && submission.UpdatedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && submission.UpdatedAt.After(*filter.To) {
			continue
		}
		r.enrichSubmissionLocked(submission)
		matching = append(matching, &domain.FormSubmissionReport{
			FormSubmission: *submission,
		})
	}

	sort.Slice(matching, func(i, j int) bool {
		return matching[i].UpdatedAt.After(matching[j].UpdatedAt)
	})

	total := int64(len(matching))
	start := (filter.Page - 1) * filter.PageSize
	if start >= len(matching) {
		return &domain.SubmissionReportPage{
			Items:    []*domain.FormSubmissionReport{},
			Total:    total,
			Page:     filter.Page,
			PageSize: filter.PageSize,
		}, nil
	}
	end := start + filter.PageSize
	if end > len(matching) {
		end = len(matching)
	}

	return &domain.SubmissionReportPage{
		Items:    matching[start:end],
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (r *MemoryRepository) ListSubmissionHistory(ctx context.Context, entryID string) ([]*domain.FormSubmission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.FormSubmission, 0)
	for _, submission := range r.submissions {
		if submission.EntryID.String() == entryID {
			r.enrichSubmissionLocked(submission)
			result = append(result, submission)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Version > result[j].Version
	})
	return result, nil
}

func (r *MemoryRepository) enrichSubmissionLocked(submission *domain.FormSubmission) {
	if template, ok := r.templates[submission.TemplateID.String()]; ok {
		submission.TemplateName = template.Name
		submission.FormType = template.FormType
	}
}

var _ domain.Repository = (*MemoryRepository)(nil)
