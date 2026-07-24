package domain

import (
	"context"
	"time"
)

// ComplianceItem is one operational control check for HIPAA-oriented accountability.
type ComplianceItem struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Label   string `json:"label"`
	Message string `json:"message"`
	Control string `json:"control,omitempty"`
}

// ComplianceCheckResult summarizes clinic operational compliance controls.
type ComplianceCheckResult struct {
	CheckedAt time.Time        `json:"checkedAt"`
	Status    string           `json:"status"`
	Checks    []ComplianceItem `json:"checks"`
}

// ComplianceChecker runs operational compliance checks.
type ComplianceChecker interface {
	RunComplianceCheck(ctx context.Context) (*ComplianceCheckResult, error)
}
