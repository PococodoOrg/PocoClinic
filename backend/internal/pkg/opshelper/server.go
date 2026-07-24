package opshelper

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/backup"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/gin-gonic/gin"
)

// Status is the friendly dashboard snapshot for the local ops helper.
type Status struct {
	DatabaseConnected bool       `json:"databaseConnected"`
	StorageMode       string     `json:"storageMode"`
	LastBackupAt      *time.Time `json:"lastBackupAt,omitempty"`
	BackupAgeHours    *float64   `json:"backupAgeHours,omitempty"`
	BackupStatus      string     `json:"backupStatus"`
	LastBackupFile    string     `json:"lastBackupFile,omitempty"`
	BackupCount       int        `json:"backupCount"`
	PatientCount      int64      `json:"patientCount"`
	AppVersion        string     `json:"appVersion"`
	BackupDir         string     `json:"backupDir"`
	MainAppURL        string     `json:"mainAppUrl"`
	MainAppOnline     bool       `json:"mainAppOnline"`
}

// VerifyResult reports manifest checksum verification for the helper UI.
type VerifyResult struct {
	Filename        string                  `json:"filename"`
	Valid           bool                    `json:"valid"`
	CheckedAt       time.Time               `json:"checkedAt"`
	Message         string                  `json:"message,omitempty"`
	Summary         *backup.ManifestSummary `json:"summary,omitempty"`
	IntegrityIssues []string                `json:"integrityIssues,omitempty"`
}

// BackupEntry is a simplified backup row for the helper UI.
type BackupEntry struct {
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"createdAt"`
	AppVersion string    `json:"appVersion,omitempty"`
	FriendlyAt string    `json:"friendlyAt"`
}

// Server exposes localhost-only backup and restore helpers.
type Server struct {
	pool          *database.DB
	backupDir     string
	documentsDir  string
	appVersion    string
	mainAppURL      string
	healthCheckURL  string
	mu              sync.Mutex
}

func NewServer(pool *database.DB, backupDir, documentsDir, appVersion, mainAppURL, healthCheckURL string) *Server {
	return &Server{
		pool:           pool,
		backupDir:      backupDir,
		documentsDir:   documentsDir,
		appVersion:     appVersion,
		mainAppURL:     mainAppURL,
		healthCheckURL: healthCheckURL,
	}
}

func (s *Server) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/status", s.getStatus)
		api.GET("/backups", s.listBackups)
		api.POST("/backups", s.createBackup)
		api.POST("/backups/verify", s.verifyBackup)
		api.POST("/restore", s.restoreBackup)
	}
}

func (s *Server) getStatus(c *gin.Context) {
	status, err := s.buildStatus(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, status)
}

func (s *Server) listBackups(c *gin.Context) {
	if s.pool == nil {
		respondServiceUnavailable(c, "Database is not configured. Set DATABASE_URL and restart the helper.")
		return
	}

	items, err := backup.List(s.backupDir)
	if err != nil {
		respondInternalError(c)
		return
	}

	result := make([]BackupEntry, 0, len(items))
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		result = append(result, mapBackupEntry(item))
	}
	c.JSON(http.StatusOK, gin.H{"backups": result})
}

func (s *Server) createBackup(c *gin.Context) {
	if s.pool == nil {
		respondServiceUnavailable(c, "Database is not configured.")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := backup.Create(c.Request.Context(), s.pool, s.backupDir, s.documentsDir, s.appVersion)
	if err != nil {
		respondOperationError(c, err)
		return
	}

	info, _ := backup.Open(path)
	c.JSON(http.StatusCreated, gin.H{
		"filename":  filepath.Base(path),
		"path":      path,
		"createdAt": info.Manifest.CreatedAt,
		"message":   "Backup completed successfully. Copy the file to your USB drive when ready.",
	})
}

type verifyBackupRequest struct {
	Filename string `json:"filename" binding:"required"`
}

func (s *Server) verifyBackup(c *gin.Context) {
	var req verifyBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "Please choose a backup file to verify.")
		return
	}
	if err := validateBackupFilename(req.Filename); err != nil {
		respondBadRequest(c, "Invalid backup filename.")
		return
	}

	result, err := s.verifyBackupFile(req.Filename)
	if err != nil {
		respondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, result)
}

type restoreRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ConfirmText string `json:"confirmText" binding:"required"`
}

func (s *Server) restoreBackup(c *gin.Context) {
	if s.pool == nil {
		respondServiceUnavailable(c, "Database is not configured.")
		return
	}

	var req restoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "Please choose a backup file and type RESTORE to confirm.")
		return
	}
	if strings.TrimSpace(strings.ToUpper(req.ConfirmText)) != "RESTORE" {
		respondBadRequest(c, "Type RESTORE exactly to confirm.")
		return
	}
	if err := validateBackupFilename(req.Filename); err != nil {
		respondBadRequest(c, "Invalid backup filename.")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.backupDir, req.Filename)
	if err := backup.Restore(c.Request.Context(), s.pool, path, s.documentsDir); err != nil {
		respondOperationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Restore completed. Stop and restart the main PocoClinic service if it was still running during restore, " +
			"then ask staff to sign in again and spot-check a few patient records.",
		"filename": req.Filename,
	})
}

func (s *Server) buildStatus(ctx context.Context) (*Status, error) {
	status := &Status{
		StorageMode:  "memory",
		BackupStatus: "unknown",
		AppVersion:   s.appVersion,
		BackupDir:    s.backupDir,
		MainAppURL:   s.mainAppURL,
	}

	items, err := backup.List(s.backupDir)
	if err != nil {
		return nil, err
	}
	status.BackupCount = len(items)
	if len(items) > 0 {
		latest := items[len(items)-1]
		status.LastBackupFile = latest.Filename
		status.LastBackupAt = &latest.CreatedAt
		age := time.Since(latest.CreatedAt).Hours()
		status.BackupAgeHours = &age
		status.BackupStatus = backupAgeStatus(age)
	}

	status.MainAppOnline = s.checkMainAppOnline(ctx)

	if s.pool == nil {
		return status, nil
	}

	status.StorageMode = "database"
	if err := database.Ping(ctx, s.pool); err == nil {
		status.DatabaseConnected = true
		_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM patients`).Scan(&status.PatientCount)
	}

	return status, nil
}

