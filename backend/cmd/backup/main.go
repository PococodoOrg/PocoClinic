package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/backup"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
)

func main() {
	logger := logging.NewLogger()
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	if cfg.Database.URL == "" {
		logger.Error("DATABASE_URL is required for backup", fmt.Errorf("missing DATABASE_URL"))
		os.Exit(1)
	}

	pool, err := database.OpenPool(ctx, cfg.Database.URL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer pool.Close()

	path, err := backup.Create(ctx, pool, cfg.Backup.Dir, cfg.Storage.DocumentsDir, cfg.App.Version)
	if err != nil {
		logger.Error("Backup failed", err)
		os.Exit(1)
	}

	logger.Info("Backup created successfully", "file", filepath.Base(path))
}
