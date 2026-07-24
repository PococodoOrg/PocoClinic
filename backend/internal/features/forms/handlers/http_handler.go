package handlers

import (
	"net/http"

	"strconv"

	auditdomain "github.com/PococodoOrg/PocoClinic/internal/features/audit/domain"
	authdomain "github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	"github.com/PococodoOrg/PocoClinic/internal/features/auth/middleware"
	"github.com/PococodoOrg/PocoClinic/internal/features/forms/commands"
	"github.com/PococodoOrg/PocoClinic/internal/features/forms/queries"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/httperr"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FormHandler struct {
	logger                       *logging.Logger
	auditLogger                  auditdomain.Logger
	createGroupHandler           commands.CreateGroupHandler
	updateGroupHandler           commands.UpdateGroupHandler
	deleteGroupHandler           commands.DeleteGroupHandler
	createTemplateHandler        commands.CreateTemplateHandler
	updateTemplateHandler        commands.UpdateTemplateHandler
	deleteTemplateHandler        commands.DeleteTemplateHandler
	saveFormHandler              commands.SaveFormHandler
	listGroupsHandler            queries.ListGroupsHandler
	getGroupHandler              queries.GetGroupHandler
	listTemplatesHandler         queries.ListTemplatesHandler
	getTemplateHandler           queries.GetTemplateHandler
	listSubmissionsHandler       queries.ListSubmissionsHandler
	listSubmissionHistoryHandler queries.ListSubmissionHistoryHandler
	listTemplateSubmissionsHandler queries.ListTemplateSubmissionsHandler
}

func NewFormHandler(
	logger *logging.Logger,
	auditLogger auditdomain.Logger,
	createGroup commands.CreateGroupHandler,
	updateGroup commands.UpdateGroupHandler,
	deleteGroup commands.DeleteGroupHandler,
	createTemplate commands.CreateTemplateHandler,
	updateTemplate commands.UpdateTemplateHandler,
	deleteTemplate commands.DeleteTemplateHandler,
	saveForm commands.SaveFormHandler,
	listGroups queries.ListGroupsHandler,
	getGroup queries.GetGroupHandler,
	listTemplates queries.ListTemplatesHandler,
	getTemplate queries.GetTemplateHandler,
	listSubmissions queries.ListSubmissionsHandler,
	listSubmissionHistory queries.ListSubmissionHistoryHandler,
	listTemplateSubmissions queries.ListTemplateSubmissionsHandler,
) *FormHandler {
	return &FormHandler{
		logger:                       logger,
		auditLogger:                  auditLogger,
		createGroupHandler:           createGroup,
		updateGroupHandler:           updateGroup,
		deleteGroupHandler:           deleteGroup,
		createTemplateHandler:        createTemplate,
		updateTemplateHandler:        updateTemplate,
		deleteTemplateHandler:        deleteTemplate,
		saveFormHandler:              saveForm,
		listGroupsHandler:            listGroups,
		getGroupHandler:              getGroup,
		listTemplatesHandler:         listTemplates,
		getTemplateHandler:           getTemplate,
		listSubmissionsHandler:       listSubmissions,
		listSubmissionHistoryHandler: listSubmissionHistory,
		listTemplateSubmissionsHandler: listTemplateSubmissions,
	}
}

func (h *FormHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	groups := router.Group("/form-groups")
	groups.Use(authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged())
	{
		groups.GET("", h.ListGroups)
		groups.GET("/:id", h.GetGroup)
		groups.POST("", authMiddleware.RequireRole(authdomain.RoleAdmin), h.CreateGroup)
		groups.PUT("/:id", authMiddleware.RequireRole(authdomain.RoleAdmin), h.UpdateGroup)
		groups.DELETE("/:id", authMiddleware.RequireRole(authdomain.RoleAdmin), h.DeleteGroup)
	}

	templates := router.Group("/form-templates")
	templates.Use(authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged())
	{
		templates.GET("", h.ListTemplates)
		templates.GET("/:id/submissions", authMiddleware.RequireRole(authdomain.RoleAdmin), h.ListTemplateSubmissions)
		templates.GET("/:id", h.GetTemplate)
		templates.POST("", authMiddleware.RequireRole(authdomain.RoleAdmin), h.CreateTemplate)
		templates.PUT("/:id", authMiddleware.RequireRole(authdomain.RoleAdmin), h.UpdateTemplate)
		templates.DELETE("/:id", authMiddleware.RequireRole(authdomain.RoleAdmin), h.DeleteTemplate)
	}

	patients := router.Group("/patients")
	patients.Use(authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged())
	{
		readers := authdomain.ChartReaderRoles()
		clinical := authdomain.ClinicalStaffRoles()
		patients.GET("/:id/form-submissions", authMiddleware.RequireRole(readers...), h.ListSubmissions)
		patients.POST("/:id/form-submissions", authMiddleware.RequireRole(clinical...), h.SaveForm)
		patients.GET("/:id/form-entries/:entryId/history", authMiddleware.RequireRole(readers...), h.ListSubmissionHistory)
	}
}

