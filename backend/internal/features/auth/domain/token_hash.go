package domain

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashRefreshToken returns a SHA-256 hex digest for storing refresh tokens at rest.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// KeyLookup returns a SHA-256 hex digest of a badge key for indexed lookup.
func KeyLookup(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
