package identityaccess_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/cryptography"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"testing"
	"time"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	identityaccesspg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
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

func newSession(t *testing.T, id string, userID string) *domain.Session {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Microsecond)
	session, err := domain.NewSession(id, userID, "family-1", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func TestSessionRepositoryCreatesFindsListsAndRevokes(t *testing.T) {
	ctx := t.Context()
	repo, db, _, user := setupSessionRepository(ctx, t)
	first := newSession(t, "session-1", user.ID())
	second := newSession(t, "session-2", user.ID())
	for _, session := range []*domain.Session{first, second} {
		if err := repo.CreateSession(ctx, session, "secret-"+session.ID()); err != nil {
			t.Fatal(err)
		}
	}
	var storedHash []byte
	if err := db.QueryRow(ctx, "SELECT secret_hash FROM sessions WHERE id = $1", first.ID()).Scan(&storedHash); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("secret-" + first.ID()))
	if !bytes.Equal(storedHash, digest[:]) {
		t.Fatal("database must store the digest, not the session secret")
	}
	got, err := repo.FindSessionByID(ctx, first.ID())
	if err != nil {
		t.Fatal(err)
	}
	if got.ID() != first.ID() || got.UserID() != first.UserID() || !got.ExpiresAt().Equal(first.ExpiresAt()) {
		t.Fatalf("unexpected session: %+v", got)
	}
	bySecret, err := repo.FindSessionBySecret(ctx, "secret-"+first.ID())
	if err != nil || bySecret.ID() != first.ID() {
		t.Fatalf("unexpected lookup result: %+v, %v", bySecret, err)
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
	if _, err := repo.FindSessionBySecret(t.Context(), "missing"); !errors.Is(err, application.ErrSessionNotFound) {
		t.Fatalf("expected missing-session error, got %v", err)
	}
}

func TestSessionRepositoryConcurrentRotationRevokesFamily(t *testing.T) {
	ctx := t.Context()
	repo, db, _, user := setupSessionRepository(ctx, t)
	now := time.Now().UTC().Truncate(time.Microsecond)
	original, err := domain.NewSession("rotation-original", user.ID(), "rotation-family", "192.0.2.1", "browser", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	independent, err := domain.NewSession("rotation-independent", user.ID(), "independent-family", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateSession(ctx, original, "rotation-secret"); err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateSession(ctx, independent, "independent-secret"); err != nil {
		t.Fatal(err)
	}
	uow := identityaccesspg.NewUnitOfWork(db, testEmailProtector{}, nil)
	uc := commands.NewRotateSession(uow, cryptography.RandomTokenGenerator{}, id.UUIDGenerator{}, clock.SystemClock{})
	start := make(chan struct{})
	results := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, err := uc.Execute(ctx, "rotation-secret", "2001:db8::1", "browser")
			results <- err
		}()
	}
	close(start)
	first, second := <-results, <-results
	if !(first == nil && errors.Is(second, domain.ErrSessionReused) || second == nil && errors.Is(first, domain.ErrSessionReused)) {
		t.Fatalf("expected one rotation and one replay: %v, %v", first, second)
	}
	sessions, err := repo.ListSessionsByUserID(ctx, user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 3 {
		t.Fatalf("expected exactly one successor, got %d sessions", len(sessions))
	}
	for _, s := range sessions {
		if s.FamilyID() == "rotation-family" && !s.IsRevoked() {
			t.Fatal("replay must revoke the entire family")
		}
		if s.FamilyID() == "independent-family" && s.IsRevoked() {
			t.Fatal("unrelated family must remain active")
		}
	}
	old, err := repo.FindSessionByID(ctx, original.ID())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := old.ReplacedBySessionID(); !ok {
		t.Fatal("replacement must survive rehydration")
	}
}

func TestSessionRepositoryUsageCannotRewindOrReviveSession(t *testing.T) {
	ctx := t.Context()
	repo, _, _, user := setupSessionRepository(ctx, t)
	session := newSession(t, "usage-session", user.ID())
	if err := repo.CreateSession(ctx, session, "usage-secret"); err != nil {
		t.Fatal(err)
	}
	stale, err := repo.FindSessionByID(ctx, session.ID())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Microsecond)
	if err := session.Use(now); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveUsage(ctx, session); err != nil {
		t.Fatal(err)
	}
	if err := stale.Use(now.Add(-time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveUsage(ctx, stale); !errors.Is(err, domain.ErrInvalidSessionState) {
		t.Fatalf("expected rejection of stale usage: %v", err)
	}
	got, err := repo.FindSessionByID(ctx, session.ID())
	if err != nil {
		t.Fatal(err)
	}
	if at, ok := got.LastUsedAt(); !ok || !at.Equal(now) {
		t.Fatal("last use was not persisted monotonically")
	}
	if err := repo.RevokeSession(ctx, session.ID(), &now); err != nil {
		t.Fatal(err)
	}
	if err := session.Use(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveUsage(ctx, session); !errors.Is(err, domain.ErrInvalidSessionState) {
		t.Fatalf("stale state must not update a revoked session: %v", err)
	}
}
