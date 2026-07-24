package handlers

import (
	"net/http"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	authdomain "github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/PococodoOrg/PocoClinic/internal/features/auth/middleware"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/commands"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/queries"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/httperr"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExerciseLogHandler serves physical-therapy exercise plans and session logs.
type ExerciseLogHandler struct {
	commands    *commands.ExerciseLogCommandHandler
	queries     *queries.ExerciseLogQueryHandler
	logger      *logging.Logger
	auditLogger auditdomain.Logger
}

func NewExerciseLogHandler(
	commandHandler *commands.ExerciseLogCommandHandler,
	queryHandler *queries.ExerciseLogQueryHandler,
	logger *logging.Logger,
	auditLogger auditdomain.Logger,
) *ExerciseLogHandler {
	return &ExerciseLogHandler{
		commands:    commandHandler,
		queries:     queryHandler,
		logger:      logger,
		auditLogger: auditLogger,
	}
}

func (h *ExerciseLogHandler) RegisterRoutes(router *gin.RouterGroup, auth *middleware.AuthMiddleware) {
	patients := router.Group("/patients")
	if auth != nil {
		patients.Use(auth.RequireAuth(), auth.RequirePINChanged())
	}

	readers := authdomain.ChartReaderRoles()
	clinical := authdomain.ClinicalStaffRoles()

	patients.GET("/:id/exercise-plans", roleHandlers(auth, readers, h.ListPlans)...)
	patients.POST("/:id/exercise-plans", roleHandlers(auth, clinical, h.CreatePlan)...)
	patients.PUT("/:id/exercise-plans/:planId", roleHandlers(auth, clinical, h.UpdatePlan)...)
	patients.DELETE("/:id/exercise-plans/:planId", roleHandlers(auth, clinical, h.DeletePlan)...)
	patients.GET("/:id/exercise-plans/:planId/entries", roleHandlers(auth, readers, h.ListEntries)...)
	patients.POST("/:id/exercise-plans/:planId/entries", roleHandlers(auth, clinical, h.CreateEntry)...)
	patients.PUT("/:id/exercise-plans/:planId/entries/:entryId", roleHandlers(auth, clinical, h.UpdateEntry)...)
	patients.DELETE("/:id/exercise-plans/:planId/entries/:entryId", roleHandlers(auth, clinical, h.DeleteEntry)...)
}

func (h *ExerciseLogHandler) requirePatientID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid patient id")
		return "", false
	}
	return id, true
}

func (h *ExerciseLogHandler) requirePlanID(c *gin.Context) (string, bool) {
	id := c.Param("planId")
	if _, err := uuid.Parse(id); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid plan id")
		return "", false
	}
	return id, true
}

func (h *ExerciseLogHandler) requireEntryID(c *gin.Context) (string, bool) {
	id := c.Param("entryId")
	if _, err := uuid.Parse(id); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid entry id")
		return "", false
	}
	return id, true
}

func (h *ExerciseLogHandler) authorIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	raw, ok := c.Get("userID")
	if !ok {
		httperr.RespondValidation(c, h.logger, "missing authenticated user")
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw.(string))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid authenticated user")
		return uuid.Nil, false
	}
	return id, true
}

func (h *ExerciseLogHandler) ListPlans(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	plans, err := h.queries.ListPlans(c.Request.Context(), queries.ListExercisePlansQuery{PatientID: patientID})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *ExerciseLogHandler) CreatePlan(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	authorID, ok := h.authorIDFromContext(c)
	if !ok {
		return
	}
	var cmd commands.CreateExercisePlanCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.PatientID = patientID
	cmd.CreatedBy = authorID

	plan, err := h.commands.CreatePlan(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusCreated, plan)
	h.audit(c, auditdomain.EventPatientExercisePlanCreated, patientID, true, map[string]string{
		"planId": plan.ID.String(),
		"name":   plan.Name,
	})
}

func (h *ExerciseLogHandler) UpdatePlan(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	planID, ok := h.requirePlanID(c)
	if !ok {
		return
	}
	var cmd commands.UpdateExercisePlanCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.PatientID = patientID
	cmd.PlanID = planID

	plan, err := h.commands.UpdatePlan(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, plan)
	h.audit(c, auditdomain.EventPatientExercisePlanUpdated, patientID, true, map[string]string{
		"planId": plan.ID.String(),
		"status": plan.Status,
	})
}

func (h *ExerciseLogHandler) DeletePlan(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	planID, ok := h.requirePlanID(c)
	if !ok {
		return
	}
	if err := h.commands.DeletePlan(c.Request.Context(), commands.DeleteExercisePlanCommand{
		PatientID: patientID,
		PlanID:    planID,
	}); err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.Status(http.StatusNoContent)
	h.audit(c, auditdomain.EventPatientExercisePlanDeleted, patientID, true, map[string]string{
		"planId": planID,
	})
}

func (h *ExerciseLogHandler) ListEntries(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	planID, ok := h.requirePlanID(c)
	if !ok {
		return
	}
	entries, err := h.queries.ListEntries(c.Request.Context(), queries.ListExerciseEntriesQuery{
		PatientID: patientID,
		PlanID:    planID,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *ExerciseLogHandler) CreateEntry(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	planID, ok := h.requirePlanID(c)
	if !ok {
		return
	}
	authorID, ok := h.authorIDFromContext(c)
	if !ok {
		return
	}
	var cmd commands.CreateExerciseEntryCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.PatientID = patientID
	cmd.PlanID = planID
	cmd.RecordedBy = authorID

	entry, err := h.commands.CreateEntry(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusCreated, entry)
	h.audit(c, auditdomain.EventPatientExerciseEntryCreated, patientID, true, map[string]string{
		"planId":       planID,
		"entryId":      entry.ID.String(),
		"exerciseName": entry.ExerciseName,
	})
}

func (h *ExerciseLogHandler) UpdateEntry(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	planID, ok := h.requirePlanID(c)
	if !ok {
		return
	}
	entryID, ok := h.requireEntryID(c)
	if !ok {
		return
	}
	var cmd commands.UpdateExerciseEntryCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.PatientID = patientID
	cmd.PlanID = planID
	cmd.EntryID = entryID

	entry, err := h.commands.UpdateEntry(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, entry)
	h.audit(c, auditdomain.EventPatientExerciseEntryUpdated, patientID, true, map[string]string{
		"planId":  planID,
		"entryId": entryID,
	})
}

func (h *ExerciseLogHandler) DeleteEntry(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	planID, ok := h.requirePlanID(c)
	if !ok {
		return
	}
	entryID, ok := h.requireEntryID(c)
	if !ok {
		return
	}
	if err := h.commands.DeleteEntry(c.Request.Context(), commands.DeleteExerciseEntryCommand{
		PatientID: patientID,
		PlanID:    planID,
		EntryID:   entryID,
	}); err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.Status(http.StatusNoContent)
	h.audit(c, auditdomain.EventPatientExerciseEntryDeleted, patientID, true, map[string]string{
		"planId":  planID,
		"entryId": entryID,
	})
}

func (h *ExerciseLogHandler) audit(c *gin.Context, eventType, resourceID string, success bool, details map[string]string) {
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
		ResourceType: "patient",
		ResourceID:   resourceID,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		Details:      details,
		Success:      success,
	})
}
