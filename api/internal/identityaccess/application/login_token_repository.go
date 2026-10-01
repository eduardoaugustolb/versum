package application

import (
	"context"
	"errors"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// LoginTokenRepository persists and retrieves LoginToken aggregates.
type LoginTokenRepository interface {
	Create(ctx context.Context, token *domain.LoginToken, secret string) error
	FindByID(ctx context.Context, id string) (*domain.LoginToken, error)
	FindByToken(ctx context.Context, secret string) (*domain.LoginToken, error)
	// Save persists consumption and rejects an already consumed token.
	Save(ctx context.Context, token *domain.LoginToken) error
}

var ErrInvalidLoginTokenSecret = errors.New("invalid login token secret")
