package domain_test

import (
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestParseEmailNormalizesValidInput(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "trims surrounding whitespace",
			raw:  "  ana@example.com  ",
			want: "ana@example.com",
		},
		{
			name: "normalizes address case",
			raw:  "Ana@EXAMPLE.COM",
			want: "ana@example.com",
		},
		{
			name: "accepts common address characters",
			raw:  "ana.silva+versum@example.com",
			want: "ana.silva+versum@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := domain.ParseEmail(tt.raw)
			if err != nil {
				t.Fatal(err)
			}
			if got := email.String(); got != tt.want {
				t.Fatalf("expected normalized email %q, got %q", tt.want, got)
			}
		})
	}
}

func TestParseEmailRejectsInvalidInput(t *testing.T) {
	tests := []string{
		"",
		"   ",
		"ana",
		"@example.com",
		"ana@",
		"ana@@example.com",
		"ana silva@example.com",
		"ana@example com",
	}

	for _, raw := range tests {
		t.Run(raw, func(t *testing.T) {
			_, err := domain.ParseEmail(raw)
			if !errors.Is(err, domain.ErrInvalidEmail) {
				t.Fatalf("expected invalid-email error for %q, got %v", raw, err)
			}
		})
	}
}