func (h *FormHandler) ListGroups(c *gin.Context) {
	groups, err := h.listGroupsHandler.Handle(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func (h *FormHandler) GetGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid group id")
		return
	}

	group, err := h.getGroupHandler.Handle(c.Request.Context(), queries.GetGroupQuery{ID: id.String()})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, group)
}

func (h *FormHandler) CreateGroup(c *gin.Context) {
	var cmd commands.CreateGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	group, err := h.createGroupHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusCreated, group)
}

func (h *FormHandler) UpdateGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid group id")
		return
	}

	var cmd commands.UpdateGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.ID = id.String()

	group, err := h.updateGroupHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, group)
}

func (h *FormHandler) DeleteGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid group id")
		return
	}

	if err := h.deleteGroupHandler.Handle(c.Request.Context(), commands.DeleteGroupCommand{ID: id.String()}); err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *FormHandler) ListTemplates(c *gin.Context) {
	templates, err := h.listTemplatesHandler.Handle(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *FormHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid template id")
		return
	}

	template, err := h.getTemplateHandler.Handle(c.Request.Context(), queries.GetTemplateQuery{ID: id.String()})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *FormHandler) CreateTemplate(c *gin.Context) {
	var cmd commands.CreateTemplateCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	template, err := h.createTemplateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusCreated, template)
}

func (h *FormHandler) UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid template id")
		return
	}

	var cmd commands.UpdateTemplateCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.ID = id.String()

	template, err := h.updateTemplateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *FormHandler) DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid template id")
		return
	}

	if err := h.deleteTemplateHandler.Handle(c.Request.Context(), commands.DeleteTemplateCommand{ID: id.String()}); err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *FormHandler) ListSubmissions(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid patient id")
		return
	}

	submissions, err := h.listSubmissionsHandler.Handle(c.Request.Context(), queries.ListSubmissionsQuery{
		PatientID: patientID.String(),
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"submissions": submissions})
}

func (h *FormHandler) SaveForm(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid patient id")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		httperr.Respond(c, h.logger, authdomain.NewAuthError(authdomain.ErrInvalidToken, "missing user context"))
		return
	}

	var cmd commands.SaveFormCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}
	cmd.PatientID = patientID.String()
	cmd.SubmittedBy = userID.(string)

	submission, err := h.saveFormHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusCreated, submission)
	h.auditFormSubmission(c, patientID.String(), submission.ID.String(), cmd.TemplateID, true)
}

func (h *FormHandler) auditFormSubmission(c *gin.Context, patientID, submissionID, templateID string, success bool) {
	if h.auditLogger == nil {
		return
	}

	var actorID *uuid.UUID
	if rawActorID, exists := c.Get("userID"); exists {
		if parsed, err := uuid.Parse(rawActorID.(string)); err == nil {
			actorID = &parsed
		}
	}

	h.auditLogger.Log(c.Request.Context(), auditdomain.Event{
		EventType:    auditdomain.EventFormSubmissionCreated,
		UserID:       actorID,
		ResourceType: "form_submission",
		ResourceID:   submissionID,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		Details: map[string]string{
			"patientId":  patientID,
			"templateId": templateID,
		},
		Success: success,
	})
}

func (h *FormHandler) ListSubmissionHistory(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid patient id")
		return
	}
	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid entry id")
		return
	}

	history, err := h.listSubmissionHistoryHandler.Handle(c.Request.Context(), queries.ListSubmissionHistoryQuery{
		PatientID: patientID.String(),
		EntryID:   entryID.String(),
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": history})
}

func (h *FormHandler) ListTemplateSubmissions(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid template id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "25"))

	var from, to *string
	if value := c.Query("from"); value != "" {
		from = &value
	}
	if value := c.Query("to"); value != "" {
		to = &value
	}

	report, err := h.listTemplateSubmissionsHandler.Handle(c.Request.Context(), queries.ListTemplateSubmissionsQuery{
		TemplateID: templateID.String(),
		Page:       page,
		PageSize:   pageSize,
		From:       from,
		To:         to,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, report)
}
