package errors

// APIError represents a standardized API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Common error codes
const (
	ErrInternalServer = "INTERNAL_ERROR"
	ErrValidation     = "VALIDATION_ERROR"
	ErrNotFound       = "NOT_FOUND"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrRateLimit      = "RATE_LIMIT_EXCEEDED"
)

// NewAPIError creates a new API error
func NewAPIError(code string, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}
