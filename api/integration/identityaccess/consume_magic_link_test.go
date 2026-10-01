package identityaccess_test

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/cryptography"
	identitypg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type fixedSessionID struct{ value string }

func (g fixedSessionID) Generate() string { return g.value }

func TestConsumeMagicLinkPersistsSessionAndRejectsReuse(t *testing.T) {
	ctx := t.Context()
	sessions, db, _, user := setupSessionRepository(ctx, t)
	tokens := identitypg.NewLoginTokenRepository(db)
	token := newLoginToken(t, "consume-link", user.ID())
	if err := tokens.Create(ctx, token, "link-secret"); err != nil {
		t.Fatal(err)
	}
	uc, err := commands.NewConsumeMagicLink(identitypg.NewUnitOfWork(db, testEmailProtector{}, nil), cryptography.RandomTokenGenerator{}, id.UUIDGenerator{}, clock.SystemClock{}, 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	got, err := uc.Execute(ctx, "link-secret", domain.SessionClient{IPAddress: "2001:db8::1", UserAgent: "integration-browser"})
	if err != nil {
		t.Fatal(err)
	}
	persisted, err := sessions.FindSessionBySecret(ctx, got.Secret)
	if err != nil {
		t.Fatal(err)
	}
	if persisted.ID() != got.Session.ID() || persisted.FamilyID() != got.Session.FamilyID() || persisted.UserID() != user.ID() || persisted.IPAddress() != "2001:db8::1" || persisted.UserAgent() != "integration-browser" || !persisted.IsValidAt(time.Now()) {
		t.Fatal("session was not fully persisted")
	}
	var digest []byte
	if err := db.QueryRow(ctx, "SELECT secret_hash FROM sessions WHERE id = $1", persisted.ID()).Scan(&digest); err != nil {
		t.Fatal(err)
	}
	want := sha256.Sum256([]byte(got.Secret))
	if !bytes.Equal(digest, want[:]) {
		t.Fatal("session secret must only be stored as a digest")
	}
	consumed, err := tokens.FindByID(ctx, token.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !consumed.IsConsumed() {
		t.Fatal("token consumption was not committed")
	}
	if result, err := uc.Execute(ctx, "link-secret", domain.SessionClient{}); result != nil || !errors.Is(err, domain.ErrLoginTokenAlreadyConsumed) {
		t.Fatalf("expected consumed-token error: %v", err)
	}
	all, err := sessions.ListSessionsByUserID(ctx, user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatal("reuse must not issue a second session")
	}
}

func TestConsumeMagicLinkRollsBackTokenOnSessionConflict(t *testing.T) {
	ctx := t.Context()
	sessions, db, _, user := setupSessionRepository(ctx, t)
	existing := newSession(t, "existing-session", user.ID())
	if err := sessions.CreateSession(ctx, existing, "existing-secret"); err != nil {
		t.Fatal(err)
	}
	tokens := identitypg.NewLoginTokenRepository(db)
	token := newLoginToken(t, "rollback-link", user.ID())
	if err := tokens.Create(ctx, token, "link-secret"); err != nil {
		t.Fatal(err)
	}
	uc, err := commands.NewConsumeMagicLink(identitypg.NewUnitOfWork(db, testEmailProtector{}, nil), cryptography.RandomTokenGenerator{}, fixedSessionID{value: existing.ID()}, clock.SystemClock{}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if result, err := uc.Execute(ctx, "link-secret", domain.SessionClient{}); result != nil || !errors.Is(err, application.ErrSessionAlreadyExists) {
		t.Fatalf("expected session conflict: %v", err)
	}
	restored, err := tokens.FindByID(ctx, token.ID())
	if err != nil {
		t.Fatal(err)
	}
	if restored.IsConsumed() {
		t.Fatal("failed session creation must roll back consumption")
	}
	all, err := sessions.ListSessionsByUserID(ctx, user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatal("rollback left an extra session")
	}
	// The same link remains usable after the failed transaction.
	retry, err := commands.NewConsumeMagicLink(identitypg.NewUnitOfWork(db, testEmailProtector{}, nil), cryptography.RandomTokenGenerator{}, id.UUIDGenerator{}, clock.SystemClock{}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := retry.Execute(ctx, "link-secret", domain.SessionClient{}); err != nil {
		t.Fatal(err)
	}
}

func TestConsumeMagicLinkConcurrentConsumptionIssuesOneSession(t *testing.T) {
	ctx := t.Context()
	sessions, db, _, user := setupSessionRepository(ctx, t)
	tokens := identitypg.NewLoginTokenRepository(db)
	token := newLoginToken(t, "concurrent-link", user.ID())
	if err := tokens.Create(ctx, token, "link-secret"); err != nil {
		t.Fatal(err)
	}
	uc, err := commands.NewConsumeMagicLink(identitypg.NewUnitOfWork(db, testEmailProtector{}, nil), cryptography.RandomTokenGenerator{}, id.UUIDGenerator{}, clock.SystemClock{}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	for range 2 {
		go func() { <-start; _, err := uc.Execute(ctx, "link-secret", domain.SessionClient{}); results <- err }()
	}
	close(start)
	first, second := <-results, <-results
	if !(first == nil && errors.Is(second, domain.ErrLoginTokenAlreadyConsumed) || second == nil && errors.Is(first, domain.ErrLoginTokenAlreadyConsumed)) {
		t.Fatalf("expected one success and one consumed-token error: %v, %v", first, second)
	}
	all, err := sessions.ListSessionsByUserID(ctx, user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatalf("expected one session, got %d", len(all))
	}
}

func TestConsumeMagicLinkRejectsExpiredAndMissingTokens(t *testing.T) {
	ctx := t.Context()
	sessions, db, _, user := setupSessionRepository(ctx, t)
	tokens := identitypg.NewLoginTokenRepository(db)
	token := newLoginToken(t, "expired-link", user.ID())
	if err := tokens.Create(ctx, token, "expired-secret"); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(ctx, "UPDATE login_tokens SET created_at = $2, expires_at = $3 WHERE id = $1", token.ID(), time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	uc, err := commands.NewConsumeMagicLink(identitypg.NewUnitOfWork(db, testEmailProtector{}, nil), cryptography.RandomTokenGenerator{}, id.UUIDGenerator{}, clock.SystemClock{}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		secret string
		want   error
	}{{"expired-secret", domain.ErrLoginTokenExpired}, {"missing", application.ErrLoginTokenNotFound}} {
		if got, err := uc.Execute(ctx, tt.secret, domain.SessionClient{}); got != nil || !errors.Is(err, tt.want) {
			t.Fatalf("want %v, got %v", tt.want, err)
		}
	}
	all, err := sessions.ListSessionsByUserID(ctx, user.ID())
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Fatal("invalid links must not issue sessions")
	}
	unchanged, err := tokens.FindByID(ctx, token.ID())
	if err != nil {
		t.Fatal(err)
	}
	if unchanged.IsConsumed() {
		t.Fatal("expired token must remain unconsumed")
	}
}
