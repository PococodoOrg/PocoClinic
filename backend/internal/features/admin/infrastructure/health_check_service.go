package infrastructure

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	admindomain "github.com/dksch/pococlinic/internal/features/admin/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

// HealthCheckService runs live operational checks against the clinic server.
type HealthCheckService struct {
	pool         *database.DB
	backupDir    string
	documentsDir string
}

func NewHealthCheckService(pool *database.DB, backupDir, documentsDir string) *HealthCheckService {
	return &HealthCheckService{
		pool:         pool,
		backupDir:    backupDir,
		documentsDir: documentsDir,
	}
}

func (s *HealthCheckService) RunHealthCheck(ctx context.Context) (*admindomain.HealthCheckResult, error) {
	checks := []admindomain.HealthCheckItem{
		s.checkDatabase(ctx),
		s.checkBackupDir(),
		s.checkDocumentsStorage(),
		s.checkStorageUsage(ctx),
	}
	if s.pool != nil {
		checks = append(checks, s.checkMigrations(ctx), s.checkDocumentIntegrity(ctx))
	}

	status := "ok"
	for _, item := range checks {
		switch item.Status {
		case "critical":
			status = "critical"
		case "warning":
			if status != "critical" {
				status = "warning"
			}
		}
	}

	return &admindomain.HealthCheckResult{
		CheckedAt: time.Now().UTC(),
		Status:    status,
		Checks:    checks,
	}, nil
}

func (s *HealthCheckService) checkDatabase(ctx context.Context) admindomain.HealthCheckItem {
	item := admindomain.HealthCheckItem{
		ID:    "database",
		Label: "Database",
	}
	if s.pool == nil {
		item.Status = "critical"
		item.Message = "Running in memory mode — set DATABASE_URL for production."
		return item
	}
	if err := database.Ping(ctx, s.pool); err != nil {
		item.Status = "critical"
		item.Message = "Database is unreachable."
		return item
	}
	item.Status = "ok"
	item.Message = "Connected"
	return item
}

func (s *HealthCheckService) checkMigrations(ctx context.Context) admindomain.HealthCheckItem {
	item := admindomain.HealthCheckItem{
		ID:    "migrations",
		Label: "Migrations",
	}
	pending, err := database.PendingMigrations(ctx, s.pool)
	if err != nil {
		item.Status = "warning"
		item.Message = "Could not read migration status."
		return item
	}
	if len(pending) > 0 {
		item.Status = "warning"
		item.Message = fmt.Sprintf("%d pending migration(s) — run go run ./cmd/migrate", len(pending))
		return item
	}
	item.Status = "ok"
	item.Message = "Up to date"
	return item
}

func (s *HealthCheckService) checkBackupDir() admindomain.HealthCheckItem {
	item := admindomain.HealthCheckItem{
		ID:    "backup_dir",
		Label: "Backup directory",
	}
	if strings.TrimSpace(s.backupDir) == "" {
		item.Status = "warning"
		item.Message = "BACKUP_DIR is not configured."
		return item
	}
	if err := os.MkdirAll(s.backupDir, 0o755); err != nil {
		item.Status = "critical"
		item.Message = "Backup directory is not writable."
		return item
	}
	probe := filepath.Join(s.backupDir, ".write-probe")
	if err := os.WriteFile(probe, []byte("ok"), 0o640); err != nil {
		item.Status = "critical"
		item.Message = "Backup directory is not writable."
		return item
	}
	_ = os.Remove(probe)
	item.Status = "ok"
	item.Message = "Writable"
	return item
}

func (s *HealthCheckService) checkDocumentsStorage() admindomain.HealthCheckItem {
	item := admindomain.HealthCheckItem{
		ID:    "documents_storage",
		Label: "Document storage",
	}
	item.Status = "ok"
	item.Message = "Encrypted blobs in database (AES-256-GCM)"
	return item
}

