package infrastructure

import (
	"context"
	"runtime"
	"time"

	admindomain "github.com/PococodoOrg/PocoClinic/internal/features/admin/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

// PerformanceService collects lightweight runtime metrics.
type PerformanceService struct {
	pool *database.DB
}

func NewPerformanceService(pool *database.DB) *PerformanceService {
	return &PerformanceService{pool: pool}
}

func (s *PerformanceService) Snapshot(_ context.Context) *admindomain.PerformanceSnapshot {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	snapshot := &admindomain.PerformanceSnapshot{
		CollectedAt:   time.Now().UTC(),
		Goroutines:    runtime.NumGoroutine(),
		MemoryAllocMB: float64(mem.Alloc) / 1024 / 1024,
		MemorySysMB:   float64(mem.Sys) / 1024 / 1024,
	}

	// SQLite has no client connection pool stats comparable to pgx; omit DatabasePool.
	_ = s.pool
	return snapshot
}

var _ admindomain.PerformanceReader = (*PerformanceService)(nil)
