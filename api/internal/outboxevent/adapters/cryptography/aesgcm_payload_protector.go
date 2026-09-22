package cryptography

import (
	"context"

	corecryptography "github.com/eduardoaugustolb/versum/api/internal/cryptography"
	"github.com/eduardoaugustolb/versum/api/internal/cryptography/keyring"
	ports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
)

type AESGCMPayloadProtector struct {
	key keyring.Keyring
}

var _ ports.PayloadProtector = (*AESGCMPayloadProtector)(nil)

func NewAESGCMPayloadProtector(k keyring.Keyring) *AESGCMPayloadProtector {
	return &AESGCMPayloadProtector{key: k}
}

func (p *AESGCMPayloadProtector) Protect(_ context.Context, payload []byte) (ports.ProtectedPayload, error) {
	encryptionKeyVersion, encryptionKey, err := p.key.CurrentEncryptionKey()
	if err != nil {
		return ports.ProtectedPayload{}, err
	}

	ciphertext, err := corecryptography.EncryptGCM(payload, encryptionKey)
	if err != nil {
		return ports.ProtectedPayload{}, err
	}

	return ports.ProtectedPayload{
		Ciphertext: ciphertext,
		KeyVersion: encryptionKeyVersion,
	}, nil
}

func (p *AESGCMPayloadProtector) Unprotect(_ context.Context, protected ports.ProtectedPayload) ([]byte, error) {
	encryptionKey, err := p.key.EncryptionKey(protected.KeyVersion)
	if err != nil {
		return nil, err
	}

	payload, err := corecryptography.DecryptGCM(protected.Ciphertext, encryptionKey)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
