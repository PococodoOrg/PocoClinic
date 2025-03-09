package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Security SecurityConfig
}

// ServerConfig holds all server-related configuration
type ServerConfig struct {
	Port int
	Host string
}

// SecurityConfig holds all security-related configuration
type SecurityConfig struct {
	AllowedOrigins []string
	RateLimit      RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
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

	return config, nil
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
