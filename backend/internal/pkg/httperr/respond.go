package httperr

import (
	"errors"
	"net/http"
	"strings"

	authdomain "github.com/dksch/pococlinic/internal/features/auth/domain"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
)

// Response is the standard API error payload.
type Response struct {
	Code    string `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// MappedError contains the HTTP response for an error.
type MappedError struct {
	Status  int
	Code    string
	Message string
	Log     bool
}

// Map converts an application error into an HTTP response descriptor.
func Map(err error) MappedError {
	if err == nil {
		return internalError("An internal error occurred")
	}

	if mapped, ok := mapAuthError(err); ok {
		return mapped
	}

	var apiErr *pkgerrors.APIError
	if errors.As(err, &apiErr) {
		return MappedError{
			Status:  statusForCode(apiErr.Code),
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Log:     apiErr.Code == pkgerrors.ErrInternalServer,
		}
	}

	message := err.Error()
	switch message {
	case "form template not found", "form group not found", "form entry not found", "patient not found", "Patient not found":
		return MappedError{Status: http.StatusNotFound, Code: pkgerrors.ErrNotFound, Message: message}
	case "invalid patient id", "invalid template id", "invalid group id", "invalid entry id", "invalid user id":
		return MappedError{Status: http.StatusBadRequest, Code: pkgerrors.ErrValidation, Message: message}
	case "missing user context":
		return MappedError{Status: http.StatusUnauthorized, Code: pkgerrors.ErrUnauthorized, Message: message}
	case "user not found":
		return MappedError{Status: http.StatusNotFound, Code: pkgerrors.ErrNotFound, Message: message}
	case "cannot delete the default form group":
		return MappedError{Status: http.StatusBadRequest, Code: pkgerrors.ErrValidation, Message: message}
	case "invalid request body", "Invalid request body", "invalid query parameters", "Invalid page number", "Invalid page size":
		return MappedError{Status: http.StatusBadRequest, Code: pkgerrors.ErrValidation, Message: message}
	}

	if strings.HasPrefix(message, "form groups are not initialized") {
		return MappedError{
			Status:  http.StatusServiceUnavailable,
			Code:    pkgerrors.ErrInternalServer,
			Message: "Form storage is not ready yet. Restart the server or check database migrations.",
			Log:     true,
		}
	}

	if strings.HasPrefix(message, "field ") ||
		strings.HasPrefix(message, "form type must") ||
		strings.HasPrefix(message, "at least one field") ||
		strings.HasPrefix(message, "duplicate field") ||
		strings.HasPrefix(message, "entry does not match") ||
		strings.HasPrefix(message, "invalid ") ||
		strings.HasPrefix(message, "email ") {
		return MappedError{Status: http.StatusBadRequest, Code: pkgerrors.ErrValidation, Message: message}
	}

	if strings.Contains(message, "SQLSTATE") || strings.Contains(strings.ToLower(message), "database") {
		return internalError("A database error occurred while processing the request")
	}

	return MappedError{
		Status:  http.StatusBadRequest,
		Code:    pkgerrors.ErrValidation,
		Message: "The request could not be processed",
		Log:     true,
	}
}

// MapAuth maps authentication flow errors, defaulting unknown failures to invalid credentials.
func MapAuth(err error) MappedError {
	if mapped, ok := mapAuthError(err); ok {
		return mapped
	}

	var apiErr *pkgerrors.APIError
	if errors.As(err, &apiErr) {
		return Map(err)
	}

	return MappedError{
		Status:  http.StatusUnauthorized,
		Code:    authdomain.ErrInvalidCredentials,
		Message: "invalid credentials",
	}
}

// MapUserManagement maps staff account management errors.
func MapUserManagement(err error) MappedError {
	if mapped, ok := mapAuthError(err); ok {
		return mapped
	}

	var apiErr *pkgerrors.APIError
	if errors.As(err, &apiErr) {
		return Map(err)
	}

	return internalError("An internal error occurred")
}

func mapAuthError(err error) (MappedError, bool) {
	var authErr *authdomain.AuthError
	if !errors.As(err, &authErr) {
		return MappedError{}, false
	}

	switch authErr.Code {
	case authdomain.ErrInvalidCredentials, authdomain.ErrInvalidPIN, authdomain.ErrInvalidToken, authdomain.ErrSessionExpired:
		return MappedError{Status: http.StatusUnauthorized, Code: authErr.Code, Message: authErr.Message}, true
	case authdomain.ErrAccountLocked, authdomain.ErrAccountInactive:
		return MappedError{Status: http.StatusForbidden, Code: authErr.Code, Message: authErr.Message}, true
	case authdomain.ErrUserNotFound, authdomain.ErrSessionNotFound:
		return MappedError{Status: http.StatusNotFound, Code: authErr.Code, Message: authErr.Message}, true
	case authdomain.ErrEmailTaken:
		return MappedError{Status: http.StatusConflict, Code: authErr.Code, Message: authErr.Message}, true
	case authdomain.ErrLastAdmin, authdomain.ErrCannotDeleteSelf:
		return MappedError{Status: http.StatusForbidden, Code: authErr.Code, Message: authErr.Message}, true
	default:
		return MappedError{Status: http.StatusBadRequest, Code: authErr.Code, Message: authErr.Message}, true
	}
}

// Respond writes a mapped error response and logs server-side failures.
func Respond(c *gin.Context, logger *logging.Logger, err error) {
	writeMapped(c, logger, Map(err), err)
}

// RespondAuth writes an authentication-specific error response.
func RespondAuth(c *gin.Context, logger *logging.Logger, err error) {
	writeMapped(c, logger, MapAuth(err), err)
}

// RespondUserManagement writes a staff management error response.
func RespondUserManagement(c *gin.Context, logger *logging.Logger, err error) {
	writeMapped(c, logger, MapUserManagement(err), err)
}

// RespondValidation writes a validation error for malformed requests.
func RespondValidation(c *gin.Context, logger *logging.Logger, message string) {
	writeMapped(c, logger, MappedError{
		Status:  http.StatusBadRequest,
		Code:    pkgerrors.ErrValidation,
		Message: message,
	}, nil)
}

func writeMapped(c *gin.Context, logger *logging.Logger, mapped MappedError, err error) {
	if mapped.Log && logger != nil && err != nil {
		logger.Error("request failed",
			err,
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", mapped.Status,
		)
	}

	c.JSON(mapped.Status, Response{
		Code:    mapped.Code,
		Error:   mapped.Message,
		Message: mapped.Message,
	})
}

func internalError(message string) MappedError {
	return MappedError{
		Status:  http.StatusInternalServerError,
		Code:    pkgerrors.ErrInternalServer,
		Message: message,
		Log:     true,
	}
}

func statusForCode(code string) int {
	switch code {
	case pkgerrors.ErrValidation:
		return http.StatusBadRequest
	case pkgerrors.ErrNotFound:
		return http.StatusNotFound
	case pkgerrors.ErrUnauthorized:
		return http.StatusUnauthorized
	case pkgerrors.ErrForbidden:
		return http.StatusForbidden
	case pkgerrors.ErrRateLimit:
		return http.StatusTooManyRequests
	case authdomain.ErrEmailTaken:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
