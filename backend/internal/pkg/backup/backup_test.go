package backup

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendDocumentFilesEmptyDir(t *testing.T) {
	dir := t.TempDir()
	checksums := map[string]string{}
	if err := appendDocumentFiles(nil, dir, checksums); err != nil {
		t.Fatalf("empty dir: %v", err)
	}
	if len(checksums) != 0 {
		t.Fatalf("expected no checksums, got %d", len(checksums))
	}
}

func TestAppendDocumentFilesMissingDir(t *testing.T) {
	checksums := map[string]string{}
	if err := appendDocumentFiles(nil, t.TempDir()+"/missing", checksums); err != nil {
		t.Fatalf("missing dir should be ignored: %v", err)
	}
}

func TestDecodeJSONValue_BytesAndTime(t *testing.T) {
	bytesValue, err := decodeJSONValue(map[string]interface{}{
		"__type": "bytes",
		"hex":    "48656c6c6f",
	})
	if err != nil {
		t.Fatalf("decode bytes: %v", err)
	}
	asBytes, ok := bytesValue.([]byte)
	if !ok || string(asBytes) != "Hello" {
		t.Fatalf("unexpected bytes value: %#v", bytesValue)
	}

	timeValue, err := decodeJSONValue("2026-07-22T12:00:00Z")
	if err != nil {
		t.Fatalf("decode time: %v", err)
	}
	if _, ok := timeValue.(string); ok {
		t.Fatal("expected parsed time, got string")
	}
}

func TestArchiveVerifyDetectsTampering(t *testing.T) {
	files := map[string][]byte{
		"database/patients.jsonl": []byte(`{"id":"1"}` + "\n"),
	}
	manifest := Manifest{
		Version:        1,
		AppVersion:     "test",
		DatabaseFormat: "jsonl",
		Checksums: map[string]string{
			"database/patients.jsonl": "deadbeef",
		},
	}

	archive := &Archive{
		Path:     "test.tar.gz",
		Manifest: manifest,
		Files:    files,
	}
	if err := archive.Verify(); err == nil {
		t.Fatal("expected checksum verification to fail")
	}
}

func TestReadArchiveRejectsZipSlip(t *testing.T) {
	root := t.TempDir()
	filename := "pococlinic-backup-evil.tar.gz"
	path := filepath.Join(root, filename)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	gz := gzip.NewWriter(file)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{
		Name:     "../escape.txt",
		Typeflag: tar.TypeReg,
		Size:     4,
	}); err != nil {
		t.Fatalf("write header: %v", err)
	}
	if _, err := tw.Write([]byte("evil")); err != nil {
		t.Fatalf("write body: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}

	if _, err := readArchiveFiles(root, filename); err == nil {
		t.Fatal("expected zip-slip archive to be rejected")
	}
}
