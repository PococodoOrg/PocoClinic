package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/patients/commands"
	"github.com/dksch/pococlinic/internal/features/patients/domain"
	patientinfra "github.com/dksch/pococlinic/internal/features/patients/infrastructure"
	"github.com/dksch/pococlinic/internal/features/patients/queries"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupExerciseHTTP(t *testing.T) (*gin.Engine, *domain.Patient, uuid.UUID) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	patientRepo := patientinfra.NewMemoryRepository()
	patient := domain.NewPatient("HTTP", "Exercise", time.Date(1975, 5, 5, 0, 0, 0, 0, time.UTC), domain.GenderMale)
	require.NoError(t, patientRepo.Create(context.Background(), patient))

	repo := patientinfra.NewMemoryExerciseLogRepository(patientRepo)
	handler := NewExerciseLogHandler(
		commands.NewExerciseLogCommandHandler(repo, patientRepo),
		queries.NewExerciseLogQueryHandler(repo),
		logging.NewLogger(),
		nil,
	)

	authorID := uuid.New()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", authorID.String())
		c.Next()
	})
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api, nil)
	return router, patient, authorID
}

func TestExerciseLogHTTP_RejectsMalformedJSONAndIDs(t *testing.T) {
	router, patient, _ := setupExerciseHTTP(t)
	base := "/api/v1/patients/" + patient.ID.String() + "/exercise-plans"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, base, bytes.NewBufferString(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.GreaterOrEqual(t, w.Code, 400)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, base+"/not-a-uuid", bytes.NewBufferString(`{"name":"x","status":"active"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	planID := uuid.New().String()
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, base+"/"+planID+"/entries/not-a-uuid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, base+"/"+planID+"/entries/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestExerciseLogHTTP_CreateListAndLogSession(t *testing.T) {
	router, patient, _ := setupExerciseHTTP(t)

	createBody := `{"name":"Knee PT","description":"Week 1"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patients/"+patient.ID.String()+"/exercise-plans", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var plan domain.ExercisePlan
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &plan))
	assert.Equal(t, "Knee PT", plan.Name)
	assert.Equal(t, domain.ExercisePlanStatusActive, plan.Status)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/patients/"+patient.ID.String()+"/exercise-plans", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var plans []*domain.ExercisePlan
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &plans))
	require.Len(t, plans, 1)

	entryBody := `{"exerciseName":"Quad sets","sets":3,"reps":10,"resistance":"yellow band","difficulty":4}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodPost,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String()+"/entries",
		bytes.NewBufferString(entryBody),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var entry domain.ExerciseLogEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &entry))
	assert.Equal(t, "Quad sets", entry.ExerciseName)
	assert.Equal(t, 3, entry.Sets)
	assert.Equal(t, 10, entry.Reps)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodGet,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String()+"/entries",
		nil,
	)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var entries []*domain.ExerciseLogEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &entries))
	require.Len(t, entries, 1)
}

func TestExerciseLogHTTP_ValidationAndNotFound(t *testing.T) {
	router, patient, _ := setupExerciseHTTP(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patients/not-a-uuid/exercise-plans", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodPost,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans",
		bytes.NewBufferString(`{"name":""}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.GreaterOrEqual(t, w.Code, 400, "expected error for empty plan name, got %d body=%s", w.Code, w.Body.String())

	missingPlan := uuid.New().String()
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodGet,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+missingPlan+"/entries",
		nil,
	)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestExerciseLogHTTP_UpdateAndDeleteEntry(t *testing.T) {
	router, patient, _ := setupExerciseHTTP(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans",
		bytes.NewBufferString(`{"name":"Hip PT"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var plan domain.ExercisePlan
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &plan))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodPost,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String()+"/entries",
		bytes.NewBufferString(`{"exerciseName":"Clamshell","sets":2,"reps":12}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var entry domain.ExerciseLogEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &entry))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodPut,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String()+"/entries/"+entry.ID.String(),
		bytes.NewBufferString(`{"exerciseName":"Clamshell","sets":3,"reps":15,"resistance":"green band","difficulty":5}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var updated domain.ExerciseLogEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &updated))
	assert.Equal(t, 15, updated.Reps)
	assert.Equal(t, "green band", updated.Resistance)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodDelete,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String()+"/entries/"+entry.ID.String(),
		nil,
	)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodPost,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String()+"/entries",
		bytes.NewBufferString(`{"exerciseName":"Clamshell","sets":-2,"reps":1}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.GreaterOrEqual(t, w.Code, 400)
}

func TestExerciseLogHTTP_UpdateDeletePlan(t *testing.T) {
	router, patient, _ := setupExerciseHTTP(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans",
		bytes.NewBufferString(`{"name":"Temp plan"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var plan domain.ExercisePlan
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &plan))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodPut,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String(),
		bytes.NewBufferString(`{"name":"Temp plan","description":"","status":"archived"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var updated domain.ExercisePlan
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &updated))
	assert.Equal(t, domain.ExercisePlanStatusArchived, updated.Status)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		http.MethodDelete,
		"/api/v1/patients/"+patient.ID.String()+"/exercise-plans/"+plan.ID.String(),
		nil,
	)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
