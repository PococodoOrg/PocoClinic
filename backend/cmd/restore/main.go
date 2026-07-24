package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/backup"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/pathsafe"
)

func main() {
	logger := logging.NewLogger()
	ctx := context.Background()

	file := flag.String("file", "", "backup file to restore (default: latest in BACKUP_DIR)")
	confirm := flag.Bool("confirm", false, "required flag to perform restore")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	if cfg.Database.URL == "" {
		logger.Error("DATABASE_URL is required for restore", fmt.Errorf("missing DATABASE_URL"))
		os.Exit(1)
	}
	if !*confirm {
		logger.Error("Restore aborted", fmt.Errorf("pass -confirm to replace all database contents"))
		os.Exit(1)
	}

	target := *file
	if target == "" {
		latest, err := backup.FindLatest(cfg.Backup.Dir)
		if err != nil {
			logger.Error("Failed to locate latest backup", err)
			os.Exit(1)
		}
		if latest == nil {
			logger.Error("No backup found", fmt.Errorf("no backups in %s", cfg.Backup.Dir))
			os.Exit(1)
		}
		target = latest.Path
	} else if filepath.IsAbs(target) {
		if err := pathsafe.ValidateBackupFilename(filepath.Base(target)); err != nil {
			logger.Error("Invalid backup file", err)
			os.Exit(1)
		}
	} else {
		if err := pathsafe.ValidateBackupFilename(target); err != nil {
			logger.Error("Invalid backup file", err)
			os.Exit(1)
		}
		resolved, err := pathsafe.JoinRoot(cfg.Backup.Dir, target)
		if err != nil {
			logger.Error("Invalid backup file path", err)
			os.Exit(1)
		}
		target = resolved
	}

	pool, err := database.OpenPool(ctx, cfg.Database.URL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("Restoring database from backup", "file", filepath.Base(target))
	if err := backup.Restore(ctx, pool, target, cfg.Storage.DocumentsDir); err != nil {
		logger.Error("Restore failed", err)
		os.Exit(1)
	}

	logger.Info("Restore completed successfully", "file", filepath.Base(target))
}
