package infrastructure

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// LocalFileStorage stores patient document bytes on the clinic server disk.
type LocalFileStorage struct {
	rootDir string
}

func NewLocalFileStorage(rootDir string) (*LocalFileStorage, error) {
	if err := os.MkdirAll(rootDir, 0o750); err != nil {
		return nil, err
	}
	return &LocalFileStorage{rootDir: rootDir}, nil
}

func (s *LocalFileStorage) Save(ctx context.Context, patientID string, content io.Reader, maxBytes int64) (storageKey string, size int64, err error) {
	_ = ctx
	docID := uuid.New().String()
	storageKey = filepath.ToSlash(filepath.Join(patientID, docID))

	dir := filepath.Join(s.rootDir, patientID)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", 0, err
	}

	path := filepath.Join(s.rootDir, patientID, docID)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o640)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	written, err := io.Copy(file, io.LimitReader(content, maxBytes+1))
	if err != nil {
		_ = os.Remove(path)
		return "", 0, err
	}
	if written > maxBytes {
		_ = os.Remove(path)
		return "", 0, fmt.Errorf("file exceeds maximum size")
	}

	return storageKey, written, nil
}

func (s *LocalFileStorage) Open(_ context.Context, storageKey string) (io.ReadCloser, error) {
	path, err := s.resolve(storageKey)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (s *LocalFileStorage) Delete(_ context.Context, storageKey string) error {
	path, err := s.resolve(storageKey)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (s *LocalFileStorage) RootDir() string {
	return s.rootDir
}

func (s *LocalFileStorage) resolve(storageKey string) (string, error) {
	clean := filepath.Clean(storageKey)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid storage key")
	}
	full := filepath.Join(s.rootDir, filepath.FromSlash(clean))
	absRoot, err := filepath.Abs(s.rootDir)
	if err != nil {
		return "", err
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absFull, absRoot+string(os.PathSeparator)) && absFull != absRoot {
		return "", fmt.Errorf("invalid storage path")
	}
	return full, nil
}
