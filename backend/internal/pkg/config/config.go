package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/pkg/doccrypto"
	"github.com/gin-contrib/cors"
)

const (
	defaultAccessTokenSecret  = "dev-access-secret-change-me"
	defaultRefreshTokenSecret = "dev-refresh-secret-change-me"
	// Exactly 32 bytes; base64 form is the default DOCUMENT_ENCRYPTION_KEY in development.
	defaultDocumentEncryptionKeyMaterial = "dev-document-key-change-me!!!!!!"
	minProductionSecretLen               = 32
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Security SecurityConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Backup   BackupConfig
	Storage  StorageConfig
	App      AppConfig
}

// StorageConfig holds document encryption and legacy disk paths.
type StorageConfig struct {
	// DocumentsDir is used for legacy disk-backed files and bootstrap credential paths.
	DocumentsDir string
	// DocumentEncryptionKey is a 32-byte AES key (plaintext never stored in the DB).
	DocumentEncryptionKey []byte
}

// AppConfig holds application metadata
type AppConfig struct {
	Version string
	Env     string
}

// BackupConfig holds backup directory settings
type BackupConfig struct {
	Dir string
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	URL           string
	RunMigrations bool
}

// AuthConfig holds JWT and authentication settings
type AuthConfig struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	SessionInactivityTTL   time.Duration
	Issuer                 string
	CookieSecure           bool
	RefreshTokenCookieName string
	AccessTokenCookieName  string
}

// ServerConfig holds all server-related configuration
type ServerConfig struct {
	Port int
	Host string
}

// SecurityConfig holds all security-related configuration
type SecurityConfig struct {
	AllowedOrigins []string
	TrustedProxies []string
	RateLimit      RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
}

