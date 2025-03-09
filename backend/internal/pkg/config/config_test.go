package config

import (
	"os"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save current env vars
	oldPort := os.Getenv("SERVER_PORT")
	oldHost := os.Getenv("SERVER_HOST")
	oldOrigin := os.Getenv("ALLOWED_ORIGIN")
	oldRPS := os.Getenv("RATE_LIMIT_RPS")
	oldBurst := os.Getenv("RATE_LIMIT_BURST")

	// Restore env vars after test
	defer func() {
		os.Setenv("SERVER_PORT", oldPort)
		os.Setenv("SERVER_HOST", oldHost)
		os.Setenv("ALLOWED_ORIGIN", oldOrigin)
		os.Setenv("RATE_LIMIT_RPS", oldRPS)
		os.Setenv("RATE_LIMIT_BURST", oldBurst)
	}()

	tests := []struct {
		name      string
		envVars   map[string]string
		wantError bool
		validate  func(*testing.T, *Config)
	}{
		{
			name:      "Default configuration",
			envVars:   map[string]string{},
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 8080, cfg.Server.Port)
				assert.Equal(t, "localhost", cfg.Server.Host)
				assert.Equal(t, []string{"http://localhost:3000"}, cfg.Security.AllowedOrigins)
				assert.Equal(t, 10, cfg.Security.RateLimit.RequestsPerSecond)
				assert.Equal(t, 20, cfg.Security.RateLimit.BurstSize)
			},
		},
		{
			name: "Custom configuration",
			envVars: map[string]string{
				"SERVER_PORT":      "9000",
				"SERVER_HOST":      "0.0.0.0",
				"ALLOWED_ORIGIN":   "https://example.com",
				"RATE_LIMIT_RPS":   "100",
				"RATE_LIMIT_BURST": "50",
			},
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 9000, cfg.Server.Port)
				assert.Equal(t, "0.0.0.0", cfg.Server.Host)
				assert.Equal(t, []string{"https://example.com"}, cfg.Security.AllowedOrigins)
				assert.Equal(t, 100, cfg.Security.RateLimit.RequestsPerSecond)
				assert.Equal(t, 50, cfg.Security.RateLimit.BurstSize)
			},
		},
		{
			name: "Invalid port",
			envVars: map[string]string{
				"SERVER_PORT": "invalid",
			},
			wantError: true,
		},
		{
			name: "Invalid rate limit",
			envVars: map[string]string{
				"RATE_LIMIT_RPS": "invalid",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg, err := LoadConfig()
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			tt.validate(t, cfg)
		})
	}
}

func TestConfigureCORS(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		validateConfig func(*testing.T, cors.Config)
	}{
		{
			name: "Default CORS configuration",
			config: &Config{
				Security: SecurityConfig{
					AllowedOrigins: []string{"http://localhost:3000"},
				},
			},
			validateConfig: func(t *testing.T, c cors.Config) {
				assert.Equal(t, []string{"http://localhost:3000"}, c.AllowOrigins)
				assert.ElementsMatch(t, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, c.AllowMethods)
				assert.ElementsMatch(t, []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, c.AllowHeaders)
				assert.True(t, c.AllowCredentials)
				assert.Equal(t, 12*time.Hour, c.MaxAge)
			},
		},
		{
			name: "Multiple allowed origins",
			config: &Config{
				Security: SecurityConfig{
					AllowedOrigins: []string{"http://localhost:3000", "https://example.com"},
				},
			},
			validateConfig: func(t *testing.T, c cors.Config) {
				assert.Equal(t, []string{"http://localhost:3000", "https://example.com"}, c.AllowOrigins)
				assert.ElementsMatch(t, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, c.AllowMethods)
				assert.ElementsMatch(t, []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, c.AllowHeaders)
				assert.True(t, c.AllowCredentials)
				assert.Equal(t, 12*time.Hour, c.MaxAge)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corsConfig := tt.config.ConfigureCORS()
			tt.validateConfig(t, corsConfig)
		})
	}
}
