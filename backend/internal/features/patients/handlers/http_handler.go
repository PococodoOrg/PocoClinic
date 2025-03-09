package handlers

import (
	"net/http"
	"strconv"

	"github.com/dksch/pococlinic/internal/features/patients/commands"
	"github.com/dksch/pococlinic/internal/features/patients/queries"
	"github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
)

// PatientHandler handles HTTP requests for patient operations
type PatientHandler struct {
	createPatientHandler commands.CreatePatientHandler
	getPatientsHandler   queries.GetPatientsHandler
	logger               *logging.Logger
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(
	createHandler commands.CreatePatientHandler,
	getHandler queries.GetPatientsHandler,
	logger *logging.Logger,
) *PatientHandler {
	return &PatientHandler{
		createPatientHandler: createHandler,
		getPatientsHandler:   getHandler,
		logger:               logger,
	}
}

// RegisterRoutes registers the patient routes with the given router group
func (h *PatientHandler) RegisterRoutes(router *gin.RouterGroup) {
	patients := router.Group("/patients")
	{
		patients.POST("", h.CreatePatient)
		patients.GET("", h.ListPatients)
		// Add more routes as needed
	}
}

// CreatePatient handles the creation of a new patient
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var cmd commands.CreatePatientCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		h.logger.Error("Invalid request body", err)
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrValidation, "Invalid request body"))
		return
	}

	patient, err := h.createPatientHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to create patient", err)
		// Check for specific error types and return appropriate status codes
		switch err.(type) {
		case *errors.APIError:
			apiErr := err.(*errors.APIError)
			c.JSON(getStatusCodeForError(apiErr.Code), apiErr)
		default:
			c.JSON(http.StatusInternalServerError, errors.NewAPIError(errors.ErrInternalServer, "Failed to create patient"))
		}
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// ListPatients handles retrieving a paginated list of patients
func (h *PatientHandler) ListPatients(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrValidation, "Invalid page number"))
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrValidation, "Invalid page size"))
		return
	}

	search := c.Query("search")

	query := queries.GetPatientsQuery{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	result, err := h.getPatientsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to get patients", err)
		switch err.(type) {
		case *errors.APIError:
			apiErr := err.(*errors.APIError)
			c.JSON(getStatusCodeForError(apiErr.Code), apiErr)
		default:
			c.JSON(http.StatusInternalServerError, errors.NewAPIError(errors.ErrInternalServer, "Failed to retrieve patients"))
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// getStatusCodeForError returns the appropriate HTTP status code for an error code
func getStatusCodeForError(code string) int {
	switch code {
	case errors.ErrValidation:
		return http.StatusBadRequest
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrForbidden:
		return http.StatusForbidden
	case errors.ErrRateLimit:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
