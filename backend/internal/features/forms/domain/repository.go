package domain

import "context"

// Repository persists form groups, templates, and submissions.
type Repository interface {
	EnsureDefaultGroup(ctx context.Context) error

	CreateGroup(ctx context.Context, group *FormGroup) error
	UpdateGroup(ctx context.Context, group *FormGroup) error
	DeleteGroup(ctx context.Context, id string) error
	GetGroupByID(ctx context.Context, id string) (*FormGroup, error)
	ListGroups(ctx context.Context) ([]*FormGroup, error)

	CreateTemplate(ctx context.Context, template *FormTemplate) error
	UpdateTemplate(ctx context.Context, template *FormTemplate) error
	DeleteTemplate(ctx context.Context, id string) error
	GetTemplateByID(ctx context.Context, id string) (*FormTemplate, error)
	ListTemplates(ctx context.Context) ([]*FormTemplate, error)

	CreateSubmission(ctx context.Context, submission *FormSubmission) error
	MarkEntryNotCurrent(ctx context.Context, entryID string) error
	GetCurrentSubmissionByEntryID(ctx context.Context, entryID string) (*FormSubmission, error)
	GetCurrentSingletonSubmission(ctx context.Context, patientID, templateID string) (*FormSubmission, error)
	ListCurrentSubmissionsByPatient(ctx context.Context, patientID string) ([]*FormSubmission, error)
	ListCurrentSubmissionsByTemplate(ctx context.Context, templateID string, filter TemplateSubmissionFilter) (*SubmissionReportPage, error)
	ListSubmissionHistory(ctx context.Context, entryID string) ([]*FormSubmission, error)
}
