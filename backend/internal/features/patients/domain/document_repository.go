package domain

import "context"

// DocumentRepository persists patient document metadata and encrypted blobs.
type DocumentRepository interface {
	Create(ctx context.Context, doc *Document) error
	ListByPatientID(ctx context.Context, patientID string) ([]*Document, error)
	GetByID(ctx context.Context, patientID, documentID string) (*Document, error)
	// GetEncryptedContent returns metadata plus ciphertext (never plaintext).
	GetEncryptedContent(ctx context.Context, patientID, documentID string) (*Document, []byte, error)
	Delete(ctx context.Context, patientID, documentID string) (*Document, error)
}
