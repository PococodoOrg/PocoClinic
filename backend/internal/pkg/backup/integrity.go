package backup

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dksch/pococlinic/internal/pkg/database"
)

// ManifestSummary captures row and file counts at backup time for integrity checks.
type ManifestSummary struct {
	TableRows     map[string]int `json:"tableRows"`
	DocumentFiles int            `json:"documentFiles"`
}

// IntegrityResult reports structural checks beyond checksum verification.
type IntegrityResult struct {
	Valid   bool            `json:"valid"`
	Issues  []string        `json:"issues,omitempty"`
	Summary ManifestSummary `json:"summary"`
}

func documentFileCount(checksums map[string]string) int {
	count := 0
	for name := range checksums {
		if strings.HasPrefix(name, "documents/") {
			count++
		}
	}
	return count
}

// VerifyIntegrity checks archive structure for the active database format.
func (a *Archive) VerifyIntegrity() IntegrityResult {
	result := IntegrityResult{
		Valid: true,
		Summary: ManifestSummary{
			TableRows:     map[string]int{},
			DocumentFiles: 0,
		},
	}

	for name := range a.Files {
		if strings.HasPrefix(name, "documents/") && !strings.HasSuffix(name, "/") {
			result.Summary.DocumentFiles++
		}
	}

	if a.Manifest.Summary != nil {
		result.Summary.TableRows = a.Manifest.Summary.TableRows
		if a.Manifest.Summary.DocumentFiles != result.Summary.DocumentFiles {
			result.Valid = false
			result.Issues = append(result.Issues, fmt.Sprintf(
				"document file count mismatch: manifest %d, archive %d",
				a.Manifest.Summary.DocumentFiles,
				result.Summary.DocumentFiles,
			))
		}
	}

	switch a.Manifest.DatabaseFormat {
	case "sqlite":
		payload, ok := a.Files[sqlitePayloadName]
		if !ok || len(payload) == 0 {
			result.Valid = false
			result.Issues = append(result.Issues, "sqlite database payload missing")
			return result
		}
		if expected, ok := a.Manifest.Checksums[sqlitePayloadName]; ok && sha256Hex(payload) != expected {
			result.Valid = false
			result.Issues = append(result.Issues, "sqlite database checksum mismatch")
		}
		if err := quickCheckPayload(payload); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, err.Error())
		}
	case "jsonl":
		result.Valid = false
		result.Issues = append(result.Issues, "jsonl (legacy Cockroach) backups cannot be verified for restore on SQLite")
	default:
		result.Valid = false
		result.Issues = append(result.Issues, fmt.Sprintf("unknown database format %q", a.Manifest.DatabaseFormat))
	}

	return result
}

func quickCheckPayload(payload []byte) error {
	tmp, err := os.CreateTemp("", "pococlinic-verify-*.sqlite")
	if err != nil {
		return fmt.Errorf("temp file for quick_check: %w", err)
	}
	path := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(path)

	if err := os.WriteFile(path, payload, 0o600); err != nil {
		return err
	}
	if err := database.QuickCheckFile(context.Background(), path); err != nil {
		return fmt.Errorf("sqlite quick_check failed: %w", err)
	}
	return nil
}
