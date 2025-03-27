package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/features/patients/queries"
	"github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetPatientHandler is a mock implementation of the GetPatientHandler interface
type MockGetPatientHandler struct {
	mock.Mock
}

func (m *MockGetPatientHandler) Handle(ctx context.Context, query queries.GetPatientQuery) (*domain.Patient, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func TestGetPatient(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logger := logging.NewLogger()

	// Create test UUIDs
	testID := uuid.New()
	otherID := uuid.New()
	errorID := uuid.New()

	tests := []struct {
		name          string
		patientID     string
		setupMock     func(*MockGetPatientHandler)
		expectedCode  int
		expectedError string
	}{
		{
			name:      "Success",
			patientID: testID.String(),
			setupMock: func(m *MockGetPatientHandler) {
				m.On("Handle", mock.Anything, queries.GetPatientQuery{ID: testID.String()}).Return(
					&domain.Patient{
						ID:        testID,
						FirstName: "John",
						LastName:  "Doe",
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "Not Found",
			patientID: otherID.String(),
			setupMock: func(m *MockGetPatientHandler) {
				m.On("Handle", mock.Anything, queries.GetPatientQuery{ID: otherID.String()}).Return(
					nil, errors.NewAPIError(errors.ErrNotFound, "Patient not found"))
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "Patient not found",
		},
		{
			name:      "Internal Server Error",
			patientID: errorID.String(),
			setupMock: func(m *MockGetPatientHandler) {
				m.On("Handle", mock.Anything, queries.GetPatientQuery{ID: errorID.String()}).Return(
					nil, errors.NewAPIError(errors.ErrInternalServer, "Database error"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockHandler := new(MockGetPatientHandler)
			tt.setupMock(mockHandler)

			// Create handler with mock
			handler := NewPatientHandler(nil, nil, mockHandler, nil, logger)

			// Setup router
			router := gin.New()
			api := router.Group("/api")
			handler.RegisterRoutes(api)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/patients/"+tt.patientID, nil)

			// Serve request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// If we expect an error, verify the error message
			if tt.expectedError != "" {
				var response errors.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Message)
			}

			// Verify all mocked calls were made
			mockHandler.AssertExpectations(t)
		})
	}
}
