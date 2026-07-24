package infrastructure

import (
	"bytes"
	"context"
	"io"
	"sync"

	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/google/uuid"
)

// MemoryFileStorage is a legacy in-memory disk stand-in for pre-migration files.
type memoryFile struct {
	data []byte
}

type MemoryFileStorage struct {
	mu    sync.RWMutex
	files map[string]memoryFile
}

func NewMemoryFileStorage() *MemoryFileStorage {
	return &MemoryFileStorage{files: make(map[string]memoryFile)}
}

func (s *MemoryFileStorage) Save(_ context.Context, patientID string, content io.Reader, maxBytes int64) (string, int64, error) {
	docID := uuid.New().String()
	key := patientID + "/" + docID
	data, err := io.ReadAll(io.LimitReader(content, maxBytes+1))
	if err != nil {
		return "", 0, err
	}
	if int64(len(data)) > maxBytes {
		return "", 0, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File exceeds maximum size")
	}
	s.mu.Lock()
	s.files[key] = memoryFile{data: data}
	s.mu.Unlock()
	return key, int64(len(data)), nil
}

func (s *MemoryFileStorage) Open(_ context.Context, storageKey string) (io.ReadCloser, error) {
	s.mu.RLock()
	file, ok := s.files[storageKey]
	s.mu.RUnlock()
	if !ok {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document file not found")
	}
	return io.NopCloser(bytes.NewReader(file.data)), nil
}

func (s *MemoryFileStorage) Delete(_ context.Context, storageKey string) error {
	s.mu.Lock()
	delete(s.files, storageKey)
	s.mu.Unlock()
	return nil
}

type MemoryDocumentRepository struct {
	mu        sync.RWMutex
	docs      map[uuid.UUID]*domain.Document
	patients  domain.PatientRepository
	userNames map[uuid.UUID]string
}

func NewMemoryDocumentRepository(patients domain.PatientRepository) *MemoryDocumentRepository {
	return &MemoryDocumentRepository{
		docs:      make(map[uuid.UUID]*domain.Document),
		patients:  patients,
		userNames: make(map[uuid.UUID]string),
	}
}

func (r *MemoryDocumentRepository) SetUploaderName(id uuid.UUID, name string) {
	r.mu.Lock()
	r.userNames[id] = name
	r.mu.Unlock()
}

func (r *MemoryDocumentRepository) Create(ctx context.Context, doc *domain.Document) error {
	if _, err := r.patients.GetByID(ctx, doc.PatientID.String()); err != nil {
		return err
	}
	r.mu.Lock()
	copy := *doc
	if len(doc.EncryptedContent) > 0 {
		copy.EncryptedContent = append([]byte(nil), doc.EncryptedContent...)
	}
	if name := r.userNames[doc.UploadedBy]; name != "" {
		copy.UploadedByName = name
	} else {
		copy.UploadedByName = "Staff"
	}
	r.docs[doc.ID] = &copy
	r.mu.Unlock()
	return nil
}

func (r *MemoryDocumentRepository) ListByPatientID(ctx context.Context, patientID string) ([]*domain.Document, error) {
	if _, err := r.patients.GetByID(ctx, patientID); err != nil {
		return nil, err
	}
	parsed, err := uuid.Parse(patientID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid patient ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Document, 0)
	for _, doc := range r.docs {
		if doc.PatientID == parsed {
			copy := *doc
			copy.EncryptedContent = nil
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (r *MemoryDocumentRepository) GetByID(ctx context.Context, patientID, documentID string) (*domain.Document, error) {
	if _, err := r.patients.GetByID(ctx, patientID); err != nil {
		return nil, err
	}
	parsedDoc, err := uuid.Parse(documentID)
	if err != nil {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid document ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	doc, ok := r.docs[parsedDoc]
	if !ok {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document not found")
	}
	copy := *doc
	copy.EncryptedContent = nil
	return &copy, nil
}

func (r *MemoryDocumentRepository) GetEncryptedContent(ctx context.Context, patientID, documentID string) (*domain.Document, []byte, error) {
	if _, err := r.patients.GetByID(ctx, patientID); err != nil {
		return nil, nil, err
	}
	parsedDoc, err := uuid.Parse(documentID)
	if err != nil {
		return nil, nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Invalid document ID")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	doc, ok := r.docs[parsedDoc]
	if !ok {
		return nil, nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document not found")
	}
	copy := *doc
	copy.EncryptedContent = nil
	var encrypted []byte
	if len(doc.EncryptedContent) > 0 {
		encrypted = append([]byte(nil), doc.EncryptedContent...)
	}
	return &copy, encrypted, nil
}

func (r *MemoryDocumentRepository) Delete(ctx context.Context, patientID, documentID string) (*domain.Document, error) {
	doc, err := r.GetByID(ctx, patientID, documentID)
	if err != nil {
		return nil, err
	}
	parsedDoc, _ := uuid.Parse(documentID)
	r.mu.Lock()
	delete(r.docs, parsedDoc)
	r.mu.Unlock()
	return doc, nil
}

var _ domain.DocumentRepository = (*MemoryDocumentRepository)(nil)
