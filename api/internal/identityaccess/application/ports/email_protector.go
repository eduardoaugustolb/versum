package ports

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// ProtectedEmail is the storage representation of an email. It belongs to
// the protection boundary rather than to the domain entity.
type ProtectedEmail struct {
	Ciphertext           []byte
	LookupHMAC           []byte
	EncryptionKeyVersion int
	LookupKeyVersion     int
}

// EmailProtector protects normalized emails before persistence and restores
// them when hydrating a domain user.
type EmailProtector interface {
	Protect(ctx context.Context, email domain.Email) (ProtectedEmail, error)
	Unprotect(ctx context.Context, protected ProtectedEmail) (domain.Email, error)
	LookupHMAC(ctx context.Context, email domain.Email) ([]byte, error)
}
