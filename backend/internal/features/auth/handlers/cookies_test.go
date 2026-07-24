package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

func testAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		AccessTokenCookieName:  "poco_access_token",
		RefreshTokenCookieName: "poco_refresh_token",
		AccessTokenTTL:         15 * time.Minute,
		RefreshTokenTTL:        24 * time.Hour,
		CookieSecure:           false,
	}
}

func TestSetAuthCookies_EmitsHttpOnlyPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	cfg := testAuthConfig()

	setAuthCookies(c, cfg, "access-token-value", "refresh-token-value")

	setCookies := w.Header()["Set-Cookie"]
	if len(setCookies) != 2 {
		t.Fatalf("expected 2 cookies, got %v", setCookies)
	}

	joined := strings.Join(setCookies, "\n")
	if !strings.Contains(joined, cfg.AccessTokenCookieName+"=access-token-value") {
		t.Fatalf("missing access cookie: %s", joined)
	}
	if !strings.Contains(joined, "Path=/api/v1;") || !strings.Contains(joined, "HttpOnly") {
		t.Fatalf("access cookie missing path/httponly: %s", joined)
	}
	if !strings.Contains(joined, cfg.RefreshTokenCookieName+"=refresh-token-value") {
		t.Fatalf("missing refresh cookie: %s", joined)
	}
	if !strings.Contains(joined, "Path=/api/v1/auth") {
		t.Fatalf("refresh cookie missing auth path: %s", joined)
	}
}

func TestClearAuthCookies_ExpiresBoth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	cfg := testAuthConfig()

	clearAuthCookies(c, cfg)

	setCookies := w.Header()["Set-Cookie"]
	if len(setCookies) != 2 {
		t.Fatalf("expected 2 cleared cookies, got %v", setCookies)
	}
	for _, cookie := range setCookies {
		if !strings.Contains(cookie, "Max-Age=0") {
			t.Fatalf("expected expired cookie, got %q", cookie)
		}
	}
}

func TestReadRefreshTokenCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testAuthConfig()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/api/v1/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: cfg.RefreshTokenCookieName, Value: "abc123"})
	c.Request = req

	if got := readRefreshTokenCookie(c, cfg); got != "abc123" {
		t.Fatalf("expected refresh token abc123, got %q", got)
	}

	emptyW := httptest.NewRecorder()
	emptyC, _ := gin.CreateTestContext(emptyW)
	emptyC.Request = httptest.NewRequest("GET", "/api/v1/auth/refresh", nil)
	if got := readRefreshTokenCookie(emptyC, cfg); got != "" {
		t.Fatalf("expected empty refresh token, got %q", got)
	}
}
