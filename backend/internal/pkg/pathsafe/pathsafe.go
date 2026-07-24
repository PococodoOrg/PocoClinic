package pathsafe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// JoinRoot joins elem under rootDir and ensures the result cannot escape rootDir.
func JoinRoot(rootDir string, elem ...string) (string, error) {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" {
		return "", fmt.Errorf("root directory is required")
	}

	for _, part := range elem {
		if filepath.IsAbs(part) {
			return "", fmt.Errorf("path escapes root directory")
		}
	}

	rootAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("resolve root directory: %w", err)
	}
	rootAbs = filepath.Clean(rootAbs)

	joined := filepath.Join(append([]string{rootAbs}, elem...)...)
	joinedAbs, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	joinedAbs = filepath.Clean(joinedAbs)

	rel, err := filepath.Rel(rootAbs, joinedAbs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes root directory")
	}
	return joinedAbs, nil
}

// ValidateBackupFilename ensures name is a safe backup bundle basename.
func ValidateBackupFilename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("backup filename is required")
	}
	if name != filepath.Base(name) {
		return fmt.Errorf("invalid backup filename")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return fmt.Errorf("invalid backup filename")
	}
	if !strings.HasPrefix(name, "pococlinic-backup-") || !strings.HasSuffix(name, ".tar.gz") {
		return fmt.Errorf("invalid backup filename")
	}
	return nil
}

// OpenBackupFile opens a validated backup bundle under rootDir.
func OpenBackupFile(rootDir, filename string) (*os.File, error) {
	if err := ValidateBackupFilename(filename); err != nil {
		return nil, err
	}
	return openBasenameUnderRoot(rootDir, filename)
}

// StatBackupFile returns metadata for a validated backup bundle under rootDir.
func StatBackupFile(rootDir, filename string) (os.FileInfo, error) {
	if err := ValidateBackupFilename(filename); err != nil {
		return nil, err
	}
	return statBasenameUnderRoot(rootDir, filename)
}

// validateBasename rejects path separators and parent-directory references in a
// single path component before it is joined under a root directory.
func validateBasename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("path component is required")
	}
	if filepath.IsAbs(name) {
		return fmt.Errorf("invalid path component")
	}
	if name != filepath.Base(name) {
		return fmt.Errorf("invalid path component")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return fmt.Errorf("invalid path component")
	}
	return nil
}

// resolveBasenameUnderRoot resolves a validated single-component name under rootDir.
func resolveBasenameUnderRoot(rootDir, basename string) (string, error) {
	if err := validateBasename(basename); err != nil {
		return "", err
	}

	rootAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("resolve root directory: %w", err)
	}
	rootAbs = filepath.Clean(rootAbs)

	joinedAbs, err := filepath.Abs(filepath.Join(rootAbs, basename))
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	joinedAbs = filepath.Clean(joinedAbs)

	rel, err := filepath.Rel(rootAbs, joinedAbs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes root directory")
	}
	return joinedAbs, nil
}

func openBasenameUnderRoot(rootDir, basename string) (*os.File, error) {
	path, err := resolveBasenameUnderRoot(rootDir, basename)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func statBasenameUnderRoot(rootDir, basename string) (os.FileInfo, error) {
	path, err := resolveBasenameUnderRoot(rootDir, basename)
	if err != nil {
		return nil, err
	}
	return os.Stat(path)
}

// ValidateArchiveEntry rejects tar paths that could escape on extract (zip slip).
func ValidateArchiveEntry(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("invalid archive entry")
	}
	if filepath.IsAbs(name) {
		return fmt.Errorf("invalid archive entry")
	}

	clean := filepath.ToSlash(filepath.Clean(name))
	if strings.HasPrefix(clean, "/") {
		return fmt.Errorf("invalid archive entry")
	}
	if strings.Contains(clean, ":") {
		return fmt.Errorf("invalid archive entry")
	}
	if clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		return fmt.Errorf("invalid archive entry")
	}
	return nil
}

// ValidateRelativeKey ensures a slash-separated storage key has no traversal segments.
func ValidateRelativeKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("storage key is required")
	}
	if filepath.IsAbs(key) {
		return fmt.Errorf("invalid storage key")
	}

	clean := filepath.ToSlash(filepath.Clean(key))
	if strings.Contains(clean, ":") {
		return fmt.Errorf("invalid storage key")
	}
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		return fmt.Errorf("invalid storage key")
	}
	return nil
}
