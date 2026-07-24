package doccrypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestCipherRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{0x42}, KeySize)
	c, err := NewCipher(key)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}

	plain := []byte("%PDF-1.4 fake document bytes")
	enc, err := c.Encrypt(plain)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if bytes.Equal(enc, plain) {
		t.Fatal("ciphertext must differ from plaintext")
	}

	out, err := c.Decrypt(enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("round-trip mismatch: got %q", out)
	}
}

func TestCipherRejectsWrongKey(t *testing.T) {
	c1, _ := NewCipher(bytes.Repeat([]byte{1}, KeySize))
	c2, _ := NewCipher(bytes.Repeat([]byte{2}, KeySize))
	enc, err := c1.Encrypt([]byte("secret scan"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if _, err := c2.Decrypt(enc); err == nil {
		t.Fatal("expected decrypt failure with wrong key")
	}
}

func TestParseKey(t *testing.T) {
	raw := bytes.Repeat([]byte{9}, KeySize)
	encoded := base64.StdEncoding.EncodeToString(raw)
	key, err := ParseKey(encoded)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !bytes.Equal(key, raw) {
		t.Fatal("parsed key mismatch")
	}
	if _, err := ParseKey("not-base64!!"); err == nil {
		t.Fatal("expected parse error")
	}
	if _, err := ParseKey(""); err == nil {
		t.Fatal("expected empty error")
	}
	if _, err := ParseKey(base64.StdEncoding.EncodeToString([]byte("short"))); err == nil {
		t.Fatal("expected length error")
	}
}
