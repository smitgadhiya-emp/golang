package helper

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken returns a cryptographically-random 64-char hex string,
// suitable for one-time links such as password-reset tokens.
func GenerateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
