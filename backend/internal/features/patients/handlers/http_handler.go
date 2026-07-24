package handlers

import (
	"io"
	"net/http"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	authdomain "github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/features/auth/middleware"
	"github.com/dksch/pococlinic/internal/features/patients/commands"
	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/features/patients/queries"
	"github.com/dksch/pococlinic/internal/pkg/filename"
	"github.com/dksch/pococlinic/internal/pkg/httperr"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


// PatientHandler handles HTTP requests for patient operations
type PatientHandler struct {
	createPatientHandler  commands.CreatePatientHandler
	getPatientsHandler    queries.GetPatientsHandler
	getPatientHandler     queries.GetPatientHandler
	updatePatientHandler  commands.UpdatePatientHandler
	deletePatientHandler  commands.DeletePatientHandler
	createNoteHandler     commands.CreateNoteHandler
	updateNoteHandler     commands.UpdateNoteHandler
	deleteNoteHandler     commands.DeleteNoteHandler
	listNotesHandler      queries.ListNotesHandler
	uploadDocumentHandler      commands.UploadDocumentHandler
	deleteDocumentHandler      commands.DeleteDocumentHandler
	listDocumentsHandler       queries.ListDocumentsHandler
	getDocumentHandler         queries.GetDocumentHandler
	openDocumentContentHandler commands.OpenDocumentContentHandler
	getFieldRequirementsHandler  queries.GetFieldRequirementsHandler
	logger                     *logging.Logger
	auditLogger                auditdomain.Logger
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(
	createHandler commands.CreatePatientHandler,
	getHandler queries.GetPatientsHandler,
	getPatientHandler queries.GetPatientHandler,
	updateHandler commands.UpdatePatientHandler,
	deleteHandler commands.DeletePatientHandler,
	createNoteHandler commands.CreateNoteHandler,
	updateNoteHandler commands.UpdateNoteHandler,
	deleteNoteHandler commands.DeleteNoteHandler,
	listNotesHandler queries.ListNotesHandler,
	uploadDocumentHandler commands.UploadDocumentHandler,
	deleteDocumentHandler commands.DeleteDocumentHandler,
	listDocumentsHandler queries.ListDocumentsHandler,
	getDocumentHandler queries.GetDocumentHandler,
	openDocumentContentHandler commands.OpenDocumentContentHandler,
	getFieldRequirementsHandler queries.GetFieldRequirementsHandler,
	logger *logging.Logger,
	auditLogger auditdomain.Logger,
) *PatientHandler {
	return &PatientHandler{
		createPatientHandler:       createHandler,
		getPatientsHandler:         getHandler,
		getPatientHandler:          getPatientHandler,
		updatePatientHandler:       updateHandler,
		deletePatientHandler:       deleteHandler,
		createNoteHandler:          createNoteHandler,
		updateNoteHandler:          updateNoteHandler,
		deleteNoteHandler:          deleteNoteHandler,
		listNotesHandler:           listNotesHandler,
		uploadDocumentHandler:      uploadDocumentHandler,
		deleteDocumentHandler:      deleteDocumentHandler,
		listDocumentsHandler:       listDocumentsHandler,
		getDocumentHandler:         getDocumentHandler,
		openDocumentContentHandler: openDocumentContentHandler,
		getFieldRequirementsHandler:  getFieldRequirementsHandler,
		logger:                     logger,
		auditLogger:                auditLogger,
	}
}

// RegisterRoutes registers the patient routes with the given router group.
func (h *PatientHandler) RegisterRoutes(router *gin.RouterGroup, auth *middleware.AuthMiddleware) {
	patients := router.Group("/patients")
	if auth != nil {
		patients.Use(auth.RequireAuth(), auth.RequirePINChanged())
	}

	readers := authdomain.ChartReaderRoles()
	clinical := authdomain.ClinicalStaffRoles()
	uploaders := authdomain.DocumentUploaderRoles()
	deleters := authdomain.DocumentDeleterRoles()

	patients.GET("/field-requirements", roleHandlers(auth, readers, h.GetFieldRequirements)...)
	patients.GET("", roleHandlers(auth, readers, h.ListPatients)...)
	patients.GET("/:id", roleHandlers(auth, readers, h.GetPatient)...)
	patients.POST("", roleHandlers(auth, clinical, h.CreatePatient)...)
	patients.PUT("/:id", roleHandlers(auth, clinical, h.UpdatePatient)...)
	patients.DELETE("/:id", roleHandlers(auth, []authdomain.Role{authdomain.RoleAdmin}, h.DeletePatient)...)
	patients.GET("/:id/notes", roleHandlers(auth, readers, h.ListNotes)...)
	patients.POST("/:id/notes", roleHandlers(auth, clinical, h.CreateNote)...)
	patients.PUT("/:id/notes/:noteId", roleHandlers(auth, clinical, h.UpdateNote)...)
	patients.DELETE("/:id/notes/:noteId", roleHandlers(auth, clinical, h.DeleteNote)...)
	patients.GET("/:id/documents", roleHandlers(auth, readers, h.ListDocuments)...)
	patients.POST("/:id/documents", roleHandlers(auth, uploaders, h.UploadDocument)...)
	patients.GET("/:id/documents/:documentId/download", roleHandlers(auth, readers, h.DownloadDocument)...)
	patients.DELETE("/:id/documents/:documentId", roleHandlers(auth, deleters, h.DeleteDocument)...)
}

func roleHandlers(auth *middleware.AuthMiddleware, roles []authdomain.Role, handler gin.HandlerFunc) []gin.HandlerFunc {
	if auth == nil || len(roles) == 0 {
		return []gin.HandlerFunc{handler}
	}
	return []gin.HandlerFunc{auth.RequireRole(roles...), handler}
}

func (h *PatientHandler) requirePatientID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid patient id")
		return "", false
	}
	return id, true
}

