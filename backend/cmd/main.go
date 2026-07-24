package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	adminhandlers "github.com/PococodoOrg/PocoClinic/internal/features/admin/handlers"
	admininfra "github.com/PococodoOrg/PocoClinic/internal/features/admin/infrastructure"
	admindomain "github.com/PococodoOrg/PocoClinic/internal/features/admin/domain"
	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	auditinfra "github.com/PococodoOrg/PocoClinic/internal/features/audit/infrastructure"
	auditqueries "github.com/PococodoOrg/PocoClinic/internal/features/audit/queries"
	auditservice "github.com/PococodoOrg/PocoClinic/internal/features/audit/service"
	authcommands "github.com/PococodoOrg/PocoClinic/internal/features/auth/commands"
	authdomain "github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	authhandlers "github.com/PococodoOrg/PocoClinic/internal/features/auth/handlers"
	authinfra "github.com/PococodoOrg/PocoClinic/internal/features/auth/infrastructure"
	authmiddleware "github.com/PococodoOrg/PocoClinic/internal/features/auth/middleware"
	authqueries "github.com/PococodoOrg/PocoClinic/internal/features/auth/queries"
	formdomain "github.com/PococodoOrg/PocoClinic/internal/features/forms/domain"
	formcommands "github.com/PococodoOrg/PocoClinic/internal/features/forms/commands"
	formhandlers "github.com/PococodoOrg/PocoClinic/internal/features/forms/handlers"
	forminfra "github.com/PococodoOrg/PocoClinic/internal/features/forms/infrastructure"
	formqueries "github.com/PococodoOrg/PocoClinic/internal/features/forms/queries"
	patientcommands "github.com/PococodoOrg/PocoClinic/internal/features/patients/commands"
	patientdomain "github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	patienthandlers "github.com/PococodoOrg/PocoClinic/internal/features/patients/handlers"
	patientinfra "github.com/PococodoOrg/PocoClinic/internal/features/patients/infrastructure"
	patientqueries "github.com/PococodoOrg/PocoClinic/internal/features/patients/queries"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/doccrypto"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/middleware"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/staticserve"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	logger := logging.NewLogger()
	ctx := context.Background()
	startedAt := time.Now()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	tokenConfig := cfg.TokenConfig()
	rateLimiter := middleware.NewIPRateLimiter(
		rate.Limit(cfg.Security.RateLimit.RequestsPerSecond),
		cfg.Security.RateLimit.BurstSize,
	)
	rateLimiter.CleanupTask()

	var pool *database.DB
	if cfg.Database.URL != "" {
		if cfg.Database.RunMigrations {
			pool, err = database.Connect(ctx, cfg.Database.URL)
		} else {
			pool, err = database.OpenPool(ctx, cfg.Database.URL)
		}
		if err != nil {
			logger.Error("Failed to connect to database", err)
			os.Exit(1)
		}
		defer pool.Close()
		if cfg.Database.RunMigrations {
			logger.Info("Connected to database (migrations applied on startup)")
		} else {
			logger.Info("Connected to database (migrations skipped; use cmd/migrate in deploy)")
		}
	} else {
		logger.Info("DATABASE_URL not set; using in-memory storage")
	}

	userRepo, sessionRepo, patientRepo, noteRepo, documentRepo, exerciseRepo, auditRepo := buildRepositories(pool)
	legacyFiles, err := buildLegacyFileStorage(cfg)
	if err != nil {
		logger.Error("Failed to initialize legacy document storage", err)
		os.Exit(1)
	}
	docCipher, err := doccrypto.NewCipher(cfg.Storage.DocumentEncryptionKey)
	if err != nil {
		logger.Error("Failed to initialize document encryption", err)
		os.Exit(1)
	}
	formRepo := buildFormRepository(pool)
	auditLogger := auditservice.NewAuditService(auditRepo, logger)

	dataDir := filepath.Dir(cfg.Storage.DocumentsDir)
	if key, err := authinfra.SeedDefaultAdmin(ctx, userRepo); err != nil {
		logger.Error("Failed to seed default admin user", err)
		os.Exit(1)
	} else if key != "" {
		credPath, err := authinfra.WriteBootstrapCredentialsFile(dataDir, authinfra.DefaultAdminEmail, key)
		if err != nil {
			logger.Error("Default admin created but credentials file could not be written", err)
			os.Exit(1)
		}
		logger.Info("Default admin user created; one-time credentials written to file (delete after setup)",
			"path", credPath,
		)
	}

	if authinfra.BootstrapCredentialsExist(dataDir) {
		logger.Info("SECURITY: bootstrap credentials file still present; delete after administrator setup",
			"path", authinfra.BootstrapCredentialsPath(dataDir),
		)
	}

	createUserHandler := authcommands.NewCreateUserHandler(userRepo)
	loginHandler := authcommands.NewLoginHandler(userRepo, sessionRepo, tokenConfig, auditLogger)
	refreshTokenHandler := authcommands.NewRefreshTokenHandler(userRepo, sessionRepo, tokenConfig, auditLogger)
	logoutHandler := authcommands.NewLogoutHandler(sessionRepo, auditLogger)
	reissueBadgeHandler := authcommands.NewReissueBadgeHandler(userRepo, sessionRepo, auditLogger)
	changePINHandler := authcommands.NewChangePINHandler(userRepo, sessionRepo, auditLogger)
	updateUserHandler := authcommands.NewUpdateUserHandler(userRepo, sessionRepo, auditLogger)
	unlockUserHandler := authcommands.NewUnlockUserHandler(userRepo, auditLogger)
	deleteUserHandler := authcommands.NewDeleteUserHandler(userRepo, sessionRepo, auditLogger)
	getUserHandler := authqueries.NewGetUserHandler(userRepo)
	listUsersHandler := authqueries.NewGetUsersHandler(userRepo)
	listAuditHandler := auditqueries.NewListAuditHandler(auditRepo)
	getUserAuditHandler := auditqueries.NewGetUserAuditHandler(auditRepo)
	authHandler := authhandlers.NewAuthHandler(
		logger,
		cfg.Auth,
		dataDir,
		createUserHandler,
		loginHandler,
		refreshTokenHandler,
		logoutHandler,
		reissueBadgeHandler,
		changePINHandler,
		unlockUserHandler,
		updateUserHandler,
		deleteUserHandler,
		getUserHandler,
		listUsersHandler,
		getUserAuditHandler,
		auditLogger,
	)
	authMiddleware := authmiddleware.NewAuthMiddleware(tokenConfig, cfg.Auth, userRepo, sessionRepo, logger)

	loginRateLimiter := middleware.NewIPRateLimiter(rate.Every(12*time.Second), 5)
	loginRateLimiter.CleanupTask()

	settingsRepo := buildSettingsRepository(pool)
	getFieldRequirementsHandler := patientqueries.NewGetFieldRequirementsHandler(settingsRepo)
	updateFieldRequirementsHandler := patientcommands.NewUpdateFieldRequirementsHandler(settingsRepo)

	createPatientHandler := patientcommands.NewCreatePatientHandler(patientRepo, settingsRepo)
	getPatientsHandler := patientqueries.NewGetPatientsHandler(patientRepo)
	getPatientHandler := patientqueries.NewGetPatientHandler(patientRepo)
	updatePatientHandler := patientcommands.NewUpdatePatientHandler(patientRepo, settingsRepo)
	deletePatientHandler := patientcommands.NewDeletePatientHandler(patientRepo)
	createNoteHandler := patientcommands.NewCreateNoteHandler(noteRepo, patientRepo)
	updateNoteHandler := patientcommands.NewUpdateNoteHandler(noteRepo)
	deleteNoteHandler := patientcommands.NewDeleteNoteHandler(noteRepo)
	listNotesHandler := patientqueries.NewListNotesHandler(noteRepo)
	uploadDocumentHandler := patientcommands.NewUploadDocumentHandler(documentRepo, patientRepo, docCipher)
	deleteDocumentHandler := patientcommands.NewDeleteDocumentHandler(documentRepo, legacyFiles)
	listDocumentsHandler := patientqueries.NewListDocumentsHandler(documentRepo)
	getDocumentHandler := patientqueries.NewGetDocumentHandler(documentRepo)
	openDocumentContentHandler := patientcommands.NewOpenDocumentContentHandler(documentRepo, docCipher, legacyFiles)
	patientHandler := patienthandlers.NewPatientHandler(
		createPatientHandler,
		getPatientsHandler,
		getPatientHandler,
		updatePatientHandler,
		deletePatientHandler,
		createNoteHandler,
		updateNoteHandler,
		deleteNoteHandler,
		listNotesHandler,
		uploadDocumentHandler,
		deleteDocumentHandler,
		listDocumentsHandler,
		getDocumentHandler,
		openDocumentContentHandler,
		getFieldRequirementsHandler,
		logger,
		auditLogger,
	)
	exerciseCommandHandler := patientcommands.NewExerciseLogCommandHandler(exerciseRepo, patientRepo)
	exerciseQueryHandler := patientqueries.NewExerciseLogQueryHandler(exerciseRepo)
	exerciseHandler := patienthandlers.NewExerciseLogHandler(
		exerciseCommandHandler,
		exerciseQueryHandler,
		logger,
		auditLogger,
	)

	patientExists := func(ctx context.Context, patientID string) error {
		if _, err := patientRepo.GetByID(ctx, patientID); err != nil {
			return fmt.Errorf("patient not found")
		}
		return nil
	}

	formHandler := formhandlers.NewFormHandler(
		logger,
		auditLogger,
		formcommands.NewCreateGroupHandler(formRepo),
		formcommands.NewUpdateGroupHandler(formRepo),
		formcommands.NewDeleteGroupHandler(formRepo),
		formcommands.NewCreateTemplateHandler(formRepo),
		formcommands.NewUpdateTemplateHandler(formRepo),
		formcommands.NewDeleteTemplateHandler(formRepo),
		formcommands.NewSaveFormHandler(formRepo, patientExists),
		formqueries.NewListGroupsHandler(formRepo),
		formqueries.NewGetGroupHandler(formRepo),
		formqueries.NewListTemplatesHandler(formRepo),
		formqueries.NewGetTemplateHandler(formRepo),
		formqueries.NewListSubmissionsHandler(formRepo, patientExists),
		formqueries.NewListSubmissionHistoryHandler(formRepo, patientExists),
		formqueries.NewListTemplateSubmissionsHandler(formRepo),
	)

	statsRepo := admininfra.NewStatsRepository(
		pool,
		patientRepo,
		userRepo,
		formRepo,
		cfg.Backup.Dir,
		cfg.Storage.DocumentsDir,
		startedAt,
		cfg.App.Env,
		cfg.App.Version,
	)
	var backupManager admindomain.BackupManager
	if pool != nil {
		backupManager = admininfra.NewBackupService(pool, cfg.Backup.Dir, cfg.Storage.DocumentsDir, cfg.App.Version)
	}
	healthChecker := admininfra.NewHealthCheckService(pool, cfg.Backup.Dir, cfg.Storage.DocumentsDir)
	complianceChecker := admininfra.NewComplianceService(pool, cfg.Backup.Dir, cfg.Storage.DocumentsDir, startedAt)
	performanceReader := admininfra.NewPerformanceService(pool)
	adminHandler := adminhandlers.NewAdminHandler(
		logger,
		statsRepo,
		backupManager,
		healthChecker,
		complianceChecker,
		performanceReader,
		listAuditHandler,
		auditRepo,
		auditLogger,
		getFieldRequirementsHandler,
		updateFieldRequirementsHandler,
	)

	router := gin.New()
	if len(cfg.Security.TrustedProxies) > 0 {
		if err := router.SetTrustedProxies(cfg.Security.TrustedProxies); err != nil {
			logger.Error("Failed to configure trusted proxies", err)
			os.Exit(1)
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}
	router.Use(
		middleware.Recovery(logger),
		middleware.SecurityHeaders(),
		middleware.RateLimiterMiddleware(rateLimiter),
	)
	router.Use(cors.New(cfg.ConfigureCORS()))
	initializeRoutes(router, authHandler, authMiddleware, loginRateLimiter, patientHandler, exerciseHandler, formHandler, adminHandler, cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Info("Starting server", "host", cfg.Server.Host, "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}

func buildFormRepository(pool *database.DB) formdomain.Repository {
	if pool != nil {
		return forminfra.NewSQLRepository(pool)
	}
	return forminfra.NewMemoryRepository()
}

func buildSettingsRepository(pool *database.DB) patientdomain.SettingsRepository {
	if pool != nil {
		return patientinfra.NewSQLSettingsRepository(pool)
	}
	return patientinfra.NewMemorySettingsRepository()
}

func buildRepositories(pool *database.DB) (
	authdomain.UserRepository,
	authdomain.SessionRepository,
	patientdomain.PatientRepository,
	patientdomain.NoteRepository,
	patientdomain.DocumentRepository,
	patientdomain.ExerciseLogRepository,
	auditdomain.Repository,
) {
	if pool != nil {
		return authinfra.NewSQLUserRepository(pool),
			authinfra.NewSQLSessionRepository(pool),
			patientinfra.NewSQLRepository(pool),
			patientinfra.NewSQLNoteRepository(pool),
			patientinfra.NewSQLDocumentRepository(pool),
			patientinfra.NewSQLExerciseLogRepository(pool),
			auditinfra.NewSQLRepository(pool)
	}

	patientRepo := patientinfra.NewMemoryRepository()
	return authinfra.NewMemoryUserRepository(),
		authinfra.NewMemorySessionRepository(),
		patientRepo,
		patientinfra.NewMemoryNoteRepository(patientRepo),
		patientinfra.NewMemoryDocumentRepository(patientRepo),
		patientinfra.NewMemoryExerciseLogRepository(patientRepo),
		auditinfra.NewMemoryRepository()
}

func buildLegacyFileStorage(cfg *config.Config) (patientcommands.LegacyFileStore, error) {
	if cfg.Database.URL != "" {
		return patientinfra.NewLocalFileStorage(cfg.Storage.DocumentsDir)
	}
	return patientinfra.NewMemoryFileStorage(), nil
}

func initializeRoutes(
	router *gin.Engine,
	authHandler *authhandlers.AuthHandler,
	authMiddleware *authmiddleware.AuthMiddleware,
	loginRateLimiter *middleware.IPRateLimiter,
	patientHandler *patienthandlers.PatientHandler,
	exerciseHandler *patienthandlers.ExerciseLogHandler,
	formHandler *formhandlers.FormHandler,
	adminHandler *adminhandlers.AdminHandler,
	cfg *config.Config,
) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"time":    time.Now().UTC(),
			"version": cfg.App.Version,
		})
	})

	v1 := router.Group("/api/v1")
	authHandler.RegisterRoutes(v1, authMiddleware, middleware.RateLimiterMiddleware(loginRateLimiter))
	patientHandler.RegisterRoutes(v1, authMiddleware)
	exerciseHandler.RegisterRoutes(v1, authMiddleware)
	formHandler.RegisterRoutes(v1, authMiddleware)
	adminHandler.RegisterRoutes(v1, authMiddleware)

	if cfg.App.Env == "production" {
		staticDir := staticserve.ResolveDir(
			os.Getenv("STATIC_DIR"),
			"static",
			filepath.Join(".", "static"),
		)
		if staticDir != "" {
			staticserve.RegisterSPA(router, staticDir)
		}
	}
}
