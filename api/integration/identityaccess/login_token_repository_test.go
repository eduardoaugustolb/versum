package identityaccess_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	identityaccesspg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
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

func newLoginToken(t *testing.T, id string, userID string) *domain.LoginToken {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Microsecond)
	token, err := domain.NewLoginToken(id, userID, now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestLoginTokenRepositoryCreatesFindsAndConsumes(t *testing.T) {
	ctx := t.Context()
	repo, db, _, user := setupLoginTokenRepository(ctx, t)
	want := newLoginToken(t, "token-1", user.ID())
	if err := repo.Create(ctx, want, "secret"); err != nil {
		t.Fatal(err)
	}
	var storedHash []byte
	if err := db.QueryRow(ctx, "SELECT token_hash FROM login_tokens WHERE id = $1", want.ID()).Scan(&storedHash); err != nil {
		t.Fatal(err)
	}
	expectedHash := sha256.Sum256([]byte("secret"))
	if !bytes.Equal(storedHash, expectedHash[:]) {
		t.Fatal("database must contain only the secret digest")
	}
	got, err := repo.FindByID(ctx, want.ID())
	if err != nil {
		t.Fatal(err)
	}
	if got.ID() != want.ID() || got.UserID() != want.UserID() || !got.ExpiresAt().Equal(want.ExpiresAt()) {
		t.Fatalf("unexpected login token: %+v", got)
	}
	byHash, err := repo.FindByToken(ctx, "secret")
	if err != nil || byHash.ID() != want.ID() {
		t.Fatalf("unexpected lookup result: %+v, %v", byHash, err)
	}
	consumedAt := time.Now().UTC().Truncate(time.Microsecond)
	if err := want.Consume(consumedAt); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, want); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, want); !errors.Is(err, domain.ErrLoginTokenAlreadyConsumed) {
		t.Fatalf("expected duplicate consumption error, got %v", err)
	}
	consumed, err := repo.FindByID(ctx, want.ID())
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
	if _, err := repo.FindByID(t.Context(), "missing"); !errors.Is(err, application.ErrLoginTokenNotFound) {
		t.Fatalf("expected missing-token error, got %v", err)
	}
	if _, err := repo.FindByToken(t.Context(), "missing"); !errors.Is(err, application.ErrLoginTokenNotFound) {
		t.Fatalf("expected missing-token error, got %v", err)
	}
}