func (s *HealthCheckService) checkDocumentIntegrity(ctx context.Context) admindomain.HealthCheckItem {
	item := admindomain.HealthCheckItem{
		ID:    "document_integrity",
		Label: "Document integrity",
	}

	rows, err := s.pool.Query(ctx, `
		SELECT COALESCE(storage_key, ''),
			CASE WHEN encrypted_content IS NULL OR length(encrypted_content) = 0 THEN false ELSE true END
		FROM patient_documents
	`)
	if err != nil {
		item.Status = "warning"
		item.Message = "Could not read patient_documents table."
		return item
	}
	defer rows.Close()

	type docRef struct {
		storageKey string
		hasBlob    bool
	}
	diskKeys := map[string]struct{}{}
	missingDisk := 0
	missingContent := 0
	blobCount := 0

	for rows.Next() {
		var ref docRef
		if err := rows.Scan(&ref.storageKey, &ref.hasBlob); err != nil {
			item.Status = "warning"
			item.Message = "Could not read document metadata."
			return item
		}
		if ref.hasBlob {
			blobCount++
			continue
		}
		if ref.storageKey == "" {
			missingContent++
			continue
		}
		diskKeys[ref.storageKey] = struct{}{}
		path := filepath.Join(s.documentsDir, filepath.FromSlash(ref.storageKey))
		if _, err := os.Stat(path); err != nil {
			missingDisk++
		}
	}
	if err := rows.Err(); err != nil {
		item.Status = "warning"
		item.Message = "Could not read document metadata."
		return item
	}

	orphans := 0
	if info, err := os.Stat(s.documentsDir); err == nil && info.IsDir() {
		_ = filepath.Walk(s.documentsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(s.documentsDir, path)
			if err != nil {
				return nil
			}
			key := filepath.ToSlash(rel)
			if _, ok := diskKeys[key]; !ok {
				orphans++
			}
			return nil
		})
	}

	switch {
	case missingContent > 0:
		item.Status = "critical"
		item.Message = fmt.Sprintf("%d document record(s) missing encrypted content and disk file", missingContent)
	case missingDisk > 0:
		item.Status = "critical"
		item.Message = fmt.Sprintf("%d legacy document record(s) missing file(s) on disk", missingDisk)
	case orphans > 0:
		item.Status = "warning"
		item.Message = fmt.Sprintf("%d orphan file(s) on disk without database records", orphans)
	default:
		item.Status = "ok"
		item.Message = fmt.Sprintf("%d encrypted in DB, %d legacy on disk", blobCount, len(diskKeys))
	}
	return item
}

func (s *HealthCheckService) checkStorageUsage(ctx context.Context) admindomain.HealthCheckItem {
	item := admindomain.HealthCheckItem{
		ID:    "storage_usage",
		Label: "Storage usage",
	}

	backupBytes := directorySizeBytes(s.backupDir)
	legacyDiskBytes := directorySizeBytes(s.documentsDir)
	var dbDocBytes int64
	if s.pool != nil {
		_ = s.pool.QueryRow(ctx, `
			SELECT COALESCE(SUM(size_bytes), 0) FROM patient_documents
		`).Scan(&dbDocBytes)
	}
	documentsBytes := dbDocBytes + legacyDiskBytes
	totalBytes := backupBytes + documentsBytes

	item.Message = fmt.Sprintf(
		"Backups %s, documents %s (%s total)",
		formatMegabytes(backupBytes),
		formatMegabytes(documentsBytes),
		formatMegabytes(totalBytes),
	)

	switch {
	case backupBytes > backupDirWarnBytes || documentsBytes > documentsDirWarnBytes:
		item.Status = "warning"
		item.Message += " — consider archiving old backups or reviewing document volume"
	case totalBytes == 0 && s.backupDir == "" && s.documentsDir == "":
		item.Status = "warning"
		item.Message = "Storage paths are not configured."
	default:
		item.Status = "ok"
	}
	return item
}

var _ admindomain.HealthChecker = (*HealthCheckService)(nil)
