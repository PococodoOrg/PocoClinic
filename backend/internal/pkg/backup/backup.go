package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dksch/pococlinic/internal/pkg/database"
)

const (
	manifestVersion   = 1
	sqlitePayloadName = "database/pococlinic.sqlite"
)

// Manifest describes a backup bundle per ADR-0012.
type Manifest struct {
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"createdAt"`
	AppVersion     string            `json:"appVersion"`
	DatabaseFormat string            `json:"databaseFormat"`
	Checksums      map[string]string `json:"checksums"`
	Summary        *ManifestSummary  `json:"summary,omitempty"`
}

// Info summarizes a backup file on disk.
type Info struct {
	Filename  string
	Path      string
	CreatedAt time.Time
	Manifest  Manifest
}

// Create writes a compressed backup bundle to dir and returns the file path.
// The database payload is a consistent SQLite file copy (VACUUM INTO).
// When documentsDir is non-empty, legacy on-disk patient document files are included under documents/.
func Create(ctx context.Context, db *database.DB, dir, documentsDir, appVersion string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is required for backup")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	stamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("pococlinic-backup-%s.tar.gz", stamp)
	path := filepath.Join(dir, filename)

	tmpDB, err := os.CreateTemp(dir, "pococlinic-backup-*.sqlite")
	if err != nil {
		return "", err
	}
	tmpPath := tmpDB.Name()
	_ = tmpDB.Close()
	defer os.Remove(tmpPath)

	if err := vacuumInto(ctx, db, tmpPath); err != nil {
		return "", fmt.Errorf("sqlite backup: %w", err)
	}
	if err := database.QuickCheckFile(ctx, tmpPath); err != nil {
		return "", fmt.Errorf("sqlite backup failed integrity check: %w", err)
	}

	dbBytes, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", err
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gz := gzip.NewWriter(file)
	defer gz.Close()

	tarWriter := tar.NewWriter(gz)
	defer tarWriter.Close()

	checksums := map[string]string{}
	summary := ManifestSummary{TableRows: map[string]int{}}
	manifest := Manifest{
		Version:        manifestVersion,
		CreatedAt:      time.Now().UTC(),
		AppVersion:     appVersion,
		DatabaseFormat: "sqlite",
		Checksums:      checksums,
	}

	if err := writeTarBytes(tarWriter, sqlitePayloadName, dbBytes); err != nil {
		return "", err
	}
	checksums[sqlitePayloadName] = sha256Hex(dbBytes)

	if counts, err := countTables(ctx, db); err == nil {
		summary.TableRows = counts
	}

	configSnapshot, err := json.Marshal(map[string]string{
		"note": "Non-secret configuration snapshot. Secrets are excluded by design.",
	})
	if err != nil {
		return "", err
	}
	if err := writeTarBytes(tarWriter, "config/snapshot.json", configSnapshot); err != nil {
		return "", err
	}
	checksums["config/snapshot.json"] = sha256Hex(configSnapshot)

	if err := appendDocumentFiles(tarWriter, documentsDir, checksums); err != nil {
		return "", fmt.Errorf("append documents: %w", err)
	}

	summary.DocumentFiles = documentFileCount(checksums)
	manifest.Summary = &summary
	manifest.Checksums = checksums
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", err
	}
	if err := writeTarBytes(tarWriter, "manifest.json", manifestBytes); err != nil {
		return "", err
	}

	if err := tarWriter.Close(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}

	return path, nil
}

func vacuumInto(ctx context.Context, db *database.DB, destPath string) error {
	// VACUUM INTO creates a consistent copy without stopping readers for long.
	// Prefer forward slashes so Windows paths do not confuse SQLite string parsing.
	escaped := strings.ReplaceAll(filepath.ToSlash(destPath), "'", "''")
	_, err := db.Exec(ctx, fmt.Sprintf("VACUUM INTO '%s'", escaped))
	return err
}

func countTables(ctx context.Context, db *database.DB) (map[string]int, error) {
	tables := []string{
		"users", "sessions", "patients", "patient_documents", "patient_notes",
		"exercise_plans", "exercise_log_entries", "audit_logs",
		"form_groups", "form_templates", "form_submissions",
	}
	out := make(map[string]int, len(tables))
	for _, table := range tables {
		var n int
		if err := db.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&n); err != nil {
			return nil, err
		}
		out[table] = n
	}
	return out, nil
}

// FindLatest returns metadata for the newest backup in dir.
func FindLatest(dir string) (*Info, error) {
	backups, err := List(dir)
	if err != nil {
		return nil, err
	}
	if len(backups) == 0 {
		return nil, nil
	}
	latest := backups[len(backups)-1]
	return &latest, nil
}

func appendDocumentFiles(tw *tar.Writer, documentsDir string, checksums map[string]string) error {
	if strings.TrimSpace(documentsDir) == "" {
		return nil
	}
	info, err := os.Stat(documentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("documents path is not a directory")
	}

	return filepath.Walk(documentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(documentsDir, path)
		if err != nil {
			return err
		}
		name := "documents/" + filepath.ToSlash(rel)
		if err := writeTarBytes(tw, name, content); err != nil {
			return err
		}
		checksums[name] = sha256Hex(content)
		return nil
	})
}

func writeTarBytes(tw *tar.Writer, name string, content []byte) error {
	header := &tar.Header{
		Name:    name,
		Mode:    0o644,
		Size:    int64(len(content)),
		ModTime: time.Now().UTC(),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err := tw.Write(content)
	return err
}

func sha256Hex(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func readManifestFromArchive(path string) (Manifest, error) {
	files, err := readArchiveFiles(path)
	if err != nil {
		return Manifest{}, err
	}
	manifestBytes, ok := files["manifest.json"]
	if !ok {
		return Manifest{}, fmt.Errorf("manifest not found")
	}
	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func fileModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Now().UTC()
	}
	return info.ModTime().UTC()
}