// LoadConfig loads configuration from environment variables (and optional .env files).
func LoadConfig() (*Config, error) {
	loadDotEnvFiles()

	config := &Config{}

	// Server configuration
	port, err := strconv.Atoi(getEnvOrDefault("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	config.Server.Port = port
	config.Server.Host = getEnvOrDefault("SERVER_HOST", "localhost")

	// Security configuration
	config.Security.AllowedOrigins = []string{
		getEnvOrDefault("ALLOWED_ORIGIN", "http://localhost:3000"),
	}
	if proxies := strings.TrimSpace(getEnvOrDefault("TRUSTED_PROXIES", "")); proxies != "" {
		config.Security.TrustedProxies = strings.Split(proxies, ",")
		for i := range config.Security.TrustedProxies {
			config.Security.TrustedProxies[i] = strings.TrimSpace(config.Security.TrustedProxies[i])
		}
	}

	// Rate limit configuration
	rps, err := strconv.Atoi(getEnvOrDefault("RATE_LIMIT_RPS", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS: %w", err)
	}
	burst, err := strconv.Atoi(getEnvOrDefault("RATE_LIMIT_BURST", "20"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_BURST: %w", err)
	}
	config.Security.RateLimit = RateLimitConfig{
		RequestsPerSecond: rps,
		BurstSize:         burst,
	}

	accessTTL, err := time.ParseDuration(getEnvOrDefault("JWT_ACCESS_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}
	refreshTTL, err := time.ParseDuration(getEnvOrDefault("JWT_REFRESH_TTL", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}
	inactivityTTL, err := time.ParseDuration(getEnvOrDefault("SESSION_INACTIVITY_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_INACTIVITY_TTL: %w", err)
	}

	env := getEnvOrDefault("ENV", "development")
	config.App.Env = env
	config.App.Version = getEnvOrDefault("APP_VERSION", "dev")
	config.Backup.Dir = getEnvOrDefault("BACKUP_DIR", "./backups")
	config.Storage.DocumentsDir = getEnvOrDefault("DOCUMENTS_DIR", "./data/documents")

	docKey, err := loadDocumentEncryptionKey(env)
	if err != nil {
		return nil, err
	}
	config.Storage.DocumentEncryptionKey = docKey

	config.Auth = AuthConfig{
		AccessTokenSecret:      getEnvOrDefault("JWT_ACCESS_SECRET", defaultAccessTokenSecret),
		RefreshTokenSecret:     getEnvOrDefault("JWT_REFRESH_SECRET", defaultRefreshTokenSecret),
		AccessTokenTTL:           accessTTL,
		RefreshTokenTTL:          refreshTTL,
		SessionInactivityTTL:     inactivityTTL,
		Issuer:                   getEnvOrDefault("JWT_ISSUER", "pococlinic"),
		CookieSecure:             env == "production",
		RefreshTokenCookieName:   getEnvOrDefault("REFRESH_TOKEN_COOKIE", "poco_refresh_token"),
		AccessTokenCookieName:    getEnvOrDefault("ACCESS_TOKEN_COOKIE", "poco_access_token"),
	}

	config.Database.URL = getEnvOrDefault("DATABASE_URL", "")
	config.Database.RunMigrations = parseRunMigrations(env)

	if err := validateProductionSecrets(config); err != nil {
		return nil, err
	}

	if config.App.Env == "production" && config.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required in production")
	}

	return config, nil
}

func validateProductionSecrets(cfg *Config) error {
	if cfg.App.Env != "production" {
		return nil
	}

	switch {
	case cfg.Auth.AccessTokenSecret == defaultAccessTokenSecret:
		return fmt.Errorf("JWT_ACCESS_SECRET must be set to a non-default value in production")
	case cfg.Auth.RefreshTokenSecret == defaultRefreshTokenSecret:
		return fmt.Errorf("JWT_REFRESH_SECRET must be set to a non-default value in production")
	case len(cfg.Auth.AccessTokenSecret) < minProductionSecretLen:
		return fmt.Errorf("JWT_ACCESS_SECRET must be at least %d characters in production", minProductionSecretLen)
	case len(cfg.Auth.RefreshTokenSecret) < minProductionSecretLen:
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least %d characters in production", minProductionSecretLen)
	case cfg.Auth.AccessTokenSecret == cfg.Auth.RefreshTokenSecret:
		return fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must differ in production")
	case bytesEqual(cfg.Storage.DocumentEncryptionKey, []byte(defaultDocumentEncryptionKeyMaterial)):
		return fmt.Errorf("DOCUMENT_ENCRYPTION_KEY must be set to a non-default value in production")
	case len(cfg.Storage.DocumentEncryptionKey) != doccrypto.KeySize:
		return fmt.Errorf("DOCUMENT_ENCRYPTION_KEY must decode to %d bytes in production", doccrypto.KeySize)
	}

	return nil
}

func loadDocumentEncryptionKey(env string) ([]byte, error) {
	encoded := strings.TrimSpace(os.Getenv("DOCUMENT_ENCRYPTION_KEY"))
	if encoded == "" {
		if env == "production" {
			return nil, fmt.Errorf("DOCUMENT_ENCRYPTION_KEY is required in production")
		}
		return []byte(defaultDocumentEncryptionKeyMaterial), nil
	}
	key, err := doccrypto.ParseKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("DOCUMENT_ENCRYPTION_KEY: %w", err)
	}
	return key, nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// DefaultDocumentEncryptionKeyBase64 returns the development default key encoding (for docs/examples).
func DefaultDocumentEncryptionKeyBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(defaultDocumentEncryptionKeyMaterial))
}

func parseRunMigrations(env string) bool {
	if value, exists := os.LookupEnv("RUN_MIGRATIONS"); exists {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return env != "production"
		}
		return parsed
	}
	// Dev/single-node: migrate on app start for convenience.
	// Production/cluster: set RUN_MIGRATIONS=false and run cmd/migrate as a deploy step.
	return env != "production"
}

// TokenConfig returns JWT settings for the auth domain layer
func (c *Config) TokenConfig() domain.TokenConfig {
	return domain.TokenConfig{
		AccessTokenSecret:    []byte(c.Auth.AccessTokenSecret),
		RefreshTokenSecret:   []byte(c.Auth.RefreshTokenSecret),
		AccessTokenTTL:       c.Auth.AccessTokenTTL,
		RefreshTokenTTL:      c.Auth.RefreshTokenTTL,
		SessionInactivityTTL: c.Auth.SessionInactivityTTL,
		Issuer:               c.Auth.Issuer,
	}
}

// ConfigureCORS returns CORS configuration based on the current environment
func (c *Config) ConfigureCORS() cors.Config {
	return cors.Config{
		AllowOrigins:     c.Security.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		// Don't use AllowOriginFunc as it causes 403s
		// Instead, let the request through and control access via AllowOrigins
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