func (s *Server) verifyBackupFile(filename string) (*VerifyResult, error) {
	path := filepath.Join(s.backupDir, filename)
	archive, err := backup.Open(path)
	if err != nil {
		return &VerifyResult{
			Filename:  filename,
			Valid:     false,
			CheckedAt: time.Now().UTC(),
			Message:   "This backup file could not be opened.",
		}, nil
	}
	if err := archive.Verify(); err != nil {
		return &VerifyResult{
			Filename:  filename,
			Valid:     false,
			CheckedAt: time.Now().UTC(),
			Message:   "Checksums in this backup do not match the file contents.",
		}, nil
	}

	integrity := archive.VerifyIntegrity()
	result := &VerifyResult{
		Filename:  filename,
		Valid:     integrity.Valid,
		CheckedAt: time.Now().UTC(),
		Summary:   &integrity.Summary,
	}
	if integrity.Valid {
		result.Message = fmt.Sprintf(
			"Backup verified (%d patients, %d document files).",
			integrity.Summary.TableRows["patients"],
			integrity.Summary.DocumentFiles,
		)
	} else {
		result.IntegrityIssues = integrity.Issues
		result.Message = "Backup integrity check failed."
	}
	return result, nil
}

func (s *Server) checkMainAppOnline(ctx context.Context) bool {
	if strings.TrimSpace(s.healthCheckURL) == "" {
		return false
	}
	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.healthCheckURL, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func mapBackupEntry(item backup.Info) BackupEntry {
	return BackupEntry{
		Filename:   item.Filename,
		CreatedAt:  item.CreatedAt,
		AppVersion: item.Manifest.AppVersion,
		FriendlyAt: item.CreatedAt.Local().Format("Mon Jan 2, 2006 at 3:04 PM"),
	}
}

func backupAgeStatus(ageHours float64) string {
	switch {
	case ageHours <= 24:
		return "ok"
	case ageHours <= 48:
		return "warning"
	default:
		return "critical"
	}
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

// LocalhostMiddleware rejects requests that did not originate from the local machine.
func LocalhostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip != "127.0.0.1" && ip != "::1" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "The backup helper is only available on this computer (localhost).",
			})
			return
		}
		c.Next()
	}
}
