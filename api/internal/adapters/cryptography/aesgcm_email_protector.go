package cryptography

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/keyring"
)

type AESGCMEmailProtector struct {
	key keyring.Keyring
}

var _ ports.EmailProtector = (*AESGCMEmailProtector)(nil)

func NewAESGCMEmailProtector(k keyring.Keyring) *AESGCMEmailProtector {
	return &AESGCMEmailProtector{key: k}
}

func (p *AESGCMEmailProtector) Protect(ctx context.Context, email domain.Email) (ports.ProtectedEmail, error) {
	plaintext := []byte(email.String())
	encryptionKeyVersion, encryptionKey, err := p.key.CurrentEncryptionKey()
	if err != nil {
		return ports.ProtectedEmail{}, err
	}

	ciphertext, err := encryptGCM(plaintext, encryptionKey)
	if err != nil {
		return ports.ProtectedEmail{}, err
	}

	lookupHMACKeyVersion, lookupHMACKey, err := p.key.CurrentLookupKey()

	if err != nil {
		return ports.ProtectedEmail{}, err
	}

	lookupHMAC := lookupHMAC(email, lookupHMACKey)

	return ports.ProtectedEmail{
		Ciphertext:           ciphertext,
		LookupHMAC:           lookupHMAC,
		EncryptionKeyVersion: encryptionKeyVersion,
		LookupKeyVersion:     lookupHMACKeyVersion,
	}, nil
}

func (p *AESGCMEmailProtector) LookupHMAC(ctx context.Context, email domain.Email) ([]byte, error) {
	_, lookupHMACKey, err := p.key.CurrentLookupKey()
	if err != nil {
		return nil, err
	}

	lookupHMAC := lookupHMAC(email, lookupHMACKey)

	return lookupHMAC, nil
}

func lookupHMAC(email domain.Email, lookupHMACKey []byte) []byte {
	plaintext := []byte(email.String())
	lookupHMACHash := hmac.New(sha256.New, lookupHMACKey)
	lookupHMACHash.Write(plaintext)
	return lookupHMACHash.Sum(nil)
}

func (p *AESGCMEmailProtector) Unprotect(_ context.Context, protected ports.ProtectedEmail) (domain.Email, error) {
	encryptionKey, err := p.key.EncryptionKey(protected.EncryptionKeyVersion)
	if err != nil {
		return domain.Email{}, err
	}

	plaintext, err := decryptGCM(protected.Ciphertext, encryptionKey)
	if err != nil {
		return domain.Email{}, err
	}

	email, err := domain.ParseEmail(string(plaintext))
	if err != nil {
		return domain.Email{}, err
	}

	return email, nil
}

func encryptGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptGCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)

}
