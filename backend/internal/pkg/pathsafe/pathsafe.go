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

	if joinedAbs != rootAbs && !hasPathPrefix(joinedAbs, rootAbs) {
		return "", fmt.Errorf("path escapes root directory")
	}
	return joinedAbs, nil
}

func hasPathPrefix(path, prefix string) bool {
	if path == prefix {
		return true
	}
	return strings.HasPrefix(path, prefix+string(os.PathSeparator))
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
	if strings.Contains(name, "..") {
		return fmt.Errorf("invalid backup filename")
	}
	if !strings.HasPrefix(name, "pococlinic-backup-") || !strings.HasSuffix(name, ".tar.gz") {
		return fmt.Errorf("invalid backup filename")
	}
	return nil
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
