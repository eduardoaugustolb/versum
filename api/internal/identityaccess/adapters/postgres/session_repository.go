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

type SessionRepository struct{ db dbexec.Executor }

func NewSessionRepository(db dbexec.Executor) *SessionRepository {
	return &SessionRepository{db: db}
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
	if err := r.db.Exec(ctx, CreateSessionQuery, session.ID(), hash[:], session.UserID(), revokedAtValue, lastUsedAtValue, session.ExpiresAt(), session.FamilyID(), session.IPAddress(), session.UserAgent()); err != nil {
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
	return scanSession(r.db.QueryRow(ctx, FindSessionByIDQuery, id), "finding session by id")
}

func (r *SessionRepository) FindSessionBySecret(ctx context.Context, secret string) (*domain.Session, error) {
	if secret == "" {
		return nil, application.ErrInvalidSessionSecret
	}
	hash := sha256.Sum256([]byte(secret))
	return scanSession(r.db.QueryRow(ctx, FindSessionBySecretHashQuery, hash[:]), "finding session by secret")
}

func (r *SessionRepository) ListSessionsByUserID(ctx context.Context, userID string) ([]domain.Session, error) {
	if userID == "" {
		return nil, domain.ErrInvalidSessionUserID
	}
	rows, err := r.db.Query(ctx, ListSessionsByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("listing sessions by user id: %w", err)
	}
	defer rows.Close()

	sessions := []domain.Session{}
	for rows.Next() {
		session, err := scanSession(rows, "scanning listed session")
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating sessions by user id: %w", err)
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
	if err := r.db.Exec(ctx, RevokeSessionQuery, sessionID, *revokedAt); err != nil {
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
	if err := r.db.Exec(ctx, RevokeAllSessionsQuery, userID, *revokedAt); err != nil {
		return fmt.Errorf("revoking all sessions: %w", err)
	}
	return nil
}

func scanSession(row dbexec.Row, operation string) (*domain.Session, error) {
	var id, userID, familyID, ipAddress, userAgent string
	var replacement *string
	var revokedAt, lastUsedAt *time.Time
	var expiresAt time.Time
	if err := row.Scan(&id, &userID, &revokedAt, &lastUsedAt, &expiresAt, &familyID, &ipAddress, &userAgent, &replacement); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	replacedBy := ""
	if replacement != nil {
		replacedBy = *replacement
	}
	session, err := domain.RehydrateSession(id, userID, familyID, domain.SessionClient{IPAddress: ipAddress, UserAgent: userAgent}, revokedAt, lastUsedAt, expiresAt, replacedBy)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return session, nil
}

func (r *SessionRepository) LockSessionsByUserID(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.ErrInvalidSessionUserID
	}
	var id string
	err := r.db.QueryRow(ctx, LockSessionsByUserIDQuery, userID).Scan(&id)
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
	var id string
	err := r.db.QueryRow(ctx, SaveSessionRotationQuery, session.ID(), revokedAt, lastUsedAt, replacement).Scan(&id)
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
	if err := r.db.Exec(ctx, RevokeSessionFamilyQuery, userID, familyID, revokedAt.UTC()); err != nil {
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
	var id string
	err := r.db.QueryRow(ctx, SaveSessionUsageQuery, session.ID(), lastUsedAt).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrInvalidSessionState
	}
	if err != nil {
		return fmt.Errorf("saving session usage: %w", err)
	}
	return nil
}
