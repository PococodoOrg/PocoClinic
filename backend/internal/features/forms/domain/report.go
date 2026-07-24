package domain

import "time"

// FormSubmissionReport is a submission row enriched for cross-patient reporting.
type FormSubmissionReport struct {
	FormSubmission
	PatientName     string `json:"patientName"`
	SubmittedByName string `json:"submittedByName"`
}

// SubmissionReportPage is a paginated template submission report.
type SubmissionReportPage struct {
	Items    []*FormSubmissionReport `json:"items"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

// TemplateSubmissionFilter narrows a template-level report query.
type TemplateSubmissionFilter struct {
	Page     int
	PageSize int
	From     *time.Time
	To       *time.Time
}
