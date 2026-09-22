package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
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
	var createdUserID string
	if err := r.db.QueryRow(ctx, CreateUserQuery, user.ID(), protected.Ciphertext, protected.LookupHMAC, protected.EncryptionKeyVersion, protected.LookupKeyVersion).Scan(&createdUserID); err != nil {
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
	return r.scanUser(ctx, r.db.QueryRow(ctx, FindUserByIDQuery, id), "finding user by id")
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

	user, protected, err := r.scanProtectedUser(
		ctx,
		r.db.QueryRow(ctx, FindUserByLookupCandidatesQuery, lookupHMACs, lookupKeyVersions),
		"finding user by email",
	)
	if err != nil {
		return nil, err
	}

	currentCandidate := candidates[0]
	if protected.LookupKeyVersion != currentCandidate.LookupKeyVersion {
		if err := r.db.Exec(ctx, UpdateUserLookupKeyQuery, user.ID(), currentCandidate.LookupHMAC, currentCandidate.LookupKeyVersion, protected.LookupKeyVersion); err != nil {
			return nil, fmt.Errorf("migrating user email lookup key: %w", err)
		}
	}

	return user, nil
}

func (r *UserRepository) scanUser(ctx context.Context, row dbexec.Row, operation string) (*domain.User, error) {
	user, _, err := r.scanProtectedUser(ctx, row, operation)
	return user, err
}

func (r *UserRepository) scanProtectedUser(ctx context.Context, row dbexec.Row, operation string) (*domain.User, ports.ProtectedEmail, error) {
	var id string
	var protected ports.ProtectedEmail
	if err := row.Scan(&id, &protected.Ciphertext, &protected.LookupHMAC, &protected.EncryptionKeyVersion, &protected.LookupKeyVersion); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ProtectedEmail{}, application.ErrUserNotFound
		}
		return nil, ports.ProtectedEmail{}, fmt.Errorf("%s: %w", operation, err)
	}
	email, err := r.protector.Unprotect(ctx, protected)
	if err != nil {
		return nil, ports.ProtectedEmail{}, fmt.Errorf("%s: unprotecting user email: %w", operation, err)
	}
	user, err := domain.RehydrateUser(id, email.String())
	if err != nil {
		return nil, ports.ProtectedEmail{}, fmt.Errorf("%s: %w", operation, err)
	}
	return user, protected, nil
}

var _ ports.UserRepository = (*UserRepository)(nil)
