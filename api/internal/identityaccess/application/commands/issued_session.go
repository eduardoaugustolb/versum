package commands

import "github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"

// IssuedSession contains a bearer credential for delivery to the client.
// Secret must never be logged or stored in plaintext.
type IssuedSession struct {
	Session *domain.Session
	Secret  string
}
