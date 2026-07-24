package domain

import (
	"context"
	"time"

	"github.com/dksch/pococlinic/internal/pkg/backup"
)

// BackupEntry describes a backup bundle on disk.
type BackupEntry struct {
	Filename   string    `json:"filename"`
	Path       string    `json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
	AppVersion string    `json:"appVersion,omitempty"`
}

// BackupResult is returned after creating a backup.
type BackupResult struct {
	Filename  string    `json:"filename"`
	Path      string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

// BackupVerifyResult reports manifest checksum verification.
type BackupVerifyResult struct {
	Filename        string                  `json:"filename"`
	Valid           bool                    `json:"valid"`
	CheckedAt       time.Time               `json:"checkedAt"`
	Message         string                  `json:"message,omitempty"`
	Summary         *backup.ManifestSummary `json:"summary,omitempty"`
	IntegrityIssues []string                `json:"integrityIssues,omitempty"`
}

// BackupManager creates, lists, and restores database backups.
type BackupManager interface {
	CreateBackup(ctx context.Context) (*BackupResult, error)
	ListBackups(ctx context.Context) ([]BackupEntry, error)
	RestoreBackup(ctx context.Context, filename string) error
	VerifyBackup(ctx context.Context, filename string) (*BackupVerifyResult, error)
}
