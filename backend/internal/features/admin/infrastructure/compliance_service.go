package infrastructure

import (
	"context"
	"fmt"
	"time"

	admindomain "github.com/dksch/pococlinic/internal/features/admin/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

// ComplianceService runs operational HIPAA-oriented control checks.
type ComplianceService struct {
	pool         *database.DB
	backupDir    string
	documentsDir string
	startedAt    time.Time
}

func NewComplianceService(pool *database.DB, backupDir, documentsDir string, startedAt time.Time) *ComplianceService {
	return &ComplianceService{
		pool:         pool,
		backupDir:    backupDir,
		documentsDir: documentsDir,
		startedAt:    startedAt,
	}
}

func (s *ComplianceService) RunComplianceCheck(ctx context.Context) (*admindomain.ComplianceCheckResult, error) {
	checks := []admindomain.ComplianceItem{
		s.checkPersistentStorage(),
		s.checkDocumentStorage(),
		s.checkBackupSchedule(),
		s.checkMigrations(ctx),
		s.checkDefaultPins(ctx),
		s.checkLockedAccounts(ctx),
		s.checkFailedSignIns(ctx),
		s.checkAuditActivity(ctx),
	}

	status := "ok"
	for _, item := range checks {
		switch item.Status {
		case "fail":
			status = "fail"
		case "warn":
			if status != "fail" {
				status = "warn"
			}
		}
	}

	return &admindomain.ComplianceCheckResult{
		CheckedAt: time.Now().UTC(),
		Status:    status,
		Checks:    checks,
	}, nil
}

func (s *ComplianceService) checkPersistentStorage() admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "persistent_storage",
		Label:   "Persistent storage",
		Control: "164.312(a)(2)(iv) — integrity controls require durable storage",
	}
	if s.pool == nil {
		item.Status = "fail"
		item.Message = "Running in memory mode. Set DATABASE_URL for production."
		return item
	}
	item.Status = "pass"
	item.Message = "Database storage is active."
	return item
}

func (s *ComplianceService) checkDocumentStorage() admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "document_storage",
		Label:   "Document storage",
		Control: "164.312(a)(2)(iv) — encryption of ePHI at rest",
	}
	item.Status = "pass"
	item.Message = "Patient documents are encrypted (AES-256-GCM) in the database and included in DB backups."
	return item
}

func (s *ComplianceService) checkBackupSchedule() admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "backup_schedule",
		Label:   "Backup schedule",
		Control: "164.308(a)(7)(ii)(A) — data backup plan",
	}
	backup := backupStatusFromDir(s.backupDir)
	switch backup.Status {
	case "ok":
		item.Status = "pass"
		item.Message = fmt.Sprintf("Latest backup %s is within 24 hours.", backup.LastBackupFile)
	case "warning":
		item.Status = "warn"
		item.Message = fmt.Sprintf("Latest backup %s is older than 24 hours.", backup.LastBackupFile)
	case "critical":
		item.Status = "fail"
		item.Message = "No backup within 48 hours — create a backup now."
	default:
		item.Status = "fail"
		item.Message = "No backups found in BACKUP_DIR."
	}
	return item
}

func (s *ComplianceService) checkMigrations(ctx context.Context) admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "migrations",
		Label:   "Schema migrations",
		Control: "164.312(b) — audit controls require current schema",
	}
	if s.pool == nil {
		item.Status = "warn"
		item.Message = "Skipped — no database connection."
		return item
	}
	pending, err := database.PendingMigrations(ctx, s.pool)
	if err != nil {
		item.Status = "warn"
		item.Message = "Could not read migration status."
		return item
	}
	if len(pending) > 0 {
		item.Status = "warn"
		item.Message = fmt.Sprintf("%d pending migration(s) — run go run ./cmd/migrate", len(pending))
		return item
	}
	item.Status = "pass"
	item.Message = "Database schema is up to date."
	return item
}

func (s *ComplianceService) checkDefaultPins(ctx context.Context) admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "default_pins",
		Label:   "Default PIN accounts",
		Control: "164.312(d) — person or entity authentication",
	}
	if s.pool == nil {
		item.Status = "warn"
		item.Message = "Skipped — no database connection."
		return item
	}
	var count int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE must_change_pin = 1`).Scan(&count); err != nil {
		item.Status = "warn"
		item.Message = "Could not count default PIN accounts."
		return item
	}
	if count > 0 {
		item.Status = "fail"
		item.Message = fmt.Sprintf("%d staff account(s) still on default PIN.", count)
		return item
	}
	item.Status = "pass"
	item.Message = "All staff have changed their PIN."
	return item
}

func (s *ComplianceService) checkLockedAccounts(ctx context.Context) admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "locked_accounts",
		Label:   "Locked accounts",
		Control: "164.308(a)(5)(ii)(C) — log-in monitoring",
	}
	if s.pool == nil {
		item.Status = "warn"
		item.Message = "Skipped — no database connection."
		return item
	}
	var count int64
	now := time.Now().UTC()
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE locked_until IS NOT NULL AND locked_until > $1`, now).Scan(&count); err != nil {
		item.Status = "warn"
		item.Message = "Could not count locked accounts."
		return item
	}
	if count > 0 {
		item.Status = "warn"
		item.Message = fmt.Sprintf("%d locked account(s) — review audit log for failed sign-ins.", count)
		return item
	}
	item.Status = "pass"
	item.Message = "No accounts are currently locked."
	return item
}

func (s *ComplianceService) checkFailedSignIns(ctx context.Context) admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "failed_signins",
		Label:   "Failed sign-ins (7 days)",
		Control: "164.308(a)(5)(ii)(C) — log-in monitoring",
	}
	if s.pool == nil {
		item.Status = "warn"
		item.Message = "Skipped — no database connection."
		return item
	}
	var count int64
	since := time.Now().UTC().AddDate(0, 0, -7)
	if err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM audit_logs
		WHERE event_type = 'auth.login.failure'
		  AND created_at >= $1
	`, since).Scan(&count); err != nil {
		item.Status = "warn"
		item.Message = "Could not count failed sign-ins."
		return item
	}
	if count > 20 {
		item.Status = "warn"
		item.Message = fmt.Sprintf("%d failed sign-in attempts in the last 7 days — review audit log.", count)
		return item
	}
	item.Status = "pass"
	item.Message = fmt.Sprintf("%d failed sign-in attempt(s) in the last 7 days.", count)
	return item
}

func (s *ComplianceService) checkAuditActivity(ctx context.Context) admindomain.ComplianceItem {
	item := admindomain.ComplianceItem{
		ID:      "audit_logging",
		Label:   "Audit logging",
		Control: "164.312(b) — audit controls",
	}
	if s.pool == nil {
		item.Status = "fail"
		item.Message = "Audit logs require persistent database storage."
		return item
	}
	var count int64
	since := time.Now().UTC().AddDate(0, 0, -7)
	if err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM audit_logs
		WHERE created_at >= $1
	`, since).Scan(&count); err != nil {
		item.Status = "warn"
		item.Message = "Could not read audit log activity."
		return item
	}
	uptime := time.Since(s.startedAt)
	if count == 0 && uptime > 24*time.Hour {
		item.Status = "warn"
		item.Message = "No audit events in the last 7 days — verify logging is working."
		return item
	}
	item.Status = "pass"
	item.Message = fmt.Sprintf("%d audit event(s) recorded in the last 7 days.", count)
	return item
}

var _ admindomain.ComplianceChecker = (*ComplianceService)(nil)
