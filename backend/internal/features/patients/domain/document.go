package domain

import (
	"time"

	"github.com/google/uuid"
)

const MaxDocumentSizeBytes = 10 * 1024 * 1024 // 10 MB

var AllowedDocumentContentTypes = map[string]string{
	"application/pdf": "pdf",
	"image/png":       "png",
	"image/jpeg":      "jpg",
	"image/gif":       "gif",
	"text/plain":      "txt",
}

// Document is a file attached to a patient chart.
// Bytes are stored AES-GCM encrypted in the database (EncryptedContent).
// StorageKey is only set for legacy disk-backed rows.
type Document struct {
	ID               uuid.UUID `json:"id"`
	PatientID        uuid.UUID `json:"patientId"`
	UploadedBy       uuid.UUID `json:"uploadedBy"`
	UploadedByName   string    `json:"uploadedByName"`
	FileName         string    `json:"fileName"`
	ContentType      string    `json:"contentType"`
	SizeBytes        int64     `json:"sizeBytes"`
	StorageKey       string    `json:"-"`
	EncryptedContent []byte    `json:"-"`
	CreatedAt        time.Time `json:"createdAt"`
}

func NewDocument(patientID, uploadedBy uuid.UUID, fileName, contentType string, sizeBytes int64, encrypted []byte) *Document {
	return &Document{
		ID:               uuid.New(),
		PatientID:        patientID,
		UploadedBy:       uploadedBy,
		FileName:         fileName,
		ContentType:      contentType,
		SizeBytes:        sizeBytes,
		EncryptedContent: encrypted,
		CreatedAt:        time.Now(),
	}
}
