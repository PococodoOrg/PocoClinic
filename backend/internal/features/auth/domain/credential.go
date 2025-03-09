package domain

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Credential represents a hashed authentication credential (key or PIN)
type Credential struct {
	Hash []byte
	Salt []byte
}

// NewCredential creates a new credential from a plaintext value
func NewCredential(value string) (*Credential, error) {
	// Generate random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Hash the value with the salt
	hash := hashValue(value, salt)

	return &Credential{
		Hash: hash,
		Salt: salt,
	}, nil
}

// Validate checks if the provided value matches this credential
func (c *Credential) Validate(value string) bool {
	hash := hashValue(value, c.Salt)
	return subtle.ConstantTimeCompare(hash, c.Hash) == 1
}

// hashValue creates a cryptographic hash of the value using the provided salt
func hashValue(value string, salt []byte) []byte {
	// Use Argon2id for password hashing
	return argon2.IDKey(
		[]byte(value),
		salt,
		1,       // time cost
		64*1024, // memory cost
		4,       // threads
		32,      // key length
	)
}

// String returns a base64 encoded representation of the credential
func (c *Credential) String() string {
	hash := base64.StdEncoding.EncodeToString(c.Hash)
	salt := base64.StdEncoding.EncodeToString(c.Salt)
	return hash + ":" + salt
}

// ParseCredential parses a base64 encoded credential string
func ParseCredential(s string) (*Credential, error) {
	// TODO: Implement if needed for persistence
	return nil, nil
}

// GenerateKey generates a new 64-bit key and returns its credentials
func GenerateKey() (string, *Credential, error) {
	key := make([]byte, 8) // 8 bytes = 64 bits
	if _, err := rand.Read(key); err != nil {
		return "", nil, fmt.Errorf("failed to generate key: %w", err)
	}

	cred, err := NewCredential(base64.URLEncoding.EncodeToString(key))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return base64.URLEncoding.EncodeToString(key), cred, nil
}

// GeneratePINCredential generates credentials for a PIN
func GeneratePINCredential(pin string) (*Credential, error) {
	if len(pin) != 4 {
		return nil, fmt.Errorf("PIN must be exactly 4 digits")
	}

	return NewCredential(pin)
}
