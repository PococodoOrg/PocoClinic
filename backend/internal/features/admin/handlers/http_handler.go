package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	auditqueries "github.com/dksch/pococlinic/internal/features/audit/queries"
	admindomain "github.com/dksch/pococlinic/internal/features/admin/domain"
	patientcommands "github.com/dksch/pococlinic/internal/features/patients/commands"
	patientdomain "github.com/dksch/pococlinic/internal/features/patients/domain"
	patientqueries "github.com/dksch/pococlinic/internal/features/patients/queries"
	authdomain "github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/features/auth/middleware"
	"github.com/dksch/pococlinic/internal/pkg/httperr"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type restoreBackupRequest struct {
	Filename string `json:"filename" binding:"required"`
	Confirm  bool   `json:"confirm"`
}

// AdminHandler serves admin-only operational endpoints.
type AdminHandler struct {
	logger            *logging.Logger
	statsReader       admindomain.StatsReader
	backupManager     admindomain.BackupManager
	healthChecker     admindomain.HealthChecker
	complianceChecker admindomain.ComplianceChecker
	performanceReader admindomain.PerformanceReader
	listAuditHandler  auditqueries.ListAuditHandler
	auditRepository   auditdomain.Repository
	auditLogger       auditdomain.Logger
	getFieldRequirementsHandler   patientqueries.GetFieldRequirementsHandler
	updateFieldRequirementsHandler patientcommands.UpdateFieldRequirementsHandler
}

func NewAdminHandler(
	logger *logging.Logger,
	statsReader admindomain.StatsReader,
	backupManager admindomain.BackupManager,
	healthChecker admindomain.HealthChecker,
	complianceChecker admindomain.ComplianceChecker,
	performanceReader admindomain.PerformanceReader,
	listAuditHandler auditqueries.ListAuditHandler,
	auditRepository auditdomain.Repository,
	auditLogger auditdomain.Logger,
	getFieldRequirementsHandler patientqueries.GetFieldRequirementsHandler,
	updateFieldRequirementsHandler patientcommands.UpdateFieldRequirementsHandler,
) *AdminHandler {
	return &AdminHandler{
		logger:            logger,
		statsReader:       statsReader,
		backupManager:     backupManager,
		healthChecker:     healthChecker,
		complianceChecker: complianceChecker,
		performanceReader: performanceReader,
		listAuditHandler:  listAuditHandler,
		auditRepository:   auditRepository,
		auditLogger:       auditLogger,
		getFieldRequirementsHandler:   getFieldRequirementsHandler,
		updateFieldRequirementsHandler: updateFieldRequirementsHandler,
	}
}

func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	admin := router.Group("/admin")
	admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(authdomain.RoleAdmin))
	{
		admin.GET("/system-status", h.GetSystemStatus)
		admin.GET("/health-check", h.RunHealthCheck)
		admin.GET("/compliance-check", h.RunComplianceCheck)
		admin.GET("/performance", h.GetPerformance)
		admin.GET("/reports/patient-census", h.GetPatientCensus)
		admin.GET("/reports/activity-summary", h.GetActivitySummary)
		admin.GET("/reports/staff-activity", h.GetStaffActivity)
		admin.GET("/backups", h.ListBackups)
		admin.POST("/backups", h.CreateBackup)
		admin.POST("/backups/verify", h.VerifyBackup)
		admin.POST("/restore", h.RestoreBackup)
		admin.GET("/audit-logs", h.ListAuditLogs)
		admin.GET("/audit-logs/export.csv", h.ExportAuditLogs)
		admin.GET("/patient-field-requirements", h.GetPatientFieldRequirements)
		admin.PUT("/patient-field-requirements", h.UpdatePatientFieldRequirements)
	}
}