func (h *PatientHandler) requireNoteID(c *gin.Context) (string, bool) {
	id := c.Param("noteId")
	if _, err := uuid.Parse(id); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid note id")
		return "", false
	}
	return id, true
}

// GetFieldRequirements returns clinic-configured required patient fields.
func (h *PatientHandler) GetFieldRequirements(c *gin.Context) {
	reqs, err := h.getFieldRequirementsHandler.Handle(c.Request.Context())
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"requirements": reqs})
}

func (h *PatientHandler) requireDocumentID(c *gin.Context) (string, bool) {
	id := c.Param("documentId")
	if _, err := uuid.Parse(id); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid document id")
		return "", false
	}
	return id, true
}

// CreatePatient handles the creation of a new patient
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var cmd commands.CreatePatientCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	h.logger.Info("Creating patient")

	patient, err := h.createPatientHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logger.Info("Created patient", "id", patient.ID)

	c.JSON(http.StatusCreated, patient)
	h.auditPatientEvent(c, auditdomain.EventPatientCreated, patient.ID.String(), true, nil)
}

// ListPatients handles retrieving a paginated list of patients
func (h *PatientHandler) ListPatients(c *gin.Context) {
	var query queries.GetPatientsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid query parameters")
		return
	}

	if query.Page < 1 {
		httperr.RespondValidation(c, h.logger, "invalid page number")
		return
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		httperr.RespondValidation(c, h.logger, "invalid page size")
		return
	}
	if !domain.ValidGender(query.Gender) {
		httperr.RespondValidation(c, h.logger, "invalid gender")
		return
	}

	result, err := h.getPatientsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPatient handles the request to fetch a single patient by ID
func (h *PatientHandler) GetPatient(c *gin.Context) {
	id, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	query := queries.GetPatientQuery{ID: id}

	patient, err := h.getPatientHandler.Handle(c.Request.Context(), query)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, patient)
	h.auditPatientEvent(c, auditdomain.EventPatientViewed, query.ID, true, nil)
}

func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	id, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	var cmd commands.UpdatePatientCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	cmd.ID = id

	h.logger.Info("Updating patient", "id", id)

	patient, err := h.updatePatientHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	h.logger.Info("Updated patient", "id", patient.ID)

	c.JSON(http.StatusOK, patient)
	h.auditPatientEvent(c, auditdomain.EventPatientUpdated, patient.ID.String(), true, nil)
}

func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id, ok := h.requirePatientID(c)
	if !ok {
		return
	}

	if err := h.deletePatientHandler.Handle(c.Request.Context(), id); err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.Status(http.StatusNoContent)
	h.auditPatientEvent(c, auditdomain.EventPatientDeleted, id, true, nil)
}

func (h *PatientHandler) ListNotes(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	notes, err := h.listNotesHandler.Handle(c.Request.Context(), queries.ListNotesQuery{PatientID: patientID})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	if notes == nil {
		notes = []*domain.Note{}
	}
	c.JSON(http.StatusOK, notes)
}

