package ports

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// SessionRepository persists and retrieves Session aggregates.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	FindSessionByID(ctx context.Context, id string) (*domain.Session, error)
	FindSessionBySecretHash(ctx context.Context, secretHash []byte) (*domain.Session, error)
	ListSessionsByUserID(ctx context.Context, userID string) ([]domain.Session, error)
	RevokeSession(ctx context.Context, sessionID string, revokedAt *time.Time) error
	RevokeAllSessions(ctx context.Context, userID string, revokedAt *time.Time) error
}
