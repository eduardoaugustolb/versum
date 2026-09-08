package domain_test

import (
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func mustEmail(t *testing.T, raw string) domain.Email {
	t.Helper()
	email, err := domain.ParseEmail(raw)
	if err != nil {
		t.Fatal(err)
	}
	return email
}

func TestNewUser(t *testing.T) {
	email := mustEmail(t, "ana@example.com")
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"valid", "user-1", nil},
		{"missing id", "", domain.ErrInvalidUserID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.id, email)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && user == nil {
				t.Fatal("valid user was not created")
			}
		})
	}
}

func TestUserRehydratePreservesEmail(t *testing.T) {
	email := mustEmail(t, "Ana@EXAMPLE.COM")
	user, err := domain.RehydrateUser("user-1", email)
	if err != nil {
		t.Fatalf("rehydration failed: %v", err)
	}
	if got := user.Email().String(); got != "ana@example.com" {
		t.Fatalf("expected normalized email, got %q", got)
	}
}
