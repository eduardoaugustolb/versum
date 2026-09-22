package cryptography

import (
	"context"
	"sort"

	corecryptography "github.com/eduardoaugustolb/versum/api/internal/cryptography"
	"github.com/eduardoaugustolb/versum/api/internal/cryptography/keyring"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type AESGCMEmailProtector struct {
	key keyring.Keyring
}

var _ ports.EmailProtector = (*AESGCMEmailProtector)(nil)

func NewAESGCMEmailProtector(k keyring.Keyring) *AESGCMEmailProtector {
	return &AESGCMEmailProtector{key: k}
}

func (p *AESGCMEmailProtector) Protect(_ context.Context, email domain.Email) (ports.ProtectedEmail, error) {
	plaintext := []byte(email.String())
	encryptionKeyVersion, encryptionKey, err := p.key.CurrentEncryptionKey()
	if err != nil {
		return ports.ProtectedEmail{}, err
	}

	ciphertext, err := corecryptography.EncryptGCM(plaintext, encryptionKey)
	if err != nil {
		return ports.ProtectedEmail{}, err
	}

	lookupHMACKeyVersion, lookupHMACKey, err := p.key.CurrentLookupKey()

	if err != nil {
		return ports.ProtectedEmail{}, err
	}

	lookupHMAC := corecryptography.HMACSHA256(plaintext, lookupHMACKey)

	return ports.ProtectedEmail{
		Ciphertext:           ciphertext,
		LookupHMAC:           lookupHMAC,
		EncryptionKeyVersion: encryptionKeyVersion,
		LookupKeyVersion:     lookupHMACKeyVersion,
	}, nil
}

func (p *AESGCMEmailProtector) LookupHMAC(_ context.Context, email domain.Email) ([]byte, error) {
	_, lookupHMACKey, err := p.key.CurrentLookupKey()
	if err != nil {
		return nil, err
	}

	lookupHMAC := corecryptography.HMACSHA256([]byte(email.String()), lookupHMACKey)

	return lookupHMAC, nil
}

func (p *AESGCMEmailProtector) LookupCandidates(_ context.Context, email domain.Email) ([]ports.LookupCandidate, error) {
	currentLookupKeyVersion, _, err := p.key.CurrentLookupKey()
	if err != nil {
		return nil, err
	}
	lookupKeys, err := p.key.LookupKeys()
	if err != nil {
		return nil, err
	}

	versions := make([]int, 0, len(lookupKeys))
	for version := range lookupKeys {
		versions = append(versions, version)
	}
	sort.Slice(versions, func(i, j int) bool {
		if versions[i] == currentLookupKeyVersion {
			return true
		}
		if versions[j] == currentLookupKeyVersion {
			return false
		}
		return versions[i] < versions[j]
	})

	lookups := make([]ports.LookupCandidate, 0, len(versions))
	for _, version := range versions {
		lookups = append(lookups, ports.LookupCandidate{
			LookupKeyVersion: version,
			LookupHMAC:       corecryptography.HMACSHA256([]byte(email.String()), lookupKeys[version]),
		})
	}

	return lookups, nil
}

func (p *AESGCMEmailProtector) Unprotect(_ context.Context, protected ports.ProtectedEmail) (domain.Email, error) {
	encryptionKey, err := p.key.EncryptionKey(protected.EncryptionKeyVersion)
	if err != nil {
		return domain.Email{}, err
	}

	plaintext, err := corecryptography.DecryptGCM(protected.Ciphertext, encryptionKey)
	if err != nil {
		return domain.Email{}, err
	}

	email, err := domain.ParseEmail(string(plaintext))
	if err != nil {
		return domain.Email{}, err
	}

	return email, nil
}
