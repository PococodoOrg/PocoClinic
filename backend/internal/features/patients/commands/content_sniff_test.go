package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestSniffAndValidateContent_AcceptsPDF(t *testing.T) {
	content := []byte("%PDF-1.4 test content")
	detected, reader, err := sniffAndValidateContent(bytes.NewReader(content), "application/pdf", "report.pdf")
	if err != nil {
		t.Fatalf("expected pdf accepted: %v", err)
	}
	if detected != "application/pdf" {
		t.Fatalf("detected=%q", detected)
	}
	out, err := readAllLimited(reader, 64)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.HasPrefix(out, []byte("%PDF")) {
		t.Fatalf("reader prefix lost: %q", out)
	}
}

func TestSniffAndValidateContent_AcceptsPNG(t *testing.T) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00}
	_, _, err := sniffAndValidateContent(bytes.NewReader(pngHeader), "image/png", "scan.png")
	if err != nil {
		t.Fatalf("expected png accepted: %v", err)
	}
}

func TestSniffAndValidateContent_RejectsHTML(t *testing.T) {
	content := []byte("<html><body>not a chart document</body></html>")
	_, _, err := sniffAndValidateContent(bytes.NewReader(content), "application/pdf", "fake.pdf")
	if err == nil {
		t.Fatal("expected html disguised as pdf to be rejected")
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSniffAndValidateContent_RejectsEmptyFile(t *testing.T) {
	_, _, err := sniffAndValidateContent(bytes.NewReader(nil), "application/pdf", "empty.pdf")
	if err == nil {
		t.Fatal("expected empty file rejection")
	}
}

func TestSniffAndValidateContent_RejectsDeclaredTypeMismatch(t *testing.T) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00}
	_, _, err := sniffAndValidateContent(bytes.NewReader(pngHeader), "application/pdf", "scan.pdf")
	if err == nil {
		t.Fatal("expected declared type mismatch rejection")
	}
	if !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func readAllLimited(r interface{ Read([]byte) (int, error) }, n int) ([]byte, error) {
	buf := make([]byte, n)
	total := 0
	for total < n {
		count, err := r.Read(buf[total:])
		total += count
		if err != nil {
			break
		}
	}
	return buf[:total], nil
}
