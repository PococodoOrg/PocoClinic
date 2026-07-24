package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
)

// Archive is an opened backup bundle.
type Archive struct {
	Path     string
	Manifest Manifest
	Files    map[string][]byte
}

// Open reads and verifies a backup archive.
func Open(path string) (*Archive, error) {
	files, err := readArchiveFiles(path)
	if err != nil {
		return nil, err
	}

	manifestBytes, ok := files["manifest.json"]
	if !ok {
		return nil, fmt.Errorf("manifest not found")
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	archive := &Archive{
		Path:     path,
		Manifest: manifest,
		Files:    files,
	}
	if err := archive.Verify(); err != nil {
		return nil, err
	}
	return archive, nil
}

// Verify checks manifest checksums against archive contents.
func (a *Archive) Verify() error {
	for name, expected := range a.Manifest.Checksums {
		content, ok := a.Files[name]
		if !ok {
			return fmt.Errorf("missing file %s referenced in manifest", name)
		}
		if sha256Hex(content) != expected {
			return fmt.Errorf("checksum mismatch for %s", name)
		}
	}
	return nil
}

// Restore replaces the live SQLite database file from the archive and restores legacy document files when present.
// The *database.DB handle is closed and reopened against the replaced file.
func Restore(ctx context.Context, db *database.DB, path, documentsDir string) error {
	if db == nil {
		return fmt.Errorf("database connection is required for restore")
	}

	archive, err := Open(path)
	if err != nil {
		return err
	}

	switch archive.Manifest.DatabaseFormat {
	case "sqlite":
		payload, ok := archive.Files[sqlitePayloadName]
		if !ok || len(payload) == 0 {
			return fmt.Errorf("sqlite database payload missing from archive")
		}
		if err := db.ReplaceWithSQLiteBytes(ctx, payload); err != nil {
			return err
		}
	case "jsonl":
		return fmt.Errorf("jsonl backups from the Cockroach era are no longer supported; restore requires a SQLite backup")
	default:
		return fmt.Errorf("unsupported database format %q", archive.Manifest.DatabaseFormat)
	}

	return restoreDocumentFiles(archive, documentsDir)
}

func restoreDocumentFiles(archive *Archive, documentsDir string) error {
	if strings.TrimSpace(documentsDir) == "" {
		return nil
	}

	hasDocuments := false
	for name := range archive.Files {
		if strings.HasPrefix(name, "documents/") {
			hasDocuments = true
			break
		}
	}
	if !hasDocuments {
		return nil
	}

	if err := os.RemoveAll(documentsDir); err != nil {
		return fmt.Errorf("clear documents dir: %w", err)
	}
	if err := os.MkdirAll(documentsDir, 0o750); err != nil {
		return err
	}

	for name, content := range archive.Files {
		if !strings.HasPrefix(name, "documents/") {
			continue
		}
		rel := strings.TrimPrefix(name, "documents/")
		if rel == "" || strings.Contains(rel, "..") {
			continue
		}
		target := filepath.Join(documentsDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
			return err
		}
		if err := os.WriteFile(target, content, 0o640); err != nil {
			return fmt.Errorf("write document %s: %w", rel, err)
		}
	}
	return nil
}

// List returns backup bundles in dir sorted by filename (oldest first).
func List(dir string) ([]Info, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Info{}, nil
		}
		return nil, err
	}

	candidates := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, "pococlinic-backup-") && strings.HasSuffix(name, ".tar.gz") {
			candidates = append(candidates, name)
		}
	}
	sort.Strings(candidates)

	result := make([]Info, 0, len(candidates))
	for _, name := range candidates {
		path := filepathJoin(dir, name)
		info := Info{
			Filename: name,
			Path:     path,
		}
		if manifest, err := readManifestFromArchive(path); err == nil {
			info.CreatedAt = manifest.CreatedAt
			info.Manifest = manifest
		} else {
			info.CreatedAt = fileModTime(path)
		}
		result = append(result, info)
	}
	return result, nil
}

func readArchiveFiles(path string) (map[string][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	files := map[string][]byte{}
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}

		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		files[header.Name] = content
	}
	return files, nil
}

func filepathJoin(dir, name string) string {
	return strings.TrimRight(dir, `/\`) + string(os.PathSeparator) + name
}
