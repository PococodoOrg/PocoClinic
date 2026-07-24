package infrastructure

import (
	"context"
	"errors"

	pkgerrors "github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/dksch/pococlinic/internal/features/patients/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

type SQLDocumentRepository struct {
	pool *database.DB
}

func NewSQLDocumentRepository(pool *database.DB) *SQLDocumentRepository {
	return &SQLDocumentRepository{pool: pool}
}

func (r *SQLDocumentRepository) Create(ctx context.Context, doc *domain.Document) error {
	var storageKey *string
	if doc.StorageKey != "" {
		storageKey = &doc.StorageKey
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO patient_documents (
			id, patient_id, uploaded_by, file_name, content_type, size_bytes,
			storage_key, encrypted_content, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, doc.ID, doc.PatientID, doc.UploadedBy, doc.FileName, doc.ContentType, doc.SizeBytes,
		storageKey, doc.EncryptedContent, doc.CreatedAt)
	return err
}

func (r *SQLDocumentRepository) ListByPatientID(ctx context.Context, patientID string) ([]*domain.Document, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT d.id, d.patient_id, d.uploaded_by, u.name, d.file_name, d.content_type,
			d.size_bytes, COALESCE(d.storage_key, ''), d.created_at
		FROM patient_documents d
		JOIN users u ON u.id = d.uploaded_by
		WHERE d.patient_id = $1
		ORDER BY d.created_at DESC
	`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := make([]*domain.Document, 0)
	for rows.Next() {
		doc := &domain.Document{}
		if err := rows.Scan(
			&doc.ID, &doc.PatientID, &doc.UploadedBy, &doc.UploadedByName,
			&doc.FileName, &doc.ContentType, &doc.SizeBytes, &doc.StorageKey, database.Time(&doc.CreatedAt),
		); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (r *SQLDocumentRepository) GetByID(ctx context.Context, patientID, documentID string) (*domain.Document, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT d.id, d.patient_id, d.uploaded_by, u.name, d.file_name, d.content_type,
			d.size_bytes, COALESCE(d.storage_key, ''), d.created_at
		FROM patient_documents d
		JOIN users u ON u.id = d.uploaded_by
		WHERE d.patient_id = $1 AND d.id = $2
	`, patientID, documentID)
	return scanDocumentMeta(row)
}

func (r *SQLDocumentRepository) GetEncryptedContent(ctx context.Context, patientID, documentID string) (*domain.Document, []byte, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT d.id, d.patient_id, d.uploaded_by, u.name, d.file_name, d.content_type,
			d.size_bytes, COALESCE(d.storage_key, ''), d.encrypted_content, d.created_at
		FROM patient_documents d
		JOIN users u ON u.id = d.uploaded_by
		WHERE d.patient_id = $1 AND d.id = $2
	`, patientID, documentID)

	doc := &domain.Document{}
	var encrypted []byte
	err := row.Scan(
		&doc.ID, &doc.PatientID, &doc.UploadedBy, &doc.UploadedByName,
		&doc.FileName, &doc.ContentType, &doc.SizeBytes, &doc.StorageKey, &encrypted, database.Time(&doc.CreatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document not found")
		}
		return nil, nil, err
	}
	return doc, encrypted, nil
}

func (r *SQLDocumentRepository) Delete(ctx context.Context, patientID, documentID string) (*domain.Document, error) {
	doc, err := r.GetByID(ctx, patientID, documentID)
	if err != nil {
		return nil, err
	}
	tag, err := r.pool.Exec(ctx, `DELETE FROM patient_documents WHERE patient_id = $1 AND id = $2`, patientID, documentID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document not found")
	}
	return doc, nil
}

func scanDocumentMeta(row database.Row) (*domain.Document, error) {
	doc := &domain.Document{}
	err := row.Scan(
		&doc.ID, &doc.PatientID, &doc.UploadedBy, &doc.UploadedByName,
		&doc.FileName, &doc.ContentType, &doc.SizeBytes, &doc.StorageKey, database.Time(&doc.CreatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, pkgerrors.NewAPIError(pkgerrors.ErrNotFound, "Document not found")
		}
		return nil, err
	}
	return doc, nil
}

var _ domain.DocumentRepository = (*SQLDocumentRepository)(nil)
