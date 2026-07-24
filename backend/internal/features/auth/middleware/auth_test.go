package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	authinfra "github.com/dksch/pococlinic/internal/features/auth/infrastructure"
	"github.com/dksch/pococlinic/internal/pkg/config"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authTestFixture struct {
	router     *gin.Engine
	cookieName string
	tokenCfg   domain.TokenConfig
	userRepo   *authinfra.MemoryUserRepository
	sessionRepo *authinfra.MemorySessionRepository
	authMW     *AuthMiddleware
}

func setupAuthMiddlewareTest(t *testing.T) *authTestFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)

	userRepo := authinfra.NewMemoryUserRepository()
	sessionRepo := authinfra.NewMemorySessionRepository()
	tokenCfg := domain.TokenConfig{
		AccessTokenSecret:    []byte("middleware-test-access"),
		RefreshTokenSecret:   []byte("middleware-test-refresh"),
		AccessTokenTTL:       15 * time.Minute,
		RefreshTokenTTL:      24 * time.Hour,
		SessionInactivityTTL: 15 * time.Minute,
		Issuer:               "pococlinic-middleware-test",
	}
	authCfg := config.AuthConfig{
		AccessTokenCookieName:  "poco_access_token",
		RefreshTokenCookieName: "poco_refresh_token",
	}
	authMW := NewAuthMiddleware(tokenCfg, authCfg, userRepo, sessionRepo, logging.NewLogger())

	router := gin.New()
	api := router.Group("/api/v1")
	api.GET("/protected", authMW.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"userID": c.GetString("userID")})
	})
	api.GET("/admin-only", authMW.RequireAuth(), authMW.RequireRole(domain.RoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	api.GET("/clinical", authMW.RequireAuth(), authMW.RequirePINChanged(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return &authTestFixture{
		router:      router,
		cookieName:  authCfg.AccessTokenCookieName,
		tokenCfg:    tokenCfg,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		authMW:      authMW,
	}
}

func (f *authTestFixture) issueCookie(t *testing.T, role domain.Role, mustChangePIN bool) *http.Cookie {
	t.Helper()
	email := fmt.Sprintf("%s-%d@example.com", role, time.Now().UnixNano())
	user := domain.NewUser(email, string(role), role)
	user.MustChangePIN = mustChangePIN
	user.IsActive = true
	pin, err := domain.GeneratePINCredential("1234")
	require.NoError(t, err)
	user.SetPINCredential(pin)
	key, keyCred, err := domain.GenerateKey()
	require.NoError(t, err)
	user.SetKeyCredential(keyCred, key)
	require.NoError(t, f.userRepo.Create(context.Background(), user))

	session := domain.NewSession(user.ID, "test", "127.0.0.1", time.Now().Add(time.Hour))
	access, _, err := session.GenerateTokens(user, f.tokenCfg)
	require.NoError(t, err)
	require.NoError(t, f.sessionRepo.Create(context.Background(), session))
	return &http.Cookie{Name: f.cookieName, Value: access}
}

func (f *authTestFixture) get(path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	f.router.ServeHTTP(w, req)
	return w
}

func TestRequireAuth_MissingCookie(t *testing.T) {
	f := setupAuthMiddlewareTest(t)
	w := f.get("/api/v1/protected", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireAuth_ValidSession(t *testing.T) {
	f := setupAuthMiddlewareTest(t)
	cookie := f.issueCookie(t, domain.RoleStaff, false)
	w := f.get("/api/v1/protected", cookie)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_AdminOnly(t *testing.T) {
	f := setupAuthMiddlewareTest(t)

	staffCookie := f.issueCookie(t, domain.RoleStaff, false)
	w := f.get("/api/v1/admin-only", staffCookie)
	assert.Equal(t, http.StatusForbidden, w.Code)

	adminCookie := f.issueCookie(t, domain.RoleAdmin, false)
	w = f.get("/api/v1/admin-only", adminCookie)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequirePINChanged_BlocksDefaultPIN(t *testing.T) {
	f := setupAuthMiddlewareTest(t)
	cookie := f.issueCookie(t, domain.RoleDoctor, true)
	w := f.get("/api/v1/clinical", cookie)
	assert.Equal(t, http.StatusForbidden, w.Code)

	cookie = f.issueCookie(t, domain.RoleDoctor, false)
	w = f.get("/api/v1/clinical", cookie)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireAuth_InactiveUser(t *testing.T) {
	f := setupAuthMiddlewareTest(t)
	user := domain.NewUser("inactive@example.com", "Inactive", domain.RoleStaff)
	user.IsActive = false
	pin, err := domain.GeneratePINCredential("1234")
	require.NoError(t, err)
	user.SetPINCredential(pin)
	key, keyCred, err := domain.GenerateKey()
	require.NoError(t, err)
	user.SetKeyCredential(keyCred, key)
	require.NoError(t, f.userRepo.Create(context.Background(), user))

	session := domain.NewSession(user.ID, "test", "127.0.0.1", time.Now().Add(time.Hour))
	access, _, err := session.GenerateTokens(user, f.tokenCfg)
	require.NoError(t, err)
	require.NoError(t, f.sessionRepo.Create(context.Background(), session))

	w := f.get("/api/v1/protected", &http.Cookie{Name: f.cookieName, Value: access})
	assert.Equal(t, http.StatusForbidden, w.Code)
}
