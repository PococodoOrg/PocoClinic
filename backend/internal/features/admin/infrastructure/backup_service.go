package infrastructure

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	admindomain "github.com/dksch/pococlinic/internal/features/admin/domain"
	"github.com/dksch/pococlinic/internal/pkg/backup"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

// BackupService runs backup operations against the configured database.
type BackupService struct {
	pool          *database.DB
	dir           string
	documentsDir  string
	appVersion    string
	mu            sync.Mutex
}

func NewBackupService(pool *database.DB, dir, documentsDir, appVersion string) *BackupService {
	return &BackupService{
		pool:         pool,
		dir:          dir,
		documentsDir: documentsDir,
		appVersion:   appVersion,
	}
}

func (s *BackupService) CreateBackup(ctx context.Context) (*admindomain.BackupResult, error) {
	if s.pool == nil {
		return nil, fmt.Errorf("backup requires persistent database storage")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := backup.Create(ctx, s.pool, s.dir, s.documentsDir, s.appVersion)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC()
	if info, err := backup.Open(path); err == nil {
		createdAt = info.Manifest.CreatedAt
	}

	return &admindomain.BackupResult{
		Filename:  filepath.Base(path),
		Path:      path,
		CreatedAt: createdAt,
	}, nil
}

func (s *BackupService) ListBackups(ctx context.Context) ([]admindomain.BackupEntry, error) {
	_ = ctx
	items, err := backup.List(s.dir)
	if err != nil {
		return nil, err
	}

	result := make([]admindomain.BackupEntry, 0, len(items))
	for _, item := range items {
		result = append(result, admindomain.BackupEntry{
			Filename:   item.Filename,
			Path:       item.Path,
			CreatedAt:  item.CreatedAt,
			AppVersion: item.Manifest.AppVersion,
		})
	}
	return result, nil
}

func (s *BackupService) RestoreBackup(ctx context.Context, filename string) error {
	if s.pool == nil {
		return fmt.Errorf("restore requires persistent database storage")
	}
	if err := validateBackupFilename(filename); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dir, filename)
	return backup.Restore(ctx, s.pool, path, s.documentsDir)
}

func (s *BackupService) VerifyBackup(ctx context.Context, filename string) (*admindomain.BackupVerifyResult, error) {
	_ = ctx
	if err := validateBackupFilename(filename); err != nil {
		return nil, err
	}
	path := filepath.Join(s.dir, filename)
	archive, err := backup.Open(path)
	if err != nil {
		return &admindomain.BackupVerifyResult{
			Filename:  filename,
			Valid:     false,
			CheckedAt: time.Now().UTC(),
			Message:   "backup archive could not be opened",
		}, nil
	}
	if err := archive.Verify(); err != nil {
		return &admindomain.BackupVerifyResult{
			Filename:  filename,
			Valid:     false,
			CheckedAt: time.Now().UTC(),
			Message:   "backup archive failed validation",
		}, nil
	}

	integrity := archive.VerifyIntegrity()
	result := &admindomain.BackupVerifyResult{
		Filename:  filename,
		Valid:     integrity.Valid,
		CheckedAt: time.Now().UTC(),
		Summary:   &integrity.Summary,
	}
	if integrity.Valid {
		result.Message = fmt.Sprintf(
			"Checksums match manifest (%d patients, %d documents)",
			integrity.Summary.TableRows["patients"],
			integrity.Summary.DocumentFiles,
		)
	} else {
		result.IntegrityIssues = integrity.Issues
		result.Message = "backup integrity check failed"
	}
	return result, nil
}

func validateBackupFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("backup filename is required")
	}
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid backup filename")
	}
	if !strings.HasPrefix(filename, "pococlinic-backup-") || !strings.HasSuffix(filename, ".tar.gz") {
		return fmt.Errorf("invalid backup filename")
	}
	return nil
}

var _ admindomain.BackupManager = (*BackupService)(nil)
