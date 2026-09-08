package ports

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// LoginTokenRepository persists and retrieves LoginToken aggregates.
type LoginTokenRepository interface {
	CreateLoginToken(ctx context.Context, token *domain.LoginToken) error
	FindLoginTokenByID(ctx context.Context, id string) (*domain.LoginToken, error)
	FindLoginTokenByTokenHash(ctx context.Context, tokenHash []byte) (*domain.LoginToken, error)
	ConsumeLoginTokenByTokenHash(ctx context.Context, tokenHash []byte, consumedAt *time.Time) error
}
