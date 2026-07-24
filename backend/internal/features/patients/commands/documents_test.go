package commands

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	patientinfra "github.com/PococodoOrg/PocoClinic/internal/features/patients/infrastructure"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/doccrypto"
	"github.com/google/uuid"
)

func TestUploadAndOpenEncryptedDocument(t *testing.T) {
	patients := patientinfra.NewMemoryRepository()
	patient := domain.NewPatient("Ada", "Lovelace", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), domain.GenderFemale)
	if err := patients.Create(context.Background(), patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}

	docs := patientinfra.NewMemoryDocumentRepository(patients)
	docs.SetUploaderName(uuid.MustParse("11111111-1111-1111-1111-111111111111"), "Therapist")

	key := bytes.Repeat([]byte{7}, doccrypto.KeySize)
	cipher, err := doccrypto.NewCipher(key)
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}

	uploader := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	upload := NewUploadDocumentHandler(docs, patients, cipher)
	plain := []byte("%PDF-1.4 encrypted document body")
	doc, err := upload.Handle(context.Background(), UploadDocumentCommand{
		PatientID:   patient.ID.String(),
		UploadedBy:  uploader,
		FileName:    "plan.pdf",
		ContentType: "application/pdf",
		Content:     bytes.NewReader(plain),
	})
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	if doc.EncryptedContent != nil {
		t.Fatal("API document must not include ciphertext")
	}

	open := NewOpenDocumentContentHandler(docs, cipher, nil)
	reader, opened, err := open.Handle(context.Background(), patient.ID.String(), doc.ID.String())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer reader.Close()
	if opened.FileName != "plan.pdf" {
		t.Fatalf("unexpected file name %q", opened.FileName)
	}
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("plaintext mismatch: got %q", got)
	}

	_, enc, err := docs.GetEncryptedContent(context.Background(), patient.ID.String(), doc.ID.String())
	if err != nil {
		t.Fatalf("get encrypted: %v", err)
	}
	if len(enc) == 0 {
		t.Fatal("expected encrypted blob in repository")
	}
	if bytes.Equal(enc, plain) {
		t.Fatal("repository must store ciphertext, not plaintext")
	}
}

func TestDeleteDocumentRemovesRow(t *testing.T) {
	patients := patientinfra.NewMemoryRepository()
	patient := domain.NewPatient("Grace", "Hopper", time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC), domain.GenderFemale)
	if err := patients.Create(context.Background(), patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}
	docs := patientinfra.NewMemoryDocumentRepository(patients)
	cipher, _ := doccrypto.NewCipher(bytes.Repeat([]byte{1}, doccrypto.KeySize))
	uploader := uuid.New()
	docs.SetUploaderName(uploader, "Admin")

	upload := NewUploadDocumentHandler(docs, patients, cipher)
	doc, err := upload.Handle(context.Background(), UploadDocumentCommand{
		PatientID:   patient.ID.String(),
		UploadedBy:  uploader,
		FileName:    "note.txt",
		ContentType: "text/plain",
		Content:     bytes.NewReader([]byte("hello exercise notes")),
	})
	if err != nil {
		t.Fatalf("upload: %v", err)
	}

	del := NewDeleteDocumentHandler(docs, nil)
	if _, err := del.Handle(context.Background(), patient.ID.String(), doc.ID.String()); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := docs.GetByID(context.Background(), patient.ID.String(), doc.ID.String()); err == nil {
		t.Fatal("expected document gone after delete")
	}
}
