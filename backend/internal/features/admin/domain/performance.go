package domain

import (
	"context"
	"time"
)

// PerformanceSnapshot is a lightweight runtime health sample for administrators.
type PerformanceSnapshot struct {
	CollectedAt   time.Time  `json:"collectedAt"`
	Goroutines    int        `json:"goroutines"`
	MemoryAllocMB float64    `json:"memoryAllocMB"`
	MemorySysMB   float64    `json:"memorySysMB"`
	DatabasePool  *PoolStats `json:"databasePool,omitempty"`
}

// PoolStats summarizes database connection pool usage.
type PoolStats struct {
	TotalConns int32 `json:"totalConns"`
	IdleConns  int32 `json:"idleConns"`
	InUseConns int32 `json:"inUseConns"`
	MaxConns   int32 `json:"maxConns"`
}

// PerformanceReader collects runtime performance samples.
type PerformanceReader interface {
	Snapshot(ctx context.Context) *PerformanceSnapshot
}
