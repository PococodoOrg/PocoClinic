package domain

import "time"

// StaffActivityRow is per-user activity for an admin report.
type StaffActivityRow struct {
	UserID       string     `json:"userId"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	SignIns      int64      `json:"signIns"`
	FailedSignIns int64     `json:"failedSignIns"`
	PatientViews int64      `json:"patientViews"`
	LastSignIn   *time.Time `json:"lastSignIn,omitempty"`
}

// StaffActivityReport summarizes staff usage from audit logs.
type StaffActivityReport struct {
	PeriodDays int                `json:"periodDays"`
	Staff      []StaffActivityRow `json:"staff"`
}
