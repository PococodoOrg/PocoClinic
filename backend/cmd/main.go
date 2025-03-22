package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/commands"
	"github.com/dksch/pococlinic/internal/features/patients/handlers"
	"github.com/dksch/pococlinic/internal/features/patients/infrastructure"
	"github.com/dksch/pococlinic/internal/features/patients/queries"
	"github.com/dksch/pococlinic/internal/pkg/config"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/dksch/pococlinic/internal/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	// Set up rate limiter
	rateLimiter := middleware.NewIPRateLimiter(
		rate.Limit(cfg.Security.RateLimit.RequestsPerSecond),
		cfg.Security.RateLimit.BurstSize,
	)
	rateLimiter.CleanupTask() // Start cleanup task

	// Initialize repositories and handlers
	patientRepo := infrastructure.NewMemoryRepository()
	createPatientHandler := commands.NewCreatePatientHandler(patientRepo)
	getPatientsHandler := queries.NewGetPatientsHandler(patientRepo)
	getPatientHandler := queries.NewGetPatientHandler(patientRepo)
	patientHandler := handlers.NewPatientHandler(createPatientHandler, getPatientsHandler, getPatientHandler, logger)

	// Initialize router with security middleware
	router := gin.New() // Don't use Default() as we'll add our own middleware
	router.Use(
		middleware.Recovery(),
		middleware.SecurityHeaders(),
		middleware.RateLimiterMiddleware(rateLimiter),
	)

	// Configure CORS
	corsConfig := cfg.ConfigureCORS()
	router.Use(cors.New(corsConfig))

	// Initialize routes
	initializeRoutes(router, patientHandler)

	// Configure server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server
	go func() {
		logger.Info("Starting server",
			"host", cfg.Server.Host,
			"port", cfg.Server.Port,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}

func initializeRoutes(router *gin.Engine, patientHandler *handlers.PatientHandler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

	v1 := router.Group("/api/v1")
	patientHandler.RegisterRoutes(v1)
}
