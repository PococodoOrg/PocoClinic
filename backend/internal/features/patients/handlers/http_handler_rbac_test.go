package handlers

import (
	"bytes"
	"context"
	"fmt"
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
	"github.com/PococodoOrg/PocoClinic/internal/pkg/doccrypto"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type patientAuthFixture struct {
	router      *gin.Engine
	patient     *patientdomain.Patient
	cookieName  string
	tokenCfg    authdomain.TokenConfig
	userRepo    *authinfra.MemoryUserRepository
	sessionRepo *authinfra.MemorySessionRepository
}

func setupPatientRBACHTTP(t *testing.T) *patientAuthFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)

	patientRepo := patientinfra.NewMemoryRepository()
	patient := patientdomain.NewPatient("RBAC", "Patient", time.Date(1988, 3, 3, 0, 0, 0, 0, time.UTC), patientdomain.GenderFemale)
	require.NoError(t, patientRepo.Create(context.Background(), patient))

	noteRepo := patientinfra.NewMemoryNoteRepository(patientRepo)
	docRepo := patientinfra.NewMemoryDocumentRepository(patientRepo)
	cipher, err := doccrypto.NewCipher(bytes.Repeat([]byte{9}, doccrypto.KeySize))
	require.NoError(t, err)
	legacyFiles := patientinfra.NewMemoryFileStorage()

	userRepo := authinfra.NewMemoryUserRepository()
	sessionRepo := authinfra.NewMemorySessionRepository()
	tokenCfg := authdomain.TokenConfig{
		AccessTokenSecret:    []byte("patient-rbac-access"),
		RefreshTokenSecret:   []byte("patient-rbac-refresh"),
		AccessTokenTTL:       15 * time.Minute,
		RefreshTokenTTL:      24 * time.Hour,
		SessionInactivityTTL: 15 * time.Minute,
		Issuer:               "pococlinic-patient-rbac",
	}
	authCfg := config.AuthConfig{
		AccessTokenCookieName:  "poco_access_token",
		RefreshTokenCookieName: "poco_refresh_token",
	}

	settingsRepo := patientinfra.NewMemorySettingsRepository()
	getFieldRequirementsHandler := queries.NewGetFieldRequirementsHandler(settingsRepo)

	handler := NewPatientHandler(
		commands.NewCreatePatientHandler(patientRepo, settingsRepo),
		queries.NewGetPatientsHandler(patientRepo),
		queries.NewGetPatientHandler(patientRepo),
		commands.NewUpdatePatientHandler(patientRepo, settingsRepo),
		commands.NewDeletePatientHandler(patientRepo),
		commands.NewCreateNoteHandler(noteRepo, patientRepo),
		commands.NewUpdateNoteHandler(noteRepo),
		commands.NewDeleteNoteHandler(noteRepo),
		queries.NewListNotesHandler(noteRepo),
		commands.NewUploadDocumentHandler(docRepo, patientRepo, cipher),
		commands.NewDeleteDocumentHandler(docRepo, legacyFiles),
		queries.NewListDocumentsHandler(docRepo),
		queries.NewGetDocumentHandler(docRepo),
		commands.NewOpenDocumentContentHandler(docRepo, cipher, legacyFiles),
		getFieldRequirementsHandler,
		logging.NewLogger(),
		nil,
	)
	authMW := authmiddleware.NewAuthMiddleware(tokenCfg, authCfg, userRepo, sessionRepo, logging.NewLogger())

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api, authMW)

	return &patientAuthFixture{
		router:      router,
		patient:     patient,
		cookieName:  authCfg.AccessTokenCookieName,
		tokenCfg:    tokenCfg,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (f *patientAuthFixture) loginAs(t *testing.T, role authdomain.Role, mustChangePIN bool) *http.Cookie {
	t.Helper()
	email := fmt.Sprintf("%s-%d@example.com", role, time.Now().UnixNano())
	user := authdomain.NewUser(email, string(role)+" User", role)
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

func (f *patientAuthFixture) do(method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
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

func TestPatientHTTP_RequiresAuth(t *testing.T) {
	f := setupPatientRBACHTTP(t)
	w := f.do(http.MethodGet, "/api/v1/patients", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientHTTP_StaffCanReadButNotCreate(t *testing.T) {
	f := setupPatientRBACHTTP(t)
	cookie := f.loginAs(t, authdomain.RoleStaff, false)

	w := f.do(http.MethodGet, "/api/v1/patients/"+f.patient.ID.String(), "", cookie)
	assert.Equal(t, http.StatusOK, w.Code)

	createBody := `{"firstName":"New","lastName":"Patient","dateOfBirth":"1990-01-01","gender":"female","email":"new@example.com","phoneNumber":"555-0100"}`
	w = f.do(http.MethodPost, "/api/v1/patients", createBody, cookie)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPatientHTTP_NurseCanCreatePatient(t *testing.T) {
	f := setupPatientRBACHTTP(t)
	cookie := f.loginAs(t, authdomain.RoleNurse, false)

	createBody := `{"firstName":"New","lastName":"Patient","dateOfBirth":"1990-01-01","gender":"female","email":"new@example.com","phoneNumber":"555-0100"}`
	w := f.do(http.MethodPost, "/api/v1/patients", createBody, cookie)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPatientHTTP_OnlyAdminDeletesPatient(t *testing.T) {
	f := setupPatientRBACHTTP(t)
	nurseCookie := f.loginAs(t, authdomain.RoleNurse, false)
	adminCookie := f.loginAs(t, authdomain.RoleAdmin, false)

	w := f.do(http.MethodDelete, "/api/v1/patients/"+f.patient.ID.String(), "", nurseCookie)
	assert.Equal(t, http.StatusForbidden, w.Code)

	w = f.do(http.MethodDelete, "/api/v1/patients/"+f.patient.ID.String(), "", adminCookie)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPatientHTTP_MustChangePINBlocksAccess(t *testing.T) {
	f := setupPatientRBACHTTP(t)
	cookie := f.loginAs(t, authdomain.RoleDoctor, true)
	w := f.do(http.MethodGet, "/api/v1/patients/"+f.patient.ID.String(), "", cookie)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
