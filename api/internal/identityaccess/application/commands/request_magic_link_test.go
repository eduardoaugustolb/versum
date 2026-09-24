package commands_test

import (
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxdomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

func newMagicLinkDeps(users *FakeUserRepository, tokens *FakeLoginTokenRepository, outbox *FakeOutboxRepository) identityports.IdentityAccessUnitOfWork {
	return &FakeUnitOfWork{Inner: identityports.IdentityAccessUnitOfWorkRepositories{
		Users:       users,
		LoginTokens: tokens,
		Outbox:      outbox,
	}}
}

func TestRequestMagicLinkUsesClockAndConfiguredTTL(t *testing.T) {
	user, err := domain.NewUser("user-1", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
	users := &FakeUserRepository{ByEmail: map[string]*domain.User{user.Email().String(): user}}
	tokens := &FakeLoginTokenRepository{}
	outbox := &FakeOutboxRepository{}
	transactions := newMagicLinkDeps(users, tokens, outbox)
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		FakeLoginTokenGenerator{Token: "raw-token"},
		FakeLoginTokenHasher{Sum: []byte("hashed-token")},
		&FakeIDGenerator{IDs: []string{"ignored-user-id", "token-1", "event-1"}},
		FakeClock{NowTime: now},
		policy.DefaultMagicLinkPolicy,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), user.Email()); err != nil {
		t.Fatal(err)
	}
	if tokens.Created == nil {
		t.Fatal("expected a persisted login token")
	}
	if len(outbox.Published) != 1 || outbox.Published[0].EventType() != outboxdomain.EventTypeMagicLinkRequested {
		t.Fatalf("expected magic link event, got %#v", outbox.Published)
	}
	if tokens.Created.ID() != "token-1" || string(tokens.Created.TokenHash()) != "hashed-token" {
		t.Fatalf("unexpected token: %+v", tokens.Created)
	}
	if tokens.Created.ExpiresAt().Location() != time.UTC {
		t.Fatalf("expected UTC expiration, got %s", tokens.Created.ExpiresAt().Location())
	}
	wantExpiresAt := now.UTC().Add(15 * time.Minute)
	if !tokens.Created.ExpiresAt().Equal(wantExpiresAt) {
		t.Fatalf("expected expiration %s, got %s", wantExpiresAt, tokens.Created.ExpiresAt())
	}
}

func TestNewRequestMagicLinkRejectsNonPositiveTTL(t *testing.T) {
	p := policy.DefaultMagicLinkPolicy
	p.TTL = 0
	_, err := commands.NewRequestMagicLink(nil, nil, nil, nil, FakeClock{}, p)
	if !errors.Is(err, application.ErrInvalidMagicLinkTTL) {
		t.Fatalf("expected invalid TTL error, got %v", err)
	}
}

func TestRequestMagicLinkUsesExistingUser(t *testing.T) {
	existingUser, err := domain.NewUser("existing-user", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	users := &FakeUserRepository{ByEmail: map[string]*domain.User{existingUser.Email().String(): existingUser}}
	tokens := &FakeLoginTokenRepository{}
	outbox := &FakeOutboxRepository{}
	transactions := newMagicLinkDeps(users, tokens, outbox)
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		FakeLoginTokenGenerator{},
		FakeLoginTokenHasher{},
		&FakeIDGenerator{IDs: []string{"ignored-user-id", "token-1", "event-1"}},
		FakeClock{NowTime: time.Now()},
		policy.DefaultMagicLinkPolicy,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), existingUser.Email()); err != nil {
		t.Fatal(err)
	}
	if users.CreateCalls != 0 {
		t.Fatal("expected existing user not to be created again")
	}
	if tokens.Created == nil || tokens.Created.UserID() != existingUser.ID() {
		t.Fatalf("expected a token for existing user %q, got %#v", existingUser.ID(), tokens.Created)
	}
}

func TestRequestMagicLinkDoesNotCreateUserFoundByPreviousLookupKey(t *testing.T) {
	existingUser, err := domain.NewUser("existing-user", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	users := &FakeUserRepository{ByEmail: map[string]*domain.User{existingUser.Email().String(): existingUser}}
	tokens := &FakeLoginTokenRepository{}
	outbox := &FakeOutboxRepository{}
	transactions := newMagicLinkDeps(users, tokens, outbox)
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		FakeLoginTokenGenerator{},
		FakeLoginTokenHasher{},
		&FakeIDGenerator{IDs: []string{"ignored-user-id", "token-1", "event-1"}},
		FakeClock{NowTime: time.Now()},
		policy.DefaultMagicLinkPolicy,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), existingUser.Email()); err != nil {
		t.Fatal(err)
	}
	if users.CreateCalls != 0 {
		t.Fatal("expected existing user found by a previous lookup key not to be created again")
	}
	if tokens.Created == nil || tokens.Created.UserID() != existingUser.ID() {
		t.Fatalf("expected a token for existing user %q, got %#v", existingUser.ID(), tokens.Created)
	}
}

func TestRequestMagicLinkCreatesMissingUser(t *testing.T) {
	email, err := domain.ParseEmail("novo@example.com")
	if err != nil {
		t.Fatal(err)
	}
	users := &FakeUserRepository{FindByEmailErr: application.ErrUserNotFound}
	tokens := &FakeLoginTokenRepository{}
	outbox := &FakeOutboxRepository{}
	transactions := newMagicLinkDeps(users, tokens, outbox)
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		FakeLoginTokenGenerator{},
		FakeLoginTokenHasher{},
		&FakeIDGenerator{IDs: []string{"new-user", "token-1", "event-1"}},
		FakeClock{NowTime: time.Now()},
		policy.DefaultMagicLinkPolicy,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), email); err != nil {
		t.Fatal(err)
	}
	if users.CreateCalls != 1 {
		t.Fatalf("expected one user creation, got %d", users.CreateCalls)
	}
	if tokens.Created == nil || tokens.Created.UserID() != "new-user" {
		t.Fatalf("expected a token for the new user, got %#v", tokens.Created)
	}
}
