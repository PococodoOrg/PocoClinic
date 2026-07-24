package commands

import (
	"bytes"
	"context"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/pkg/doccrypto"
	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/google/uuid"
)

// LegacyFileStore opens/deletes pre-migration disk-backed document files.
type LegacyFileStore interface {
	Open(ctx context.Context, storageKey string) (io.ReadCloser, error)
	Delete(ctx context.Context, storageKey string) error
}

type UploadDocumentCommand struct {
	PatientID   string
	UploadedBy  uuid.UUID
	FileName    string
	ContentType string
	Content     io.Reader
}

type UploadDocumentHandler struct {
	docRepo     domain.DocumentRepository
	patientRepo domain.GetPatientRepository
	cipher      *doccrypto.Cipher
}

func NewUploadDocumentHandler(
	docRepo domain.DocumentRepository,
	patientRepo domain.GetPatientRepository,
	cipher *doccrypto.Cipher,
) UploadDocumentHandler {
	return UploadDocumentHandler{docRepo: docRepo, patientRepo: patientRepo, cipher: cipher}
}

func (h UploadDocumentHandler) Handle(ctx context.Context, cmd UploadDocumentCommand) (*domain.Document, error) {
	fileName := strings.TrimSpace(filepath.Base(cmd.FileName))
	if fileName == "" || fileName == "." || fileName == ".." {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File name is required")
	}

	if _, err := h.patientRepo.GetByID(ctx, cmd.PatientID); err != nil {
		return nil, err
	}

	contentType, content, err := sniffAndValidateContent(cmd.Content, cmd.ContentType, fileName)
	if err != nil {
		return nil, err
	}

	patientID, err := uuid.Parse(cmd.PatientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}

	plain, err := io.ReadAll(io.LimitReader(content, domain.MaxDocumentSizeBytes+1))
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Failed to read uploaded file")
	}
	if int64(len(plain)) > domain.MaxDocumentSizeBytes {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File exceeds 10 MB limit")
	}
	if len(plain) == 0 {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File is empty")
	}

	encrypted, err := h.cipher.Encrypt(plain)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to encrypt document")
	}

	doc := domain.NewDocument(patientID, cmd.UploadedBy, fileName, contentType, int64(len(plain)), encrypted)
	if err := h.docRepo.Create(ctx, doc); err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to save document")
	}

	// Do not return ciphertext on the API response.
	doc.EncryptedContent = nil
	return doc, nil
}

type DeleteDocumentHandler struct {
	docRepo domain.DocumentRepository
	legacy  LegacyFileStore
}

func NewDeleteDocumentHandler(docRepo domain.DocumentRepository, legacy LegacyFileStore) DeleteDocumentHandler {
	return DeleteDocumentHandler{docRepo: docRepo, legacy: legacy}
}

func (h DeleteDocumentHandler) Handle(ctx context.Context, patientID, documentID string) (*domain.Document, error) {
	doc, err := h.docRepo.Delete(ctx, patientID, documentID)
	if err != nil {
		return nil, err
	}
	if doc.StorageKey != "" && h.legacy != nil {
		_ = h.legacy.Delete(ctx, doc.StorageKey)
	}
	return doc, nil
}

type OpenDocumentContentHandler struct {
	docRepo domain.DocumentRepository
	cipher  *doccrypto.Cipher
	legacy  LegacyFileStore
}

func NewOpenDocumentContentHandler(
	docRepo domain.DocumentRepository,
	cipher *doccrypto.Cipher,
	legacy LegacyFileStore,
) OpenDocumentContentHandler {
	return OpenDocumentContentHandler{docRepo: docRepo, cipher: cipher, legacy: legacy}
}

func (h OpenDocumentContentHandler) Handle(ctx context.Context, patientID, documentID string) (io.ReadCloser, *domain.Document, error) {
	doc, encrypted, err := h.docRepo.GetEncryptedContent(ctx, patientID, documentID)
	if err != nil {
		return nil, nil, err
	}

	if len(encrypted) > 0 {
		plain, err := h.cipher.Decrypt(encrypted)
		if err != nil {
			return nil, nil, pkgerrors.NewAPIError(pkgerrors.ErrInternalServer, "Failed to decrypt document")
		}
		return io.NopCloser(bytes.NewReader(plain)), doc, nil
	}

	if doc.StorageKey != "" && h.legacy != nil {
		reader, err := h.legacy.Open(ctx, doc.StorageKey)
		if err != nil {
			return nil, nil, err
		}
		return reader, doc, nil
	}

	return nil, nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document content not found")
}

func normalizeContentType(contentType, fileName string) string {
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	if contentType != "" && contentType != "application/octet-stream" {
		return contentType
	}
	if guessed := mime.TypeByExtension(filepath.Ext(fileName)); guessed != "" {
		return guessed
	}
	return contentType
}
