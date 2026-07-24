package domain

// ActivitySummary is a rolling clinic activity report for administrators.
type ActivitySummary struct {
	PeriodDays        int   `json:"periodDays"`
	SignIns           int64 `json:"signIns"`
	FailedSignIns     int64 `json:"failedSignIns"`
	PatientsCreated   int64 `json:"patientsCreated"`
	ChartNotesAdded   int64 `json:"chartNotesAdded"`
	PatientViews      int64 `json:"patientViews"`
	DocumentsUploaded int64 `json:"documentsUploaded"`
	BackupsCreated    int64 `json:"backupsCreated"`
}
