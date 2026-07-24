package domain

import (
	"context"
	"time"
)

// BackupStatus summarizes the most recent backup on disk.
type BackupStatus struct {
	LastBackupAt   *time.Time `json:"lastBackupAt,omitempty"`
	BackupAgeHours *float64   `json:"backupAgeHours,omitempty"`
	Status         string     `json:"status"`
	LastBackupFile string     `json:"lastBackupFile,omitempty"`
}

// ResourceCounts holds high-level entity totals.
type ResourceCounts struct {
	Patients         int64 `json:"patients"`
	Users            int64 `json:"users"`
	FormTemplates    int64 `json:"formTemplates"`
	FormSubmissions  int64 `json:"formSubmissions"`
	ActiveSessions   int64 `json:"activeSessions"`
	LockedAccounts   int64 `json:"lockedAccounts"`
	DefaultPinAccounts int64 `json:"defaultPinAccounts"`
	PatientDocuments int64 `json:"patientDocuments"`
}

// SecurityOverview summarizes auth and document storage health.
type SecurityOverview struct {
	DocumentsStorageReady bool `json:"documentsStorageReady"`
}

// SystemStatus is the admin dashboard health snapshot.
type SystemStatus struct {
	StorageMode         string   `json:"storageMode"`
	DatabaseConnected   bool     `json:"databaseConnected"`
	PendingMigrations   []string `json:"pendingMigrations"`
	LatestMigration     string   `json:"latestMigration,omitempty"`
	UptimeSeconds       int64    `json:"uptimeSeconds"`
	Environment         string   `json:"environment"`
	AppVersion          string   `json:"appVersion"`
	Backup              BackupStatus `json:"backup"`
	Counts              ResourceCounts `json:"counts"`
	Security            SecurityOverview `json:"security"`
}

// GenderCount is a patient count grouped by gender.
type GenderCount struct {
	Gender string `json:"gender"`
	Count  int64  `json:"count"`
}

// PatientCensus is a basic patient population report.
type PatientCensus struct {
	TotalPatients   int64         `json:"totalPatients"`
	ByGender        []GenderCount `json:"byGender"`
	AddedLast30Days int64         `json:"addedLast30Days"`
}

// StatsReader loads operational metrics for admin views.
type StatsReader interface {
	GetSystemStatus(ctx context.Context) (*SystemStatus, error)
	GetPatientCensus(ctx context.Context) (*PatientCensus, error)
	GetActivitySummary(ctx context.Context, periodDays int) (*ActivitySummary, error)
	GetStaffActivity(ctx context.Context, periodDays int) (*StaffActivityReport, error)
}
