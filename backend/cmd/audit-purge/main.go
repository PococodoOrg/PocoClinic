package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	"github.com/dksch/pococlinic/internal/pkg/auditretention"
	"github.com/dksch/pococlinic/internal/pkg/config"
	"github.com/dksch/pococlinic/internal/pkg/database"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/google/uuid"
)

func main() {
	logger := logging.NewLogger()
	ctx := context.Background()

	days := flag.Int("days", 0, "retention days (default: AUDIT_RETENTION_DAYS env)")
	dryRun := flag.Bool("dry-run", false, "report counts only, do not delete")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	retentionDays := *days
	if retentionDays <= 0 {
		retentionDays, _ = strconv.Atoi(os.Getenv("AUDIT_RETENTION_DAYS"))
	}
	if retentionDays <= 0 {
		logger.Error("Retention days required", fmt.Errorf("set -days or AUDIT_RETENTION_DAYS"))
		os.Exit(1)
	}

	if cfg.Database.URL == "" {
		logger.Error("DATABASE_URL is required", fmt.Errorf("missing DATABASE_URL"))
		os.Exit(1)
	}

	pool, err := database.OpenPool(ctx, cfg.Database.URL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer pool.Close()

	total, err := auditretention.CountRows(ctx, pool)
	if err != nil {
		logger.Error("Failed to count audit rows", err)
		os.Exit(1)
	}

	if *dryRun {
		logger.Info("Audit purge dry run", "totalRows", total, "retentionDays", retentionDays)
		return
	}

	result, err := auditretention.PurgeOlderThan(ctx, pool, retentionDays)
	if err != nil {
		logger.Error("Audit purge failed", err)
		os.Exit(1)
	}

	details, _ := json.Marshal(map[string]string{
		"deletedRows":   fmt.Sprintf("%d", result.DeletedRows),
		"retentionDays": fmt.Sprintf("%d", result.RetentionDays),
		"cutoff":        result.Cutoff.Format("2006-01-02"),
	})
	if _, err := pool.Exec(ctx, `
		INSERT INTO audit_logs (id, event_type, resource_type, success, details, created_at)
		VALUES ($1, $2, 'system', 1, $3, $4)
	`, uuid.New(), auditdomain.EventAuditPurged, details, time.Now().UTC()); err != nil {
		logger.Error("Failed to log audit purge event", err)
	}

	logger.Info("Audit purge completed",
		"deletedRows", result.DeletedRows,
		"remainingRows", total-result.DeletedRows,
		"retentionDays", result.RetentionDays,
		"cutoff", result.Cutoff.Format(time.RFC3339),
	)
}
