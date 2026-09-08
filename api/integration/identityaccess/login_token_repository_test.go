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

func setupLoginTokenRepository(ctx context.Context, t *testing.T) (*identityaccesspg.LoginTokenRepository, dbexec.Executor, *pgxpool.Pool, *domain.User) {
	t.Helper()
	db, pool, err := setupPostgresDBExecutor(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := db.Exec(ctx, "DELETE FROM users WHERE id = $1", "token-user"); err != nil {
		t.Fatal(err)
	}
	user := createTestUser(ctx, t, db, "token-user", "token-user@example.com")
	t.Cleanup(func() { _ = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", user.ID()) })
	return identityaccesspg.NewLoginTokenRepository(db), db, pool, user
}

func newLoginToken(t *testing.T, id string, hash []byte, userID string) *domain.LoginToken {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Microsecond)
	token, err := domain.NewLoginToken(id, hash, userID, now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestLoginTokenRepositoryCreatesFindsAndConsumes(t *testing.T) {
	ctx := t.Context()
	repo, _, _, user := setupLoginTokenRepository(ctx, t)
	want := newLoginToken(t, "token-1", []byte("token-hash-1"), user.ID())
	if err := repo.CreateLoginToken(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindLoginTokenByID(ctx, want.ID())
	if err != nil {
		t.Fatal(err)
	}
	if got.ID() != want.ID() || got.UserID() != want.UserID() || string(got.TokenHash()) != string(want.TokenHash()) || !got.ExpiresAt().Equal(want.ExpiresAt()) {
		t.Fatalf("unexpected login token: %+v", got)
	}
	byHash, err := repo.FindLoginTokenByTokenHash(ctx, want.TokenHash())
	if err != nil || byHash.ID() != want.ID() {
		t.Fatalf("unexpected lookup result: %+v, %v", byHash, err)
	}
	consumedAt := time.Now().UTC().Truncate(time.Microsecond)
	if err := repo.ConsumeLoginTokenByTokenHash(ctx, want.TokenHash(), &consumedAt); err != nil {
		t.Fatal(err)
	}
	consumed, err := repo.FindLoginTokenByID(ctx, want.ID())
	if err != nil {
		t.Fatal(err)
	}
	at, ok := consumed.ConsumedAt()
	if !ok || !at.Equal(consumedAt) {
		t.Fatalf("expected consumed_at %v, got %v (present: %t)", consumedAt, at, ok)
	}
}

func TestLoginTokenRepositoryReturnsNotFound(t *testing.T) {
	repo, _, _, _ := setupLoginTokenRepository(t.Context(), t)
	if _, err := repo.FindLoginTokenByID(t.Context(), "missing"); !errors.Is(err, application.ErrLoginTokenNotFound) {
		t.Fatalf("expected missing-token error, got %v", err)
	}
	if _, err := repo.FindLoginTokenByTokenHash(t.Context(), []byte("missing")); !errors.Is(err, application.ErrLoginTokenNotFound) {
		t.Fatalf("expected missing-token error, got %v", err)
	}
}
