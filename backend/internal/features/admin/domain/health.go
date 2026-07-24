package domain

import (
	"context"
	"time"
)

// HealthCheckItem is a single operational check result.
type HealthCheckItem struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Label   string `json:"label"`
	Message string `json:"message"`
}

// HealthCheckResult summarizes live system health for administrators.
type HealthCheckResult struct {
	CheckedAt time.Time         `json:"checkedAt"`
	Status    string            `json:"status"`
	Checks    []HealthCheckItem `json:"checks"`
}

// HealthChecker runs live integrity and readiness checks.
type HealthChecker interface {
	RunHealthCheck(ctx context.Context) (*HealthCheckResult, error)
}
