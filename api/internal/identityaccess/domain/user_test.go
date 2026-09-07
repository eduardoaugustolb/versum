package domain_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		ciphertext []byte
		lookupHMAC []byte
		encryption int
		lookup     int
		wantErr    error
	}{
		{"valid", "user-1", []byte("ciphertext"), []byte("hmac"), 1, 1, nil},
		{"missing id", "", []byte("ciphertext"), []byte("hmac"), 1, 1, domain.ErrInvalidUserID},
		{"empty ciphertext", "user-1", nil, []byte("hmac"), 1, 1, domain.ErrInvalidEmailCiphertext},
		{"empty lookup hmac", "user-1", []byte("ciphertext"), nil, 1, 1, domain.ErrInvalidEmailLookupHMAC},
		{"invalid encryption version", "user-1", []byte("ciphertext"), []byte("hmac"), 0, 1, domain.ErrInvalidEmailEncryptionKeyVersion},
		{"invalid lookup version", "user-1", []byte("ciphertext"), []byte("hmac"), 1, 0, domain.ErrInvalidEmailLookupKeyVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.id, tt.ciphertext, tt.lookupHMAC, tt.encryption, tt.lookup)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && user == nil {
				t.Fatal("valid user was not created")
			}
		})
	}
}

func TestUserRehydrateAndDefensiveCopies(t *testing.T) {
	ciphertext := []byte("ciphertext")
	lookupHMAC := []byte("hmac")
	user, err := domain.RehydrateUser("user-1", ciphertext, lookupHMAC, 2, 3)
	if err != nil {
		t.Fatalf("rehydration failed: %v", err)
	}

	ciphertext[0] = 'X'
	lookupHMAC[0] = 'Y'
	if bytes.Equal(user.EmailCiphertext(), ciphertext) {
		t.Fatal("user retained the input ciphertext reference")
	}
	if bytes.Equal(user.EmailLookupHMAC(), lookupHMAC) {
		t.Fatal("user retained the input HMAC reference")
	}

	returnedCiphertext := user.EmailCiphertext()
	returnedCiphertext[0] = 'Z'
	if bytes.Equal(user.EmailCiphertext(), returnedCiphertext) {
		t.Fatal("getter exposed the internal ciphertext reference")
	}
	if user.EmailEncryptionKeyVersion() != 2 || user.EmailLookupKeyVersion() != 3 {
		t.Fatal("key versions were rehydrated incorrectly")
	}
}
