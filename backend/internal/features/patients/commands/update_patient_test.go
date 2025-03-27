package commands

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPatientRepository is a mock implementation of PatientRepository
type MockPatientRepository struct {
	mock.Mock
}

func (m *MockPatientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	args := m.Called(ctx, patient)
	return args.Error(0)
}

func (m *MockPatientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	args := m.Called(ctx, patient)
	return args.Error(0)
}

func (m *MockPatientRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPatientRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if patient, ok := args.Get(0).(*domain.Patient); ok {
		return patient, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPatientRepository) List(ctx context.Context) ([]*domain.Patient, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Patient), args.Error(1)
}

func (m *MockPatientRepository) ListPaginated(ctx context.Context, page, pageSize int, search string) ([]*domain.Patient, int64, error) {
	args := m.Called(ctx, page, pageSize, search)
	if args.Get(0) == nil {
		return nil, 0, args.Error(1)
	}
	return args.Get(0).([]*domain.Patient), args.Get(1).(int64), args.Error(2)
}

func (m *MockPatientRepository) GetPatientByID(ctx context.Context, id string) (*domain.Patient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func TestUpdatePatientHandler_Handle(t *testing.T) {
	// Create test data
	existingPatient := &domain.Patient{
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: domain.Date(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
		Gender:      domain.GenderMale,
		Email:       "john@example.com",
		PhoneNumber: "1234567890",
		Height:      180,
		Weight:      80,
	}

	updateCmd := UpdatePatientCommand{
		ID:          existingPatient.ID.String(),
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		Email:       "john@example.com",
		PhoneNumber: "1234567890",
		Height:      &[]float64{180}[0],
		Weight:      &[]float64{80}[0],
	}

	tests := []struct {
		name          string
		setupMock     func(*MockPatientRepository)
		cmd           UpdatePatientCommand
		expectedError error
		checkAPIError bool
	}{
		{
			name: "successful update",
			setupMock: func(mockRepo *MockPatientRepository) {
				mockRepo.On("GetByID", mock.Anything, existingPatient.ID.String()).Return(existingPatient, nil)
				mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			cmd:           updateCmd,
			expectedError: nil,
			checkAPIError: false,
		},
		{
			name: "patient not found",
			setupMock: func(mockRepo *MockPatientRepository) {
				// Return patient=nil but no error - this will trigger the nil check in the handler
				mockRepo.On("GetByID", mock.Anything, existingPatient.ID.String()).Return(nil, nil)
			},
			cmd:           updateCmd,
			expectedError: errors.NewAPIError(errors.ErrNotFound, "Patient not found"),
			checkAPIError: true,
		},
		{
			name: "error from repository",
			setupMock: func(mockRepo *MockPatientRepository) {
				repoErr := fmt.Errorf("database error")
				mockRepo.On("GetByID", mock.Anything, existingPatient.ID.String()).Return(nil, repoErr)
			},
			cmd:           updateCmd,
			expectedError: fmt.Errorf("database error"),
			checkAPIError: false,
		},
		{
			name: "invalid date format",
			setupMock: func(mockRepo *MockPatientRepository) {
				mockRepo.On("GetByID", mock.Anything, existingPatient.ID.String()).Return(existingPatient, nil)
			},
			cmd: UpdatePatientCommand{
				ID:          existingPatient.ID.String(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: "invalid-date",
				Gender:      "male",
				Email:       "john@example.com",
				PhoneNumber: "1234567890",
			},
			expectedError: errors.NewAPIError(errors.ErrValidation, "Invalid date of birth format"),
			checkAPIError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh mock repository for each test
			mockRepo := new(MockPatientRepository)

			// Setup mock expectations
			tt.setupMock(mockRepo)

			// Create the handler
			handler := NewUpdatePatientHandler(mockRepo)

			// Execute the handler
			patient, err := handler.Handle(context.Background(), tt.cmd)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, patient)

				if tt.checkAPIError {
					apiErr, ok := err.(*errors.APIError)
					if assert.True(t, ok, "Expected an APIError") {
						expectedAPIErr := tt.expectedError.(*errors.APIError)
						assert.Equal(t, expectedAPIErr.Code, apiErr.Code)
						assert.Equal(t, expectedAPIErr.Message, apiErr.Message)
					}
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, patient)
			assert.Equal(t, tt.cmd.FirstName, patient.FirstName)
			assert.Equal(t, tt.cmd.LastName, patient.LastName)
			assert.Equal(t, tt.cmd.Email, patient.Email)
			assert.Equal(t, tt.cmd.PhoneNumber, patient.PhoneNumber)
			if tt.cmd.Height != nil {
				assert.Equal(t, *tt.cmd.Height, patient.Height)
			}
			if tt.cmd.Weight != nil {
				assert.Equal(t, *tt.cmd.Weight, patient.Weight)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