func (h *PatientHandler) CreateNote(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	var body struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	rawUserID, ok := c.Get("userID")
	if !ok {
		httperr.RespondValidation(c, h.logger, "missing user context")
		return
	}
	authorID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user context")
		return
	}

	note, err := h.createNoteHandler.Handle(c.Request.Context(), commands.CreateNoteCommand{
		PatientID: patientID,
		AuthorID:  authorID,
		Body:      body.Body,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.JSON(http.StatusCreated, note)
	h.auditPatientEvent(c, auditdomain.EventPatientNoteCreated, patientID, true, map[string]string{
		"noteId": note.ID.String(),
	})
}

func (h *PatientHandler) UpdateNote(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	noteID, ok := h.requireNoteID(c)
	if !ok {
		return
	}
	var body struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	authorID, ok := h.authorIDFromContext(c)
	if !ok {
		return
	}

	note, err := h.updateNoteHandler.Handle(c.Request.Context(), commands.UpdateNoteCommand{
		PatientID: patientID,
		NoteID:    noteID,
		AuthorID:  authorID,
		Body:      body.Body,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		h.auditPatientEvent(c, auditdomain.EventPatientNoteUpdated, patientID, false, map[string]string{
			"noteId": noteID,
		})
		return
	}

	c.JSON(http.StatusOK, note)
	h.auditPatientEvent(c, auditdomain.EventPatientNoteUpdated, patientID, true, map[string]string{
		"noteId": note.ID.String(),
	})
}

func (h *PatientHandler) DeleteNote(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	noteID, ok := h.requireNoteID(c)
	if !ok {
		return
	}

	actorID, ok := h.authorIDFromContext(c)
	if !ok {
		return
	}

	note, err := h.deleteNoteHandler.Handle(c.Request.Context(), commands.DeleteNoteCommand{
		PatientID: patientID,
		NoteID:    noteID,
		ActorID:   actorID,
		IsAdmin:   h.isAdminFromContext(c),
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		h.auditPatientEvent(c, auditdomain.EventPatientNoteDeleted, patientID, false, map[string]string{
			"noteId": noteID,
		})
		return
	}

	c.Status(http.StatusNoContent)
	h.auditPatientEvent(c, auditdomain.EventPatientNoteDeleted, patientID, true, map[string]string{
		"noteId": note.ID.String(),
	})
}

func (h *PatientHandler) authorIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	rawUserID, ok := c.Get("userID")
	if !ok {
		httperr.RespondValidation(c, h.logger, "missing user context")
		return uuid.Nil, false
	}
	authorID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user context")
		return uuid.Nil, false
	}
	return authorID, true
}

func (h *PatientHandler) isAdminFromContext(c *gin.Context) bool {
	rawRole, ok := c.Get("userRole")
	if !ok {
		return false
	}
	role, ok := rawRole.(authdomain.Role)
	return ok && role == authdomain.RoleAdmin
}

func (h *PatientHandler) ListDocuments(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	docs, err := h.listDocumentsHandler.Handle(c.Request.Context(), queries.ListDocumentsQuery{PatientID: patientID})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	if docs == nil {
		docs = []*domain.Document{}
	}
	c.JSON(http.StatusOK, docs)
}

func (h *PatientHandler) UploadDocument(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		httperr.RespondValidation(c, h.logger, "File is required")
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	defer src.Close()

	rawUserID, ok := c.Get("userID")
	if !ok {
		httperr.RespondValidation(c, h.logger, "missing user context")
		return
	}
	uploaderID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user context")
		return
	}

	doc, err := h.uploadDocumentHandler.Handle(c.Request.Context(), commands.UploadDocumentCommand{
		PatientID:   patientID,
		UploadedBy:  uploaderID,
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Content:     src,
	})
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.JSON(http.StatusCreated, doc)
	h.auditPatientEvent(c, auditdomain.EventPatientDocumentUploaded, patientID, true, map[string]string{
		"documentId": doc.ID.String(),
		"fileName":   doc.FileName,
	})
}

func (h *PatientHandler) DownloadDocument(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	documentID, ok := h.requireDocumentID(c)
	if !ok {
		return
	}

	reader, doc, err := h.openDocumentContentHandler.Handle(c.Request.Context(), patientID, documentID)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}
	defer reader.Close()

	c.Header("Content-Type", doc.ContentType)
	c.Header("Content-Disposition", filename.ContentDispositionAttachment(doc.FileName))
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, reader); err != nil {
		h.logger.Error("Failed to stream document", err)
		return
	}

	h.auditPatientEvent(c, auditdomain.EventPatientDocumentViewed, patientID, true, map[string]string{
		"documentId": doc.ID.String(),
		"fileName":   doc.FileName,
	})
}

func (h *PatientHandler) DeleteDocument(c *gin.Context) {
	patientID, ok := h.requirePatientID(c)
	if !ok {
		return
	}
	documentID, ok := h.requireDocumentID(c)
	if !ok {
		return
	}

	doc, err := h.deleteDocumentHandler.Handle(c.Request.Context(), patientID, documentID)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.Status(http.StatusNoContent)
	h.auditPatientEvent(c, auditdomain.EventPatientDocumentDeleted, patientID, true, map[string]string{
		"documentId": doc.ID.String(),
		"fileName":   doc.FileName,
	})
}

func (h *PatientHandler) auditPatientEvent(c *gin.Context, eventType, resourceID string, success bool, details map[string]string) {
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
