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

type SessionRepository struct{ db dbexec.Executor }

func NewSessionRepository(db dbexec.Executor) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	revokedAt, _ := session.RevokedAt()
	usedAt, _ := session.UsedAt()
	var revokedAtValue, usedAtValue any
	if !revokedAt.IsZero() {
		revokedAtValue = revokedAt
	}
	if !usedAt.IsZero() {
		usedAtValue = usedAt
	}
	if err := r.db.Exec(ctx, CreateSessionQuery, session.ID(), session.SecretHash(), session.UserID(), revokedAtValue, usedAtValue, session.ExpiresAt()); err != nil {
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

func (r *SessionRepository) FindSessionBySecretHash(ctx context.Context, secretHash []byte) (*domain.Session, error) {
	if len(secretHash) == 0 {
		return nil, domain.ErrInvalidSessionSecretHash
	}
	return scanSession(r.db.QueryRow(ctx, FindSessionBySecretHashQuery, secretHash), "finding session by secret hash")
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
	var id, userID string
	var secretHash []byte
	var revokedAt, usedAt *time.Time
	var expiresAt time.Time
	if err := row.Scan(&id, &secretHash, &userID, &revokedAt, &usedAt, &expiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	session, err := domain.RehydrateSession(id, secretHash, userID, revokedAt, usedAt, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return session, nil
}
