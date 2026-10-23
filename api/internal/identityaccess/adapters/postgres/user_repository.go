package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type UserDAO interface {
	InsertUser(ctx context.Context, id string, email ports.ProtectedEmail) error
	UserByID(ctx context.Context, id string) (UserRecord, error)
	UserByLookupCandidates(ctx context.Context, hashes [][]byte, versions []int16) (UserRecord, error)
	UpdateUserLookup(ctx context.Context, id string, hash []byte, version, previousVersion int) error
	DeleteUser(ctx context.Context, id string) error
}

type UserRepository struct {
	dao       UserDAO
	protector ports.EmailProtector
}

func NewUserRepository(dao UserDAO, protector ports.EmailProtector) *UserRepository {
	return &UserRepository{dao: dao, protector: protector}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	protected, err := r.protector.Protect(ctx, user.Email())
	if err != nil {
		return fmt.Errorf("protecting user email: %w", err)
	}
	if err := r.dao.InsertUser(ctx, user.ID(), protected); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return application.ErrUserAlreadyExists
		}
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidUserID
	}
	record, err := r.dao.UserByID(ctx, id)
	return r.mapUser(ctx, record, err, "finding user by id")
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	candidates, err := r.protector.LookupCandidates(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("creating email lookup candidates: %w", err)
	}
	if len(candidates) == 0 {
		return nil, application.ErrUserNotFound
	}

	lookupHMACs := make([][]byte, 0, len(candidates))
	lookupKeyVersions := make([]int16, 0, len(candidates))
	for _, candidate := range candidates {
		lookupHMACs = append(lookupHMACs, candidate.LookupHMAC)
		lookupKeyVersions = append(lookupKeyVersions, int16(candidate.LookupKeyVersion))
	}

	record, findErr := r.dao.UserByLookupCandidates(ctx, lookupHMACs, lookupKeyVersions)
	user, err := r.mapUser(ctx, record, findErr, "finding user by email")
	protected := record.Email
	if err != nil {
		return nil, err
	}

	currentCandidate := candidates[0]
	if protected.LookupKeyVersion != currentCandidate.LookupKeyVersion {
		if err := r.dao.UpdateUserLookup(ctx, user.ID(), currentCandidate.LookupHMAC, currentCandidate.LookupKeyVersion, protected.LookupKeyVersion); err != nil {
			return nil, fmt.Errorf("migrating user email lookup key: %w", err)
		}
	}

	return user, nil
}

func (r *UserRepository) mapUser(ctx context.Context, record UserRecord, err error, operation string) (*domain.User, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, application.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	email, err := r.protector.Unprotect(ctx, record.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: unprotecting user email: %w", operation, err)
	}
	user, err := domain.RehydrateUser(record.ID, email.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return user, nil
}
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	err := r.dao.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

var _ ports.UserRepository = (*UserRepository)(nil)
