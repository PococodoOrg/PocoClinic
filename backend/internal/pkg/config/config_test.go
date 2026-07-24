package config

import (
	"os"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	envKeys := []string{
		"SERVER_PORT",
		"SERVER_HOST",
		"ALLOWED_ORIGIN",
		"RATE_LIMIT_RPS",
		"RATE_LIMIT_BURST",
		"ENV",
		"JWT_ACCESS_SECRET",
		"JWT_REFRESH_SECRET",
		"DATABASE_URL",
		"RUN_MIGRATIONS",
		"DOCUMENT_ENCRYPTION_KEY",
	}

	saved := make(map[string]string, len(envKeys))
	for _, key := range envKeys {
		saved[key] = os.Getenv(key)
	}
	defer func() {
		for _, key := range envKeys {
			if saved[key] == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, saved[key])
			}
		}
	}()

	clearEnv := func() {
		for _, key := range envKeys {
			os.Unsetenv(key)
		}
	}

	strongAccess := "prod-access-secret-with-enough-length-here"
	strongRefresh := "prod-refresh-secret-with-enough-length-here"
	prodDB := "./data/pococlinic-test.db"
	// Non-default 32-byte key for production tests (32 zero bytes, base64).
	prodDocKey := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

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
				assert.True(t, cfg.Database.RunMigrations)
				assert.Equal(t, []byte(defaultDocumentEncryptionKeyMaterial), cfg.Storage.DocumentEncryptionKey)
				assert.Len(t, cfg.Storage.DocumentEncryptionKey, 32)
			},
		},
		{
			name: "Production skips migrations on app start",
			envVars: map[string]string{
				"ENV":                      "production",
				"JWT_ACCESS_SECRET":        strongAccess,
				"JWT_REFRESH_SECRET":       strongRefresh,
				"DATABASE_URL":             prodDB,
				"DOCUMENT_ENCRYPTION_KEY":  prodDocKey,
			},
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.False(t, cfg.Database.RunMigrations)
				assert.Equal(t, prodDB, cfg.Database.URL)
			},
		},
		{
			name: "RUN_MIGRATIONS override",
			envVars: map[string]string{
				"ENV":                      "production",
				"RUN_MIGRATIONS":           "true",
				"JWT_ACCESS_SECRET":        strongAccess,
				"JWT_REFRESH_SECRET":       strongRefresh,
				"DATABASE_URL":             prodDB,
				"DOCUMENT_ENCRYPTION_KEY":  prodDocKey,
			},
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.Database.RunMigrations)
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
			name: "Production rejects default JWT secrets",
			envVars: map[string]string{
				"ENV":          "production",
				"DATABASE_URL": prodDB,
			},
			wantError: true,
		},
		{
			name: "Production requires DATABASE_URL",
			envVars: map[string]string{
				"ENV":                     "production",
				"JWT_ACCESS_SECRET":       strongAccess,
				"JWT_REFRESH_SECRET":      strongRefresh,
				"DOCUMENT_ENCRYPTION_KEY": prodDocKey,
			},
			wantError: true,
		},
		{
			name: "Production requires DOCUMENT_ENCRYPTION_KEY",
			envVars: map[string]string{
				"ENV":                "production",
				"JWT_ACCESS_SECRET":  strongAccess,
				"JWT_REFRESH_SECRET": strongRefresh,
				"DATABASE_URL":       prodDB,
			},
			wantError: true,
		},
		{
			name: "Production rejects default DOCUMENT_ENCRYPTION_KEY",
			envVars: map[string]string{
				"ENV":                     "production",
				"JWT_ACCESS_SECRET":       strongAccess,
				"JWT_REFRESH_SECRET":      strongRefresh,
				"DATABASE_URL":            prodDB,
				"DOCUMENT_ENCRYPTION_KEY": DefaultDocumentEncryptionKeyBase64(),
			},
			wantError: true,
		},
		{
			name: "Production accepts strong JWT secrets",
			envVars: map[string]string{
				"ENV":                     "production",
				"JWT_ACCESS_SECRET":       strongAccess,
				"JWT_REFRESH_SECRET":      strongRefresh,
				"DATABASE_URL":            prodDB,
				"DOCUMENT_ENCRYPTION_KEY": prodDocKey,
			},
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "production", cfg.App.Env)
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
			clearEnv()
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
