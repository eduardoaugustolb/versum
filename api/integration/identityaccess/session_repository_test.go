package identityaccess_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	identityaccesspg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupSessionRepository(ctx context.Context, t *testing.T) (*identityaccesspg.SessionRepository, dbexec.Executor, *pgxpool.Pool, *domain.User) {
	t.Helper()
	db, pool, err := setupPostgresDBExecutor(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := db.Exec(ctx, "DELETE FROM users WHERE id = $1", "session-user"); err != nil {
		t.Fatal(err)
	}
	user := createTestUser(ctx, t, db, "session-user", "session-user@example.com")
	t.Cleanup(func() { _ = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", user.ID()) })
	return identityaccesspg.NewSessionRepository(db), db, pool, user
}

func newSession(t *testing.T, id string, hash []byte, userID string) *domain.Session {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Microsecond)
	session, err := domain.NewSession(id, hash, userID, now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func TestSessionRepositoryCreatesFindsListsAndRevokes(t *testing.T) {
	ctx := t.Context()
	repo, _, _, user := setupSessionRepository(ctx, t)
	first := newSession(t, "session-1", []byte("session-hash-1"), user.ID())
	second := newSession(t, "session-2", []byte("session-hash-2"), user.ID())
	for _, session := range []*domain.Session{first, second} {
		if err := repo.CreateSession(ctx, session); err != nil {
			t.Fatal(err)
		}
	}
	got, err := repo.FindSessionByID(ctx, first.ID())
	if err != nil {
		t.Fatal(err)
	}
	if got.ID() != first.ID() || got.UserID() != first.UserID() || string(got.SecretHash()) != string(first.SecretHash()) || !got.ExpiresAt().Equal(first.ExpiresAt()) {
		t.Fatalf("unexpected session: %+v", got)
	}
	byHash, err := repo.FindSessionBySecretHash(ctx, first.SecretHash())
	if err != nil || byHash.ID() != first.ID() {
		t.Fatalf("unexpected lookup result: %+v, %v", byHash, err)
	}
	sessions, err := repo.ListSessionsByUserID(ctx, user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 || sessions[0].ID() != first.ID() || sessions[1].ID() != second.ID() {
		t.Fatalf("unexpected sessions: %+v", sessions)
	}
	revokedAt := time.Now().UTC().Truncate(time.Microsecond)
	if err := repo.RevokeSession(ctx, first.ID(), &revokedAt); err != nil {
		t.Fatal(err)
	}
	if err := repo.RevokeAllSessions(ctx, user.ID(), &revokedAt); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{first.ID(), second.ID()} {
		session, err := repo.FindSessionByID(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		at, ok := session.RevokedAt()
		if !ok || !at.Equal(revokedAt) {
			t.Fatalf("session %q has revoked_at %v (present: %t)", id, at, ok)
		}
	}
}

func TestSessionRepositoryReturnsNotFound(t *testing.T) {
	repo, _, _, _ := setupSessionRepository(t.Context(), t)
	if _, err := repo.FindSessionByID(t.Context(), "missing"); !errors.Is(err, application.ErrSessionNotFound) {
		t.Fatalf("expected missing-session error, got %v", err)
	}
	if _, err := repo.FindSessionBySecretHash(t.Context(), []byte("missing")); !errors.Is(err, application.ErrSessionNotFound) {
		t.Fatalf("expected missing-session error, got %v", err)
	}
}
