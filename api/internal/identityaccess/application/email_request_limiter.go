package application

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// EmailRequestLimiter enforces the issuance quota across all inbound adapters.
type EmailRequestLimiter interface {
	Allow(context.Context, domain.Email) (bool, error)
}
