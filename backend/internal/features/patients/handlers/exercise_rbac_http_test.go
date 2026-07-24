package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authdomain "github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	authinfra "github.com/PococodoOrg/PocoClinic/internal/features/auth/infrastructure"
	authmiddleware "github.com/PococodoOrg/PocoClinic/internal/features/auth/middleware"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/commands"
	patientdomain "github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	patientinfra "github.com/PococodoOrg/PocoClinic/internal/features/patients/infrastructure"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/queries"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exerciseAuthFixture struct {
	router     *gin.Engine
	patient    *patientdomain.Patient
	cookieName string
	tokenCfg   authdomain.TokenConfig
	userRepo   *authinfra.MemoryUserRepository
	sessionRepo *authinfra.MemorySessionRepository
}

func setupExerciseAuthHTTP(t *testing.T) *exerciseAuthFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)

	patientRepo := patientinfra.NewMemoryRepository()
	patient := patientdomain.NewPatient("RBAC", "Exercise", time.Date(1982, 6, 6, 0, 0, 0, 0, time.UTC), patientdomain.GenderFemale)
	require.NoError(t, patientRepo.Create(context.Background(), patient))

	userRepo := authinfra.NewMemoryUserRepository()
	sessionRepo := authinfra.NewMemorySessionRepository()
	tokenCfg := authdomain.TokenConfig{
		AccessTokenSecret:    []byte("exercise-test-access-secret"),
		RefreshTokenSecret:   []byte("exercise-test-refresh-secret"),
		AccessTokenTTL:       15 * time.Minute,
		RefreshTokenTTL:      24 * time.Hour,
		SessionInactivityTTL: 15 * time.Minute,
		Issuer:               "pococlinic-exercise-test",
	}
	authCfg := config.AuthConfig{
		AccessTokenCookieName:  "poco_access_token",
		RefreshTokenCookieName: "poco_refresh_token",
		AccessTokenSecret:      string(tokenCfg.AccessTokenSecret),
		RefreshTokenSecret:     string(tokenCfg.RefreshTokenSecret),
	}

	exerciseRepo := patientinfra.NewMemoryExerciseLogRepository(patientRepo)
	handler := NewExerciseLogHandler(
		commands.NewExerciseLogCommandHandler(exerciseRepo, patientRepo),
		queries.NewExerciseLogQueryHandler(exerciseRepo),
		logging.NewLogger(),
		nil,
	)
	authMW := authmiddleware.NewAuthMiddleware(tokenCfg, authCfg, userRepo, sessionRepo, logging.NewLogger())

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api, authMW)

	return &exerciseAuthFixture{
		router:      router,
		patient:     patient,
		cookieName:  authCfg.AccessTokenCookieName,
		tokenCfg:    tokenCfg,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (f *exerciseAuthFixture) loginAs(t *testing.T, role authdomain.Role, mustChangePIN bool) *http.Cookie {
	t.Helper()
	user := authdomain.NewUser(string(role)+"@example.com", string(role)+" User", role)
	user.MustChangePIN = mustChangePIN
	user.IsActive = true
	pin, err := authdomain.GeneratePINCredential("1234")
	require.NoError(t, err)
	user.SetPINCredential(pin)
	key, keyCred, err := authdomain.GenerateKey()
	require.NoError(t, err)
	user.SetKeyCredential(keyCred, key)
	require.NoError(t, f.userRepo.Create(context.Background(), user))

	session := authdomain.NewSession(user.ID, "test-agent", "127.0.0.1", time.Now().Add(time.Hour))
	access, _, err := session.GenerateTokens(user, f.tokenCfg)
	require.NoError(t, err)
	require.NoError(t, f.sessionRepo.Create(context.Background(), session))

	return &http.Cookie{Name: f.cookieName, Value: access}
}

func (f *exerciseAuthFixture) do(method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	f.router.ServeHTTP(w, req)
	return w
}

func TestExerciseLogHTTP_RequiresAuth(t *testing.T) {
	f := setupExerciseAuthHTTP(t)
	path := "/api/v1/patients/" + f.patient.ID.String() + "/exercise-plans"

	w := f.do(http.MethodGet, path, "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = f.do(http.MethodPost, path, `{"name":"Knee"}`, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestExerciseLogHTTP_StaffCanReadButNotWrite(t *testing.T) {
	f := setupExerciseAuthHTTP(t)
	staff := f.loginAs(t, authdomain.RoleStaff, false)
	path := "/api/v1/patients/" + f.patient.ID.String() + "/exercise-plans"

	w := f.do(http.MethodGet, path, "", staff)
	assert.Equal(t, http.StatusOK, w.Code)

	w = f.do(http.MethodPost, path, `{"name":"Knee PT"}`, staff)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestExerciseLogHTTP_NurseCanWrite(t *testing.T) {
	f := setupExerciseAuthHTTP(t)
	nurse := f.loginAs(t, authdomain.RoleNurse, false)
	path := "/api/v1/patients/" + f.patient.ID.String() + "/exercise-plans"

	w := f.do(http.MethodPost, path, `{"name":"Shoulder PT","description":"ROM"}`, nurse)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestExerciseLogHTTP_MustChangePINBlocksAccess(t *testing.T) {
	f := setupExerciseAuthHTTP(t)
	cookie := f.loginAs(t, authdomain.RoleDoctor, true)
	path := "/api/v1/patients/" + f.patient.ID.String() + "/exercise-plans"

	w := f.do(http.MethodGet, path, "", cookie)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestExerciseLogHTTP_InactiveUserRejected(t *testing.T) {
	f := setupExerciseAuthHTTP(t)
	user := authdomain.NewUser("inactive@example.com", "Inactive", authdomain.RoleDoctor)
	user.MustChangePIN = false
	user.IsActive = false
	pin, err := authdomain.GeneratePINCredential("1234")
	require.NoError(t, err)
	user.SetPINCredential(pin)
	key, keyCred, err := authdomain.GenerateKey()
	require.NoError(t, err)
	user.SetKeyCredential(keyCred, key)
	require.NoError(t, f.userRepo.Create(context.Background(), user))

	session := authdomain.NewSession(user.ID, "test-agent", "127.0.0.1", time.Now().Add(time.Hour))
	access, _, err := session.GenerateTokens(user, f.tokenCfg)
	require.NoError(t, err)
	require.NoError(t, f.sessionRepo.Create(context.Background(), session))

	w := f.do(http.MethodGet, "/api/v1/patients/"+f.patient.ID.String()+"/exercise-plans", "", &http.Cookie{
		Name:  f.cookieName,
		Value: access,
	})
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "ACCOUNT_INACTIVE")
}
