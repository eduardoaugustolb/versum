package application

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

// SessionRepository persists and retrieves Session aggregates.
type SessionRepository interface {
	// LockSessionsByUserID serializes rotation and family revocation within a transaction.
	LockSessionsByUserID(ctx context.Context, userID string) error
	SaveUsage(ctx context.Context, session *domain.Session) error
	SaveRotation(ctx context.Context, session *domain.Session) error
	RevokeSessionFamily(ctx context.Context, userID, familyID string, revokedAt time.Time) error
	CreateSession(ctx context.Context, session *domain.Session, secret string) error
	FindSessionByID(ctx context.Context, id string) (*domain.Session, error)
	FindSessionBySecret(ctx context.Context, secret string) (*domain.Session, error)
	ListSessionsByUserID(ctx context.Context, userID string) ([]domain.Session, error)
	RevokeSession(ctx context.Context, sessionID string, revokedAt *time.Time) error
	RevokeAllSessions(ctx context.Context, userID string, revokedAt *time.Time) error
}
