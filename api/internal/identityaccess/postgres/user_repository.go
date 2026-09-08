package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

type UserRepository struct {
	db        dbexec.Executor
	protector ports.EmailProtector
}

func NewUserRepository(db dbexec.Executor, protector ports.EmailProtector) *UserRepository {
	return &UserRepository{db: db, protector: protector}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	protected, err := r.protector.Protect(ctx, user.Email())
	if err != nil {
		return fmt.Errorf("protecting user email: %w", err)
	}
	if err := r.db.Exec(ctx, CreateUserQuery, user.ID(), protected.Ciphertext, protected.LookupHMAC, protected.EncryptionKeyVersion, protected.LookupKeyVersion); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidUserID
	}
	return r.scanUser(ctx, r.db.QueryRow(ctx, FindUserByIDQuery, id), "finding user by id")
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	lookupHMAC, err := r.protector.LookupHMAC(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("creating email lookup hmac: %w", err)
	}
	return r.scanUser(ctx, r.db.QueryRow(ctx, FindUserByEmailLookupHMACQuery, lookupHMAC), "finding user by email")
}

func (r *UserRepository) scanUser(ctx context.Context, row dbexec.Row, operation string) (*domain.User, error) {
	var id string
	var protected ports.ProtectedEmail
	if err := row.Scan(&id, &protected.Ciphertext, &protected.LookupHMAC, &protected.EncryptionKeyVersion, &protected.LookupKeyVersion); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	email, err := r.protector.Unprotect(ctx, protected)
	if err != nil {
		return nil, fmt.Errorf("%s: unprotecting user email: %w", operation, err)
	}
	user, err := domain.RehydrateUser(id, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return user, nil
}

var _ ports.UserReader = (*UserRepository)(nil)
var _ ports.UserWriter = (*UserRepository)(nil)
