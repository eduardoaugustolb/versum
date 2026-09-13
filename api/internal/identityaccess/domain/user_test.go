package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestNewUser(t *testing.T) {
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
			user, err := domain.NewUser(tt.id, "ana@example.com")
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
	emailRaw := "Ana@EXAMPLE.COM"
	user, err := domain.RehydrateUser("user-1", emailRaw)
	if err != nil {
		t.Fatalf("rehydration failed: %v", err)
	}
	if got := user.Email().String(); got != strings.ToLower(emailRaw) {
		t.Fatalf("expected normalized email, got %q", got)
	}
}
