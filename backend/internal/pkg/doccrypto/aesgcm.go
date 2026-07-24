package doccrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const KeySize = 32

// Cipher encrypts/decrypts document payloads with AES-256-GCM.
// Wire format: nonce || ciphertext+tag.
type Cipher struct {
	gcm cipher.AEAD
}

func NewCipher(key []byte) (*Cipher, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("document encryption key must be %d bytes", KeySize)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{gcm: gcm}, nil
}

// ParseKey decodes a base64-encoded 32-byte key.
func ParseKey(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, errors.New("document encryption key is empty")
	}
	key, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("document encryption key must be base64: %w", err)
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("document encryption key must decode to %d bytes", KeySize)
	}
	return key, nil
}

func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return c.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *Cipher) Decrypt(blob []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	if len(blob) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := blob[:nonceSize], blob[nonceSize:]
	return c.gcm.Open(nil, nonce, ciphertext, nil)
}
