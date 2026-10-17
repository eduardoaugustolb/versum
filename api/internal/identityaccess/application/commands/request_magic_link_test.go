package commands_test

import (
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/events"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
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
		&FakeIDGenerator{IDs: []string{"ignored-user-id", "token-1", "event-1"}},
		FakeClock{NowTime: now},
		policy.DefaultMagicLinkPolicy,
		&FakeEmailRequestLimiter{Allowed: true},
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
	if len(outbox.Published) != 1 || outbox.Published[0].EventType() != events.MagicLinkRequested {
		t.Fatalf("expected magic link event, got %#v", outbox.Published)
	}
	if tokens.Created.ID() != "token-1" || tokens.CreatedSecret != "raw-token" {
		t.Fatalf("unexpected token: %+v", tokens.Created)
	}
	if outbox.Published[0].UserID() != user.ID() {
		t.Fatalf("expected event owner %q, got %q", user.ID(), outbox.Published[0].UserID())
	}
	if string(outbox.Published[0].Payload()) != `{"login_token_id":"token-1"}` {
		t.Fatalf("unexpected magic link payload: %s", outbox.Published[0].Payload())
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
	_, err := commands.NewRequestMagicLink(nil, nil, nil, FakeClock{}, p, nil)
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
		&FakeIDGenerator{IDs: []string{"ignored-user-id", "token-1", "event-1"}},
		FakeClock{NowTime: time.Now()},
		policy.DefaultMagicLinkPolicy,
		&FakeEmailRequestLimiter{Allowed: true},
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
		&FakeIDGenerator{IDs: []string{"ignored-user-id", "token-1", "event-1"}},
		FakeClock{NowTime: time.Now()},
		policy.DefaultMagicLinkPolicy,
		&FakeEmailRequestLimiter{Allowed: true},
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
		&FakeIDGenerator{IDs: []string{"new-user", "token-1", "event-1"}},
		FakeClock{NowTime: time.Now()},
		policy.DefaultMagicLinkPolicy,
		&FakeEmailRequestLimiter{Allowed: true},
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

func TestMagicLinkQuotaAppliesWithoutHTTP(t *testing.T) {
	email, _ := domain.ParseEmail("ana@example.com")
	for _, test := range []struct {
		name       string
		allowed    bool
		limiterErr error
		want       error
	}{
		{name: "quota exhausted", want: application.ErrEmailRateLimitExceeded},
		{name: "quota unavailable", limiterErr: errors.New("cache unavailable"), want: application.ErrEmailRateLimitUnavailable},
	} {
		t.Run(test.name, func(t *testing.T) {
			limiter := &FakeEmailRequestLimiter{Allowed: test.allowed, Err: test.limiterErr}
			uow := &FakeUnitOfWork{}
			uc, err := commands.NewRequestMagicLink(uow, FakeLoginTokenGenerator{}, &FakeIDGenerator{}, FakeClock{}, policy.DefaultMagicLinkPolicy, limiter)
			if err != nil {
				t.Fatal(err)
			}
			if err := uc.Execute(t.Context(), email); !errors.Is(err, test.want) {
				t.Fatalf("want %v, got %v", test.want, err)
			}
			if uow.Calls != 0 {
				t.Fatal("quota failure started a transaction")
			}
		})
	}
	if _, err := commands.NewRequestMagicLink(nil, nil, nil, FakeClock{}, policy.DefaultMagicLinkPolicy, nil); !errors.Is(err, application.ErrEmailRateLimitUnavailable) {
		t.Fatal(err)
	}
}
