package postgres

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type LoginTokenRepository struct{ db dbexec.Executor }

func NewLoginTokenRepository(db dbexec.Executor) *LoginTokenRepository {
	return &LoginTokenRepository{db: db}
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
	if err := r.db.Exec(ctx, CreateLoginTokenQuery, token.ID(), hash[:], token.UserID(), token.ExpiresAt(), consumedAtValue); err != nil {
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
	return scanLoginToken(r.db.QueryRow(ctx, FindLoginTokenByIDQuery, id), "finding login token by id")
}

func (r *LoginTokenRepository) FindByToken(ctx context.Context, secret string) (*domain.LoginToken, error) {
	if secret == "" {
		return nil, application.ErrInvalidLoginTokenSecret
	}
	hash := sha256.Sum256([]byte(secret))
	return scanLoginToken(r.db.QueryRow(ctx, FindLoginTokenByTokenHashQuery, hash[:]), "finding login token by token")
}

func (r *LoginTokenRepository) Save(ctx context.Context, token *domain.LoginToken) error {
	consumedAt, ok := token.ConsumedAt()
	if !ok {
		return domain.ErrInvalidLoginTokenConsumedAt
	}
	var id string
	err := r.db.QueryRow(ctx, ConsumeLoginTokenByIDQuery, token.ID(), consumedAt).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrLoginTokenAlreadyConsumed
	}
	if err != nil {
		return fmt.Errorf("saving login token consumption: %w", err)
	}
	return nil
}

func scanLoginToken(row dbexec.Row, operation string) (*domain.LoginToken, error) {
	var id, userID string
	var expiresAt time.Time
	var consumedAt *time.Time
	if err := row.Scan(&id, &userID, &expiresAt, &consumedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrLoginTokenNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	token, err := domain.RehydrateLoginToken(id, userID, expiresAt, consumedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return token, nil
}
