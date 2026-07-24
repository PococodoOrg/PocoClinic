package main

import (
	"context"
	"fmt"
	"os"

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
		logger.Error("DATABASE_URL is required for migrations", fmt.Errorf("missing DATABASE_URL"))
		os.Exit(1)
	}

	pool, err := database.OpenPool(ctx, cfg.Database.URL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer pool.Close()

	pending, err := database.PendingMigrations(ctx, pool)
	if err != nil {
		logger.Error("Failed to list pending migrations", err)
		os.Exit(1)
	}

	if len(pending) == 0 {
		logger.Info("Database schema is up to date")
		return
	}

	logger.Info("Applying migrations", "count", len(pending), "files", pending)
	if err := database.RunMigrations(ctx, pool); err != nil {
		logger.Error("Migration failed", err)
		os.Exit(1)
	}

	logger.Info("Migrations applied successfully")
}