func (h *AdminHandler) GetSystemStatus(c *gin.Context) {
	status, err := h.statsReader.GetSystemStatus(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *AdminHandler) RunHealthCheck(c *gin.Context) {
	result, err := h.healthChecker.RunHealthCheck(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logAdminEvent(c, auditdomain.EventHealthCheckRun, result.Status != "critical", map[string]string{
		"status": result.Status,
		"checks": strconv.Itoa(len(result.Checks)),
	})

	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) RunComplianceCheck(c *gin.Context) {
	result, err := h.complianceChecker.RunComplianceCheck(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logAdminEvent(c, auditdomain.EventComplianceCheckRun, result.Status != "fail", map[string]string{
		"status": result.Status,
		"checks": strconv.Itoa(len(result.Checks)),
	})

	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) GetPerformance(c *gin.Context) {
	if h.performanceReader == nil {
		httperr.RespondValidation(c, h.logger, "performance metrics require a running server")
		return
	}
	c.JSON(http.StatusOK, h.performanceReader.Snapshot(c.Request.Context()))
}

func (h *AdminHandler) GetPatientCensus(c *gin.Context) {
	census, err := h.statsReader.GetPatientCensus(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, census)
}

func (h *AdminHandler) GetActivitySummary(c *gin.Context) {
	periodDays, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	summary, err := h.statsReader.GetActivitySummary(c.Request.Context(), periodDays)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *AdminHandler) GetStaffActivity(c *gin.Context) {
	periodDays, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	report, err := h.statsReader.GetStaffActivity(c.Request.Context(), periodDays)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *AdminHandler) ListBackups(c *gin.Context) {
	if h.backupManager == nil {
		httperr.RespondValidation(c, h.logger, "backup is unavailable without persistent database storage")
		return
	}

	backups, err := h.backupManager.ListBackups(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"backups": backups})
}

func (h *AdminHandler) CreateBackup(c *gin.Context) {
	if h.backupManager == nil {
		httperr.RespondValidation(c, h.logger, "backup is unavailable without persistent database storage")
		return
	}

	result, err := h.backupManager.CreateBackup(c.Request.Context())
	if err != nil {
		h.logAdminEvent(c, auditdomain.EventBackupCreated, false, map[string]string{"error": "backup failed"})
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logAdminEvent(c, auditdomain.EventBackupCreated, true, map[string]string{"filename": result.Filename})
	c.JSON(http.StatusCreated, result)
}

type verifyBackupRequest struct {
	Filename string `json:"filename" binding:"required"`
}

func (h *AdminHandler) VerifyBackup(c *gin.Context) {
	if h.backupManager == nil {
		httperr.RespondValidation(c, h.logger, "backup verification is unavailable without persistent database storage")
		return
	}

	var req verifyBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	result, err := h.backupManager.VerifyBackup(c.Request.Context(), req.Filename)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) RestoreBackup(c *gin.Context) {
	if h.backupManager == nil {
		httperr.RespondValidation(c, h.logger, "restore is unavailable without persistent database storage")
		return
	}

	var req restoreBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	if !req.Confirm {
		httperr.RespondValidation(c, h.logger, "confirm must be true to restore")
		return
	}

	if err := h.backupManager.RestoreBackup(c.Request.Context(), req.Filename); err != nil {
		h.logAdminEvent(c, auditdomain.EventBackupRestored, false, map[string]string{
			"filename": req.Filename,
			"error":    "restore failed",
		})
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logAdminEvent(c, auditdomain.EventBackupRestored, true, map[string]string{"filename": req.Filename})
	c.JSON(http.StatusOK, gin.H{"message": "restore completed", "filename": req.Filename})
}

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "25"))

	result, err := h.listAuditHandler.Handle(c.Request.Context(), auditqueries.ListAuditQuery{
		Page:      page,
		PageSize:  pageSize,
		EventType: c.Query("eventType"),
		UserID:    c.Query("userId"),
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) ExportAuditLogs(c *gin.Context) {
	if h.auditRepository == nil {
		httperr.RespondValidation(c, h.logger, "audit export requires persistent database storage")
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}

	events, _, err := h.auditRepository.ListRecent(c.Request.Context(), auditdomain.ListFilter{
		Page:      1,
		PageSize:  5000,
		EventType: c.Query("eventType"),
		UserID:    c.Query("userId"),
		SinceDays: days,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	filename := fmt.Sprintf("audit-log-%s.csv", time.Now().UTC().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	writer := csv.NewWriter(c.Writer)
	_ = writer.Write([]string{
		"id", "event_type", "user_id", "resource_type", "resource_id",
		"ip_address", "success", "created_at", "details",
	})
	for _, event := range events {
		userID := ""
		if event.UserID != nil {
			userID = event.UserID.String()
		}
		details := ""
		if len(event.Details) > 0 {
			parts := make([]string, 0, len(event.Details))
			for key, value := range event.Details {
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			}
			details = fmt.Sprintf("%v", parts)
		}
		_ = writer.Write([]string{
			event.ID.String(),
			event.EventType,
			userID,
			event.ResourceType,
			event.ResourceID,
			event.IPAddress,
			strconv.FormatBool(event.Success),
			event.CreatedAt.UTC().Format(time.RFC3339),
			details,
		})
	}
	writer.Flush()
}

func (h *AdminHandler) GetPatientFieldRequirements(c *gin.Context) {
	reqs, err := h.getFieldRequirementsHandler.Handle(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"requirements": reqs})
}

type updatePatientFieldRequirementsBody struct {
	Requirements patientdomain.PatientFieldRequirements `json:"requirements"`
}

func (h *AdminHandler) UpdatePatientFieldRequirements(c *gin.Context) {
	var body updatePatientFieldRequirementsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	reqs, err := h.updateFieldRequirementsHandler.Handle(c.Request.Context(), patientcommands.UpdateFieldRequirementsCommand{
		Requirements: body.Requirements,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logAdminEvent(c, "patient_field_requirements_updated", true, map[string]string{
		"action": "update_patient_field_requirements",
	})
	c.JSON(http.StatusOK, gin.H{"requirements": reqs})
}

func (h *AdminHandler) logAdminEvent(c *gin.Context, eventType string, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}

	var userID *uuid.UUID
	if rawUserID, ok := c.Get("userID"); ok {
		if parsed, err := uuid.Parse(rawUserID.(string)); err == nil {
			userID = &parsed
		}
	}

	h.auditLogger.Log(c.Request.Context(), auditdomain.Event{
		EventType:    eventType,
		UserID:       userID,
		ResourceType: "system",
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		Details:      details,
		Success:      success,
	})
}
