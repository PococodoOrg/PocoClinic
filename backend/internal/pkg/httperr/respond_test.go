package httperr

import (
	"errors"
	"net/http"
	"testing"

	authdomain "github.com/dksch/pococlinic/internal/features/auth/domain"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
)

func TestMapKnownErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
		wantLog    bool
	}{
		{
			name:       "not found",
			err:        errors.New("form template not found"),
			wantStatus: http.StatusNotFound,
			wantCode:   pkgerrors.ErrNotFound,
		},
		{
			name:       "validation",
			err:        errors.New("field \"Are you sick?\" is required"),
			wantStatus: http.StatusBadRequest,
			wantCode:   pkgerrors.ErrValidation,
		},
		{
			name:       "database",
			err:        errors.New(`ERROR: column "group_id" does not exist (SQLSTATE 42703)`),
			wantStatus: http.StatusInternalServerError,
			wantCode:   pkgerrors.ErrInternalServer,
			wantLog:    true,
		},
		{
			name:       "api error",
			err:        pkgerrors.NewAPIError(pkgerrors.ErrForbidden, "Admin only"),
			wantStatus: http.StatusForbidden,
			wantCode:   pkgerrors.ErrForbidden,
		},
		{
			name:       "patient api error",
			err:        pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Patient not found"),
			wantStatus: http.StatusNotFound,
			wantCode:   pkgerrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapped := Map(tt.err)
			if mapped.Status != tt.wantStatus {
				t.Fatalf("status = %d, want %d", mapped.Status, tt.wantStatus)
			}
			if mapped.Code != tt.wantCode {
				t.Fatalf("code = %q, want %q", mapped.Code, tt.wantCode)
			}
			if mapped.Log != tt.wantLog {
				t.Fatalf("log = %v, want %v", mapped.Log, tt.wantLog)
			}
		})
	}
}

func TestMapAuthErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "account inactive",
			err:        authdomain.ErrAccountInactiveError,
			wantStatus: http.StatusForbidden,
			wantCode:   authdomain.ErrAccountInactive,
		},
		{
			name:       "invalid credentials",
			err:        authdomain.ErrInvalidCredentialsError,
			wantStatus: http.StatusUnauthorized,
			wantCode:   authdomain.ErrInvalidCredentials,
		},
		{
			name:       "email taken",
			err:        authdomain.ErrEmailTakenError("test@example.com"),
			wantStatus: http.StatusConflict,
			wantCode:   authdomain.ErrEmailTaken,
		},
		{
			name:       "unknown defaults to invalid credentials",
			err:        errors.New("unexpected"),
			wantStatus: http.StatusUnauthorized,
			wantCode:   authdomain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapped := MapAuth(tt.err)
			if mapped.Status != tt.wantStatus {
				t.Fatalf("status = %d, want %d", mapped.Status, tt.wantStatus)
			}
			if mapped.Code != tt.wantCode {
				t.Fatalf("code = %q, want %q", mapped.Code, tt.wantCode)
			}
		})
	}
}

func TestMapUserManagementErrors(t *testing.T) {
	mapped := MapUserManagement(errors.New("unexpected"))
	if mapped.Status != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", mapped.Status, http.StatusInternalServerError)
	}
	if mapped.Code != pkgerrors.ErrInternalServer {
		t.Fatalf("code = %q, want %q", mapped.Code, pkgerrors.ErrInternalServer)
	}
}
