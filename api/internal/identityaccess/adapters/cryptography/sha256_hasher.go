package cryptography

import (
	"crypto/sha256"
)

type SHA256Hasher struct {
}

func (SHA256Hasher) Hash(value string) ([]byte, error) {
	sum := sha256.Sum256([]byte(value))
	return sum[:], nil
}
