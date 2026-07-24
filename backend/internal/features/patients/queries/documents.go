package queries

import (
	"context"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
)

type ListDocumentsQuery struct {
	PatientID string
}

type ListDocumentsHandler struct {
	docRepo domain.DocumentRepository
}

func NewListDocumentsHandler(docRepo domain.DocumentRepository) ListDocumentsHandler {
	return ListDocumentsHandler{docRepo: docRepo}
}

func (h ListDocumentsHandler) Handle(ctx context.Context, query ListDocumentsQuery) ([]*domain.Document, error) {
	return h.docRepo.ListByPatientID(ctx, query.PatientID)
}

type GetDocumentQuery struct {
	PatientID  string
	DocumentID string
}

type GetDocumentHandler struct {
	docRepo domain.DocumentRepository
}

func NewGetDocumentHandler(docRepo domain.DocumentRepository) GetDocumentHandler {
	return GetDocumentHandler{docRepo: docRepo}
}

func (h GetDocumentHandler) Handle(ctx context.Context, query GetDocumentQuery) (*domain.Document, error) {
	return h.docRepo.GetByID(ctx, query.PatientID, query.DocumentID)
}
