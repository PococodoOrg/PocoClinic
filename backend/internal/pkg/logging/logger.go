package logging

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Logger wraps slog.Logger to provide structured logging
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Add function name to log output
		AddSource: true,
	}

	// Use JSON handler in production, text handler in development
	var handler slog.Handler
	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithContext adds context values to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Add request ID if present
	if reqID, ok := ctx.Value("request_id").(string); ok {
		return &Logger{
			Logger: l.With("request_id", reqID),
		}
	}
	return l
}

// Error logs an error with additional context. Error text is sanitized for log injection.
func (l *Logger) Error(msg string, err error, args ...any) {
	if err != nil {
		newArgs := append([]any{"error", SanitizeString(err.Error())}, args...)
		l.Logger.Error(sanitizeLogMessage(msg), newArgs...)
	}
}

// Info logs an informational message. String values are sanitized for log injection.
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(sanitizeLogMessage(msg), sanitizeLogArgs(args...)...)
}

// RequestLogger creates a middleware for logging HTTP requests
func (l *Logger) RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a copy of the response writer that captures the status code
			rw := &responseWriter{ResponseWriter: w}

			// Process request
			next.ServeHTTP(rw, r)

			// Log request details (sanitize user-controlled header/path values).
			l.Info("HTTP Request",
				"method", r.Method,
				"path", SanitizeString(r.URL.Path),
				"status", rw.status,
				"duration", time.Since(start),
				"ip", SanitizeString(r.RemoteAddr),
				"user_agent", SanitizeString(r.UserAgent()),
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
