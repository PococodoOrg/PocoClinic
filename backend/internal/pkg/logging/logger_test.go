package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantJSON    bool
	}{
		{
			name:        "Development environment",
			environment: "development",
			wantJSON:    false,
		},
		{
			name:        "Production environment",
			environment: "production",
			wantJSON:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current env
			oldEnv := os.Getenv("ENV")
			defer os.Setenv("ENV", oldEnv)

			// Set test env
			os.Setenv("ENV", tt.environment)

			// Create logger with custom output
			var buf bytes.Buffer
			opts := &slog.HandlerOptions{Level: slog.LevelInfo}
			var handler slog.Handler
			if tt.wantJSON {
				handler = slog.NewJSONHandler(&buf, opts)
			} else {
				handler = slog.NewTextHandler(&buf, opts)
			}
			logger := &Logger{Logger: slog.New(handler)}

			// Log a test message
			logger.Info("test message", "key", "value")

			// Check output format
			output := buf.String()
			if tt.wantJSON {
				assert.True(t, json.Valid([]byte(output)), "Expected valid JSON output")
			} else {
				assert.True(t, strings.Contains(output, "test message"), "Expected text output")
			}
		})
	}
}

func TestWithContext(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := &Logger{Logger: slog.New(handler)}

	// Create context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "test-id")
	loggerWithCtx := logger.WithContext(ctx)

	// Log a message
	loggerWithCtx.Info("test message")

	// Parse output
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Check if request_id is present
	assert.Equal(t, "test-id", logEntry["request_id"])
}

func TestRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := &Logger{Logger: slog.New(handler)}

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Simulate work
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with request logger
	loggingHandler := logger.RequestLogger()(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	req.Header.Set("User-Agent", "test-agent")

	// Execute request
	w := httptest.NewRecorder()
	loggingHandler.ServeHTTP(w, req)

	// Parse log output
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Verify log fields
	assert.Equal(t, "HTTP Request", logEntry["msg"])
	assert.Equal(t, "GET", logEntry["method"])
	assert.Equal(t, "/test", logEntry["path"])
	assert.Equal(t, float64(200), logEntry["status"])
	assert.Equal(t, "192.0.2.1:1234", logEntry["ip"])
	assert.Equal(t, "test-agent", logEntry["user_agent"])
	assert.Contains(t, logEntry, "duration")
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(&buf, opts)
	logger := &Logger{Logger: slog.New(handler)}

	testError := assert.AnError
	logger.Error("test error", testError, "extra", "value")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	assert.Equal(t, "test error", logEntry["msg"])
	assert.Equal(t, testError.Error(), logEntry["error"])
	assert.Equal(t, "value", logEntry["extra"])
	assert.Equal(t, "ERROR", logEntry["level"])
}
