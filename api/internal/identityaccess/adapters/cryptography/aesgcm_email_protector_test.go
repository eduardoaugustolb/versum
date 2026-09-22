package cryptography

import (
	"context"
	"errors"
	"testing"

	corecryptography "github.com/eduardoaugustolb/versum/api/internal/cryptography"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

var errKeyringUnavailable = errors.New("keyring unavailable")

type testKeyring struct {
	encryptionVersion int
	encryptionKey     []byte
	lookupVersion     int
	lookupKey         []byte
	lookupKeys        map[int][]byte
	currentEncryptErr error
	currentLookupErr  error
}

func (k testKeyring) CurrentEncryptionKey() (int, []byte, error) {
	if k.currentEncryptErr != nil {
		return 0, nil, k.currentEncryptErr
	}
	return k.encryptionVersion, k.encryptionKey, nil
}

func (k testKeyring) EncryptionKey(version int) ([]byte, error) {
	if version != k.encryptionVersion {
		return nil, errKeyringUnavailable
	}
	return k.encryptionKey, nil
}

func (k testKeyring) CurrentLookupKey() (int, []byte, error) {
	if k.currentLookupErr != nil {
		return 0, nil, k.currentLookupErr
	}
	return k.lookupVersion, k.lookupKey, nil
}

func (k testKeyring) LookupKey(version int) ([]byte, error) {
	if version != k.lookupVersion {
		return nil, errKeyringUnavailable
	}
	return k.lookupKey, nil
}

func (k testKeyring) LookupKeys() (map[int][]byte, error) {
	if k.lookupKeys != nil {
		return k.lookupKeys, nil
	}
	return map[int][]byte{k.lookupVersion: k.lookupKey}, nil
}

func newTestEmailProtector(t *testing.T) (*AESGCMEmailProtector, domain.Email) {
	t.Helper()
	email, err := domain.ParseEmail("Ana@Example.com")
	if err != nil {
		t.Fatal(err)
	}
	return NewAESGCMEmailProtector(testKeyring{
		encryptionVersion: 1,
		encryptionKey:     []byte("01234567890123456789012345678901"),
		lookupVersion:     2,
		lookupKey:         []byte("abcdefghijklmnopqrstuvwxyz123456"),
	}), email
}

func TestAESGCMEmailProtectorRoundTripAndRandomNonce(t *testing.T) {
	protector, email := newTestEmailProtector(t)

	first, err := protector.Protect(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}
	second, err := protector.Protect(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}

	if first.EncryptionKeyVersion != 1 || first.LookupKeyVersion != 2 {
		t.Fatalf("unexpected key versions: encryption=%d lookup=%d", first.EncryptionKeyVersion, first.LookupKeyVersion)
	}
	if string(first.Ciphertext) == email.String() {
		t.Fatal("ciphertext must not contain the plaintext email")
	}
	if string(first.Ciphertext) == string(second.Ciphertext) {
		t.Fatal("encrypting the same email twice must use distinct nonces")
	}

	got, err := protector.Unprotect(t.Context(), first)
	if err != nil {
		t.Fatal(err)
	}
	if got != email {
		t.Fatalf("expected %q, got %q", email.String(), got.String())
	}
}

func TestAESGCMEmailProtectorRejectsTamperedCiphertext(t *testing.T) {
	protector, email := newTestEmailProtector(t)
	protected, err := protector.Protect(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}
	protected.Ciphertext[len(protected.Ciphertext)-1] ^= 1

	if _, err := protector.Unprotect(t.Context(), protected); err == nil {
		t.Fatal("expected tampered ciphertext to fail authentication")
	}
}

func TestAESGCMEmailProtectorLookupHMACIsStable(t *testing.T) {
	protector, email := newTestEmailProtector(t)
	first, err := protector.LookupHMAC(context.Background(), email)
	if err != nil {
		t.Fatal(err)
	}
	second, err := protector.LookupHMAC(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("expected lookup HMAC to be stable")
	}
}

func TestAESGCMEmailProtectorLookupCandidatesPrioritizeCurrentKey(t *testing.T) {
	email, err := domain.ParseEmail("ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	protector := NewAESGCMEmailProtector(testKeyring{
		lookupVersion: 2,
		lookupKey:     []byte("abcdefghijklmnopqrstuvwxyz123456"),
		lookupKeys: map[int][]byte{
			1: []byte("abcdefghijklmnop"),
			2: []byte("abcdefghijklmnopqrstuvwxyz123456"),
		},
	})

	candidates, err := protector.LookupCandidates(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected two lookup candidates, got %d", len(candidates))
	}
	if candidates[0].LookupKeyVersion != 2 || candidates[1].LookupKeyVersion != 1 {
		t.Fatalf("expected current key first, got versions %d and %d", candidates[0].LookupKeyVersion, candidates[1].LookupKeyVersion)
	}
}

