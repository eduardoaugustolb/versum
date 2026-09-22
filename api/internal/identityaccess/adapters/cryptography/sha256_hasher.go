package cryptography

import (
	"crypto/sha256"

	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
)

type SHA256Hasher struct {
}

var _ ports.LoginTokenHasher = SHA256Hasher{}

func (SHA256Hasher) Hash(value string) ([]byte, error) {
	sum := sha256.Sum256([]byte(value))
	return sum[:], nil
}
