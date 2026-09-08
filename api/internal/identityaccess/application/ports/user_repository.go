package ports

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// UserRepository persists and retrieves User aggregates.
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
}