func TestAESGCMEmailProtectorReturnsKeyringErrors(t *testing.T) {
	email, err := domain.ParseEmail("ana@example.com")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		key  testKeyring
	}{
		{
			name: "encryption key",
			key:  testKeyring{currentEncryptErr: errKeyringUnavailable},
		},
		{
			name: "lookup key",
			key: testKeyring{
				encryptionVersion: 1,
				encryptionKey:     []byte("01234567890123456789012345678901"),
				currentLookupErr:  errKeyringUnavailable,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAESGCMEmailProtector(tt.key).Protect(t.Context(), email)
			if !errors.Is(err, errKeyringUnavailable) {
				t.Fatalf("expected keyring error, got %v", err)
			}
		})
	}
}

func TestAESGCMEmailProtectorRejectsUnknownEncryptionKeyVersion(t *testing.T) {
	protector, email := newTestEmailProtector(t)
	protected, err := protector.Protect(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}
	protected.EncryptionKeyVersion = 999

	if _, err := protector.Unprotect(t.Context(), protected); !errors.Is(err, errKeyringUnavailable) {
		t.Fatalf("expected keyring error, got %v", err)
	}
}

func TestAESGCMEmailProtectorLookupHMACReturnsKeyringError(t *testing.T) {
	email, err := domain.ParseEmail("ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	protector := NewAESGCMEmailProtector(testKeyring{currentLookupErr: errKeyringUnavailable})

	if _, err := protector.LookupHMAC(t.Context(), email); !errors.Is(err, errKeyringUnavailable) {
		t.Fatalf("expected keyring error, got %v", err)
	}
}

func TestDecryptGCMRejectsShortCiphertext(t *testing.T) {
	if _, err := corecryptography.DecryptGCM([]byte("short"), []byte("01234567890123456789012345678901")); err == nil {
		t.Fatal("expected short ciphertext to fail")
	}
}

func TestGCMRejectsInvalidAESKeys(t *testing.T) {
	invalidKey := []byte("invalid")
	if _, err := corecryptography.EncryptGCM([]byte("ana@example.com"), invalidKey); err == nil {
		t.Fatal("expected encryption with an invalid key to fail")
	}
	if _, err := corecryptography.DecryptGCM([]byte("ciphertext"), invalidKey); err == nil {
		t.Fatal("expected decryption with an invalid key to fail")
	}
}

func TestAESGCMEmailProtectorRejectsInvalidDecryptedEmail(t *testing.T) {
	protector, _ := newTestEmailProtector(t)
	ciphertext, err := corecryptography.EncryptGCM([]byte("not-an-email"), []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = protector.Unprotect(t.Context(), ports.ProtectedEmail{
		Ciphertext:           ciphertext,
		EncryptionKeyVersion: 1,
	})
	if err == nil {
		t.Fatal("expected invalid decrypted email to fail validation")
	}
}

func FuzzDecryptGCMDoesNotPanic(f *testing.F) {
	key := []byte("01234567890123456789012345678901")
	f.Add([]byte{})
	f.Add([]byte("short"))
	f.Add([]byte("not-a-valid-ciphertext"))

	f.Fuzz(func(t *testing.T, ciphertext []byte) {
		_, _ = corecryptography.DecryptGCM(ciphertext, key)
	})
}
