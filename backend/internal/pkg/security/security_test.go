package security_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/dksch/pococlinic/internal/pkg/config"
	"github.com/dksch/pococlinic/internal/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestSecurityHeaders(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Content-Security-Policy":   "default-src 'self'",
	}

	for header, expected := range expectedHeaders {
		assert.Equal(t, expected, w.Header().Get(header), fmt.Sprintf("Expected header %s to be %s", header, expected))
	}
}

func TestRateLimiter(t *testing.T) {
	router := setupTestRouter()
	// Configure for 2 requests per second with burst of 3
	limiter := middleware.NewIPRateLimiter(rate.Limit(2), 3)
	router.Use(middleware.RateLimiterMiddleware(limiter))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		name          string
		requests      int
		delayBetween  time.Duration
		expectedCodes []int
		testIP        string // Each test gets its own IP
	}{
		{
			name:         "Initial burst allows 3 requests",
			requests:     4,
			delayBetween: 0,
			expectedCodes: []int{
				http.StatusOK,              // First request - allowed (1/3 burst)
				http.StatusOK,              // Second request - allowed (2/3 burst)
				http.StatusOK,              // Third request - allowed (3/3 burst)
				http.StatusTooManyRequests, // Fourth request - denied (burst exceeded)
			},
			testIP: "192.0.2.1",
		},
		{
			name:         "Requests with delay are allowed",
			requests:     2,
			delayBetween: time.Second, // Wait 1 second between requests
			expectedCodes: []int{
				http.StatusOK, // First request - allowed
				http.StatusOK, // Second request after 1s - allowed
			},
			testIP: "192.0.2.2", // Different IP for this test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset rate limiter for each test
			limiter = middleware.NewIPRateLimiter(rate.Limit(2), 3)

			for i := 0; i < tt.requests; i++ {
				if tt.delayBetween > 0 && i > 0 {
					time.Sleep(tt.delayBetween)
				}

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = tt.testIP + ":1234" // Use the test-specific IP
				router.ServeHTTP(w, req)

				if i < len(tt.expectedCodes) {
					assert.Equal(t,
						tt.expectedCodes[i],
						w.Code,
						fmt.Sprintf("Request %d in test '%s': expected status %d but got %d",
							i+1,
							tt.name,
							tt.expectedCodes[i],
							w.Code,
						),
					)
				}
			}
		})
	}
}

func TestConcurrentRateLimiter(t *testing.T) {
	router := setupTestRouter()
	limiter := middleware.NewIPRateLimiter(rate.Limit(10), 1)
	router.Use(middleware.RateLimiterMiddleware(limiter))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	var wg sync.WaitGroup
	successCount := 0
	rateLimitCount := 0
	mu := sync.Mutex{}

	// Launch 20 concurrent requests
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.0.2.1:1234"
			router.ServeHTTP(w, req)

			mu.Lock()
			defer mu.Unlock()
			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				rateLimitCount++
			}
		}()
	}

	wg.Wait()

	// We expect some requests to succeed and some to be rate limited
	assert.True(t, successCount > 0, "Expected some requests to succeed")
	assert.True(t, rateLimitCount > 0, "Expected some requests to be rate limited")
	assert.Equal(t, 20, successCount+rateLimitCount, "Total requests should equal successes plus rate limits")
}

func TestCORSConfiguration(t *testing.T) {
	tests := []struct {
		name              string
		origin            string
		method            string
		headers           map[string]string
		expectedStatus    int
		shouldAllowOrigin bool
		expectedHeaders   map[string]string
	}{
		{
			name:              "Allowed origin",
			origin:            "http://localhost:3000",
			method:            "GET",
			expectedStatus:    http.StatusOK,
			shouldAllowOrigin: true,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:              "Disallowed origin",
			origin:            "http://evil.com",
			method:            "GET",
			expectedStatus:    http.StatusForbidden,
			shouldAllowOrigin: false,
		},
		{
			name:              "Options request",
			origin:            "http://localhost:3000",
			method:            "OPTIONS",
			expectedStatus:    http.StatusNoContent,
			shouldAllowOrigin: true,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Methods":     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
				"Access-Control-Allow-Headers":     "Origin,Content-Length,Content-Type,Authorization",
			},
		},
		{
			name:   "Custom headers in request",
			origin: "http://localhost:3000",
			method: "OPTIONS",
			headers: map[string]string{
				"Access-Control-Request-Headers": "X-Custom-Header",
			},
			expectedStatus:    http.StatusNoContent,
			shouldAllowOrigin: true,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Headers": "Origin,Content-Length,Content-Type,Authorization",
			},
		},
		{
			name:   "Custom method in request",
			origin: "http://localhost:3000",
			method: "OPTIONS",
			headers: map[string]string{
				"Access-Control-Request-Method": "PATCH",
			},
			expectedStatus:    http.StatusNoContent,
			shouldAllowOrigin: true,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods": "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			cfg := &config.Config{
				Security: config.SecurityConfig{
					AllowedOrigins: []string{"http://localhost:3000"},
				},
			}

			router.Use(cors.New(cfg.ConfigureCORS()))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)

			// Set any custom headers
			if tt.headers != nil {
				for key, value := range tt.headers {
					req.Header.Set(key, value)
				}
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldAllowOrigin {
				assert.Equal(t, tt.origin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}

			// Check expected headers
			if tt.expectedHeaders != nil {
				for header, expected := range tt.expectedHeaders {
					assert.Equal(t, expected, w.Header().Get(header),
						fmt.Sprintf("Expected header %s to be %s", header, expected))
				}
			}
		})
	}
}

func TestRateLimiterCleanup(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(rate.Limit(1), 1)

	// Get limiters for test IPs
	limiter.GetLimiter("192.0.2.1")
	limiter.GetLimiter("192.0.2.2")

	// Wait a bit to simulate time passing
	time.Sleep(100 * time.Millisecond)

	// Try to get limiter for old IP
	newLimiter := limiter.GetLimiter("192.0.2.1")
	assert.NotNil(t, newLimiter, "Should create new limiter for cleaned up IP")
}

func TestPanicRecovery(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.Recovery())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}
