package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dksch/pococlinic/internal/pkg/config"
	"github.com/dksch/pococlinic/internal/pkg/database"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/dksch/pococlinic/internal/pkg/middleware"
	"github.com/dksch/pococlinic/internal/pkg/opshelper"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.NewLogger()
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	host := getenv("OPS_HELPER_HOST", "127.0.0.1")
	port := getenv("OPS_HELPER_PORT", "9090")
	mainAppURL := getenv("OPS_HELPER_MAIN_APP_URL", "http://localhost:3000")

	var pool *database.DB
	if cfg.Database.URL != "" {
		pool, err = database.OpenPool(ctx, cfg.Database.URL)
		if err != nil {
			logger.Error("Failed to connect to database", err)
			os.Exit(1)
		}
		defer pool.Close()
	} else {
		logger.Info("DATABASE_URL not set; backup and restore actions will be unavailable")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), middleware.SecurityHeaders(), opshelper.LocalhostMiddleware())

	opshelper.NewServer(pool, cfg.Backup.Dir, cfg.Storage.DocumentsDir, cfg.App.Version, mainAppURL).RegisterRoutes(router)

	router.GET("/", func(c *gin.Context) {
		if staticDir := resolveStaticDir(); staticDir != "" {
			c.File(staticDir + "/index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fallbackHTML(host, port)))
	})

	staticDir := resolveStaticDir()
	if staticDir != "" {
		router.Static("/assets", staticDir+"/assets")
		router.NoRoute(func(c *gin.Context) {
			if c.Request.Method != http.MethodGet {
				c.Status(http.StatusNotFound)
				return
			}
			c.File(staticDir + "/index.html")
		})
		logger.Info("Serving ops helper UI", "dir", staticDir)
	} else {
		logger.Info("Ops helper UI not built yet; run build-ops-helper.bat")
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("PocoClinic Backup Helper listening (localhost only)", "url", fmt.Sprintf("http://%s:%s", host, port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Ops helper failed to start", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func resolveStaticDir() string {
	candidates := []string{
		"static",
		"../../../ops-helper/dist",
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate + "/index.html"); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}

func fallbackHTML(host, port string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>PocoClinic Backup Helper</title></head>
<body style="font-family:Segoe UI,sans-serif;max-width:640px;margin:48px auto;padding:0 24px;">
<h1>PocoClinic Backup Helper</h1>
<p>The friendly UI is not built yet. From the project root run <code>build-ops-helper.bat</code>, then restart this helper.</p>
<p>API is available at <code>http://%s:%s/api/status</code></p>
</body></html>`, host, port)
}

var _ fs.FS
