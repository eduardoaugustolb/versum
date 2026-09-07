package ports

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type UserWriter interface {
	CreateUser(ctx context.Context, user *domain.User) error
}

type LoginTokenWriter interface {
	CreateLoginToken(ctx context.Context, token *domain.LoginToken) error
	ConsumeLoginTokenByTokenHash(ctx context.Context, tokenHash []byte, consumedAt *time.Time) error
}

type SessionWriter interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	RevokeSession(ctx context.Context, sessionID string, revokedAt *time.Time) error
	RevokeAllSessions(ctx context.Context, userID string, revokedAt *time.Time) error
}
