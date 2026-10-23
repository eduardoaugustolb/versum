package postgres

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type LoginTokenDAO interface {
	InsertLoginToken(ctx context.Context, id string, hash []byte, userID string, expiresAt time.Time, consumedAt any) error
	LoginTokenByID(ctx context.Context, id string) (LoginTokenRecord, error)
	LoginTokenByHash(ctx context.Context, hash []byte) (LoginTokenRecord, error)
	ConsumeLoginToken(ctx context.Context, id string, at time.Time) error
}

type LoginTokenRepository struct{ dao LoginTokenDAO }

func NewLoginTokenRepository(dao LoginTokenDAO) *LoginTokenRepository {
	return &LoginTokenRepository{dao: dao}
}

var _ ports.LoginTokenRepository = (*LoginTokenRepository)(nil)

func (r *LoginTokenRepository) Create(ctx context.Context, token *domain.LoginToken, secret string) error {
	if secret == "" {
		return application.ErrInvalidLoginTokenSecret
	}
	hash := sha256.Sum256([]byte(secret))
	consumedAt, _ := token.ConsumedAt()
	var consumedAtValue any
	if !consumedAt.IsZero() {
		consumedAtValue = consumedAt
	}
	if err := r.dao.InsertLoginToken(ctx, token.ID(), hash[:], token.UserID(), token.ExpiresAt(), consumedAtValue); err != nil {
		if isUniqueViolation(err) {
			return application.ErrLoginTokenAlreadyExists
		}
		return fmt.Errorf("creating login token: %w", err)
	}
	return nil
}

func (r *LoginTokenRepository) FindByID(ctx context.Context, id string) (*domain.LoginToken, error) {
	if id == "" {
		return nil, domain.ErrInvalidLoginTokenID
	}
	record, err := r.dao.LoginTokenByID(ctx, id)
	return mapLoginToken(record, err, "finding login token by id")
}

func (r *LoginTokenRepository) FindByToken(ctx context.Context, secret string) (*domain.LoginToken, error) {
	if secret == "" {
		return nil, application.ErrInvalidLoginTokenSecret
	}
	hash := sha256.Sum256([]byte(secret))
	record, err := r.dao.LoginTokenByHash(ctx, hash[:])
	return mapLoginToken(record, err, "finding login token by token")
}

func (r *LoginTokenRepository) Save(ctx context.Context, token *domain.LoginToken) error {
	consumedAt, ok := token.ConsumedAt()
	if !ok {
		return domain.ErrInvalidLoginTokenConsumedAt
	}
	err := r.dao.ConsumeLoginToken(ctx, token.ID(), consumedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrLoginTokenAlreadyConsumed
	}
	if err != nil {
		return fmt.Errorf("saving login token consumption: %w", err)
	}
	return nil
}

func mapLoginToken(record LoginTokenRecord, err error, operation string) (*domain.LoginToken, error) {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrLoginTokenNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	token, err := domain.RehydrateLoginToken(record.ID, record.UserID, record.ExpiresAt, record.ConsumedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return token, nil
}
