package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

type LoginTokenRepository struct{ db dbexec.Executor }

func NewLoginTokenRepository(db dbexec.Executor) *LoginTokenRepository {
	return &LoginTokenRepository{db: db}
}

func (r *LoginTokenRepository) CreateLoginToken(ctx context.Context, token *domain.LoginToken) error {
	consumedAt, _ := token.ConsumedAt()
	var consumedAtValue any
	if !consumedAt.IsZero() {
		consumedAtValue = consumedAt
	}
	if err := r.db.Exec(ctx, CreateLoginTokenQuery, token.ID(), token.TokenHash(), token.UserID(), token.ExpiresAt(), consumedAtValue); err != nil {
		return fmt.Errorf("creating login token: %w", err)
	}
	return nil
}

func (r *LoginTokenRepository) FindLoginTokenByID(ctx context.Context, id string) (*domain.LoginToken, error) {
	if id == "" {
		return nil, domain.ErrInvalidLoginTokenID
	}
	return scanLoginToken(r.db.QueryRow(ctx, FindLoginTokenByIDQuery, id), "finding login token by id")
}

func (r *LoginTokenRepository) FindLoginTokenByTokenHash(ctx context.Context, tokenHash []byte) (*domain.LoginToken, error) {
	if len(tokenHash) == 0 {
		return nil, domain.ErrInvalidLoginTokenHash
	}
	return scanLoginToken(r.db.QueryRow(ctx, FindLoginTokenByTokenHashQuery, tokenHash), "finding login token by token hash")
}

func (r *LoginTokenRepository) ConsumeLoginTokenByTokenHash(ctx context.Context, tokenHash []byte, consumedAt *time.Time) error {
	if len(tokenHash) == 0 {
		return domain.ErrInvalidLoginTokenHash
	}
	if consumedAt == nil || consumedAt.IsZero() {
		return domain.ErrInvalidLoginTokenConsumedAt
	}
	if err := r.db.Exec(ctx, ConsumeLoginTokenByTokenHashQuery, tokenHash, *consumedAt); err != nil {
		return fmt.Errorf("consuming login token: %w", err)
	}
	return nil
}

func scanLoginToken(row dbexec.Row, operation string) (*domain.LoginToken, error) {
	var id, userID string
	var tokenHash []byte
	var expiresAt time.Time
	var consumedAt *time.Time
	if err := row.Scan(&id, &tokenHash, &userID, &expiresAt, &consumedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrLoginTokenNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	token, err := domain.RehydrateLoginToken(id, tokenHash, userID, expiresAt, consumedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return token, nil
}
