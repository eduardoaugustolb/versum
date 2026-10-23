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

type SessionDAO interface {
	InsertSession(ctx context.Context, id string, hash []byte, userID string, revokedAt, lastUsedAt any, expiresAt time.Time, familyID, ip, ua string) error
	SessionByID(ctx context.Context, id string) (SessionRecord, error)
	SessionByHash(ctx context.Context, hash []byte) (SessionRecord, error)
	RevokeSession(ctx context.Context, id string, at time.Time) error
	RevokeAllSessions(ctx context.Context, userID string, at time.Time) error
	LockUser(ctx context.Context, userID string) error
	RotateSession(ctx context.Context, id string, revokedAt, lastUsedAt time.Time, replacement string) error
	RevokeSessionFamily(ctx context.Context, userID, familyID string, at time.Time) error
	UpdateSessionUsage(ctx context.Context, id string, at time.Time) error
	SessionsByUserID(ctx context.Context, userID string) ([]SessionRecord, error)
}

type SessionRepository struct{ dao SessionDAO }

func NewSessionRepository(dao SessionDAO) *SessionRepository {
	return &SessionRepository{dao: dao}
}

var _ ports.SessionRepository = (*SessionRepository)(nil)

func (r *SessionRepository) CreateSession(ctx context.Context, session *domain.Session, secret string) error {
	if secret == "" {
		return application.ErrInvalidSessionSecret
	}
	hash := sha256.Sum256([]byte(secret))
	if _, replaced := session.ReplacedBySessionID(); replaced {
		return domain.ErrInvalidSessionReplacement
	}
	revokedAt, _ := session.RevokedAt()
	lastUsedAt, _ := session.LastUsedAt()
	var revokedAtValue, lastUsedAtValue any
	if !revokedAt.IsZero() {
		revokedAtValue = revokedAt
	}
	if !lastUsedAt.IsZero() {
		lastUsedAtValue = lastUsedAt
	}
	if err := r.dao.InsertSession(ctx, session.ID(), hash[:], session.UserID(), revokedAtValue, lastUsedAtValue, session.ExpiresAt(), session.FamilyID(), session.IPAddress(), session.UserAgent()); err != nil {
		if isUniqueViolation(err) {
			return application.ErrSessionAlreadyExists
		}
		return fmt.Errorf("creating session: %w", err)
	}
	return nil
}

func (r *SessionRepository) FindSessionByID(ctx context.Context, id string) (*domain.Session, error) {
	if id == "" {
		return nil, domain.ErrInvalidSessionID
	}
	record, err := r.dao.SessionByID(ctx, id)
	return mapSession(record, err, "finding session by id")
}

func (r *SessionRepository) FindSessionBySecret(ctx context.Context, secret string) (*domain.Session, error) {
	if secret == "" {
		return nil, application.ErrInvalidSessionSecret
	}
	hash := sha256.Sum256([]byte(secret))
	record, err := r.dao.SessionByHash(ctx, hash[:])
	return mapSession(record, err, "finding session by secret")
}

func (r *SessionRepository) ListSessionsByUserID(ctx context.Context, userID string) ([]domain.Session, error) {
	if userID == "" {
		return nil, domain.ErrInvalidSessionUserID
	}
	records, err := r.dao.SessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing sessions by user id: %w", err)
	}
	sessions := make([]domain.Session, 0, len(records))
	for _, record := range records {
		session, err := mapSession(record, nil, "mapping listed session")
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}
	return sessions, nil
}

func (r *SessionRepository) RevokeSession(ctx context.Context, sessionID string, revokedAt *time.Time) error {
	if sessionID == "" {
		return domain.ErrInvalidSessionID
	}
	if revokedAt == nil || revokedAt.IsZero() {
		return domain.ErrInvalidSessionRevokedAt
	}
	if err := r.dao.RevokeSession(ctx, sessionID, *revokedAt); err != nil {
		return fmt.Errorf("revoking session: %w", err)
	}
	return nil
}

func (r *SessionRepository) RevokeAllSessions(ctx context.Context, userID string, revokedAt *time.Time) error {
	if userID == "" {
		return domain.ErrInvalidSessionUserID
	}
	if revokedAt == nil || revokedAt.IsZero() {
		return domain.ErrInvalidSessionRevokedAt
	}
	if err := r.dao.RevokeAllSessions(ctx, userID, *revokedAt); err != nil {
		return fmt.Errorf("revoking all sessions: %w", err)
	}
	return nil
}

func mapSession(record SessionRecord, err error, operation string) (*domain.Session, error) {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	replacedBy := ""
	if record.Replacement != nil {
		replacedBy = *record.Replacement
	}
	session, err := domain.RehydrateSession(record.ID, record.UserID, record.FamilyID, record.IPAddress, record.UserAgent, record.RevokedAt, record.LastUsedAt, record.ExpiresAt, replacedBy)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return session, nil
}

func (r *SessionRepository) LockSessionsByUserID(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.ErrInvalidSessionUserID
	}
	err := r.dao.LockUser(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return application.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("locking user sessions: %w", err)
	}
	return nil
}

func (r *SessionRepository) SaveRotation(ctx context.Context, session *domain.Session) error {
	replacement, ok := session.ReplacedBySessionID()
	if !ok {
		return domain.ErrInvalidSessionReplacement
	}
	revokedAt, _ := session.RevokedAt()
	lastUsedAt, _ := session.LastUsedAt()
	err := r.dao.RotateSession(ctx, session.ID(), revokedAt, lastUsedAt, replacement)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrInvalidSessionState
	}
	if err != nil {
		return fmt.Errorf("saving session rotation: %w", err)
	}
	return nil
}

func (r *SessionRepository) RevokeSessionFamily(ctx context.Context, userID, familyID string, revokedAt time.Time) error {
	if userID == "" {
		return domain.ErrInvalidSessionUserID
	}
	if familyID == "" {
		return domain.ErrInvalidSessionFamilyID
	}
	if revokedAt.IsZero() {
		return domain.ErrInvalidSessionRevokedAt
	}
	if err := r.dao.RevokeSessionFamily(ctx, userID, familyID, revokedAt.UTC()); err != nil {
		return fmt.Errorf("revoking session family: %w", err)
	}
	return nil
}

func (r *SessionRepository) SaveUsage(ctx context.Context, session *domain.Session) error {
	lastUsedAt, ok := session.LastUsedAt()
	if !ok {
		return domain.ErrInvalidSessionLastUsedAt
	}
	if !session.IsValidAt(lastUsedAt) {
		return domain.ErrInvalidSessionState
	}
	err := r.dao.UpdateSessionUsage(ctx, session.ID(), lastUsedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrInvalidSessionState
	}
	if err != nil {
		return fmt.Errorf("saving session usage: %w", err)
	}
	return nil
}
