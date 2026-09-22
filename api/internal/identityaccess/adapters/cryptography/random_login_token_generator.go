package cryptography

import (
	"crypto/rand"
	"encoding/base64"

	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
)

type RandomTokenGenerator struct {
}

var _ ports.LoginTokenGenerator = RandomTokenGenerator{}

func (RandomTokenGenerator) GenerateLoginToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}
