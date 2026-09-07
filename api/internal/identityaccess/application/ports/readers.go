package ports

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type UserReader interface {
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByEmailLookupHMAC(ctx context.Context, emailLookupHMAC []byte) (*domain.User, error)
}

type SessionReader interface {
	FindSessionByID(ctx context.Context, id string) (*domain.Session, error)
	FindSessionBySecretHash(ctx context.Context, secretHash []byte) (*domain.Session, error)
	ListSessionsByUserID(ctx context.Context, userID string) ([]domain.Session, error)
}

type LoginTokenReader interface {
	FindLoginTokenByID(ctx context.Context, id string) (*domain.LoginToken, error)
	FindLoginTokenByTokenHash(ctx context.Context, tokenHash []byte) (*domain.LoginToken, error)
}
