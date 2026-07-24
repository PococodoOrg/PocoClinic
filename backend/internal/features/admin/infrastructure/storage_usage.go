package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	backupDirWarnBytes    = 2 * 1024 * 1024 * 1024  // 2 GB
	documentsDirWarnBytes = 5 * 1024 * 1024 * 1024  // 5 GB
)

func directorySizeBytes(root string) int64 {
	if root == "" {
		return 0
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return 0
	}

	var total int64
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total
}

func formatMegabytes(bytes int64) string {
	return fmt.Sprintf("%.1f MB", float64(bytes)/1024/1024)
}
