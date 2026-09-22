package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxdomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

type userRepository struct {
	user      *domain.User
	createErr error
	findErr   error
}

func (r userRepository) CreateUser(context.Context, *domain.User) error { return r.createErr }
func (r userRepository) FindUserByID(context.Context, string) (*domain.User, error) {
	return r.user, nil
}
func (r userRepository) FindUserByEmail(context.Context, domain.Email) (*domain.User, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.user, nil
}

type lookupBeforeCreateUserRepository struct {
	user         *domain.User
	createCalled bool
}

func (r *lookupBeforeCreateUserRepository) CreateUser(context.Context, *domain.User) error {
	r.createCalled = true
	return nil
}

func (r *lookupBeforeCreateUserRepository) FindUserByID(context.Context, string) (*domain.User, error) {
	return r.user, nil
}

func (r *lookupBeforeCreateUserRepository) FindUserByEmail(context.Context, domain.Email) (*domain.User, error) {
	return r.user, nil
}

type loginTokenRepository struct{ token *domain.LoginToken }

func (r *loginTokenRepository) CreateLoginToken(_ context.Context, token *domain.LoginToken) error {
	r.token = token
	return nil
}
func (r *loginTokenRepository) FindLoginTokenByID(context.Context, string) (*domain.LoginToken, error) {
	return nil, application.ErrLoginTokenNotFound
}
func (r *loginTokenRepository) FindLoginTokenByTokenHash(context.Context, []byte) (*domain.LoginToken, error) {
	return nil, application.ErrLoginTokenNotFound
}
func (r *loginTokenRepository) ConsumeLoginTokenByTokenHash(context.Context, []byte, *time.Time) error {
	return nil
}

type tokenGenerator struct{}

func (tokenGenerator) GenerateLoginToken() (string, error) { return "raw-token", nil }

type tokenHasher struct{}

func (tokenHasher) Hash(token string) ([]byte, error) {
	if token != "raw-token" {
		return nil, errors.New("unexpected token")
	}
	return []byte("hashed-token"), nil
}

type idGenerator struct{}

func (idGenerator) Generate() string { return "token-1" }

type outboxRepository struct{ event *outboxdomain.Event }

func (r *outboxRepository) Publish(_ context.Context, event *outboxdomain.Event) error {
	r.event = event
	return nil
}

type unitOfWork struct {
	users  identityports.UserRepository
	tokens identityports.LoginTokenRepository
	outbox identityports.IdentityAccessUnitOfWorkRepositories
}

func (u unitOfWork) WithinTransaction(ctx context.Context, fn func(identityports.IdentityAccessUnitOfWorkRepositories) error) error {
	return fn(u.outbox)
}

func TestRequestMagicLinkUsesClockAndConfiguredTTL(t *testing.T) {
	user, err := domain.NewUser("user-1", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
	tokens := &loginTokenRepository{}
	outbox := &outboxRepository{}
	transactions := unitOfWork{outbox: identityports.IdentityAccessUnitOfWorkRepositories{Users: userRepository{user: user}, LoginTokens: tokens, Outbox: outbox}}
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		tokenGenerator{},
		tokenHasher{},
		idGenerator{},
		fixedClock{now: now},
		15*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), user.Email().String()); err != nil {
		t.Fatal(err)
	}
	if tokens.token == nil {
		t.Fatal("expected a persisted login token")
	}
	if outbox.event == nil || outbox.event.EventType() != outboxdomain.EventTypeMagicLinkRequested {
		t.Fatalf("expected magic link event, got %#v", outbox.event)
	}
	if tokens.token.ID() != "token-1" || string(tokens.token.TokenHash()) != "hashed-token" {
		t.Fatalf("unexpected token: %+v", tokens.token)
	}
	if tokens.token.ExpiresAt().Location() != time.UTC {
		t.Fatalf("expected UTC expiration, got %s", tokens.token.ExpiresAt().Location())
	}
	wantExpiresAt := now.UTC().Add(15 * time.Minute)
	if !tokens.token.ExpiresAt().Equal(wantExpiresAt) {
		t.Fatalf("expected expiration %s, got %s", wantExpiresAt, tokens.token.ExpiresAt())
	}
}

func TestNewRequestMagicLinkRejectsNonPositiveTTL(t *testing.T) {
	_, err := commands.NewRequestMagicLink(nil, nil, nil, nil, fixedClock{}, 0)
	if !errors.Is(err, application.ErrInvalidMagicLinkTTL) {
		t.Fatalf("expected invalid TTL error, got %v", err)
	}
}

func TestRequestMagicLinkUsesExistingUser(t *testing.T) {
	existingUser, err := domain.NewUser("existing-user", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	tokens := &loginTokenRepository{}
	outbox := &outboxRepository{}
	transactions := unitOfWork{outbox: identityports.IdentityAccessUnitOfWorkRepositories{
		Users:       userRepository{user: existingUser, createErr: application.ErrUserAlreadyExists},
		LoginTokens: tokens,
		Outbox:      outbox,
	}}
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		tokenGenerator{},
		tokenHasher{},
		idGenerator{},
		fixedClock{now: time.Now()},
		15*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), existingUser.Email().String()); err != nil {
		t.Fatal(err)
	}
	if tokens.token == nil || tokens.token.UserID() != existingUser.ID() {
		t.Fatalf("expected a token for existing user %q, got %#v", existingUser.ID(), tokens.token)
	}
}

func TestRequestMagicLinkDoesNotCreateUserFoundByPreviousLookupKey(t *testing.T) {
	existingUser, err := domain.NewUser("existing-user", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	users := &lookupBeforeCreateUserRepository{user: existingUser}
	tokens := &loginTokenRepository{}
	outbox := &outboxRepository{}
	transactions := unitOfWork{outbox: identityports.IdentityAccessUnitOfWorkRepositories{
		Users:       users,
		LoginTokens: tokens,
		Outbox:      outbox,
	}}
	useCase, err := commands.NewRequestMagicLink(
		transactions,
		tokenGenerator{},
		tokenHasher{},
		idGenerator{},
		fixedClock{now: time.Now()},
		15*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), existingUser.Email().String()); err != nil {
		t.Fatal(err)
	}
	if users.createCalled {
		t.Fatal("expected existing user found by a previous lookup key not to be created again")
	}
	if tokens.token == nil || tokens.token.UserID() != existingUser.ID() {
		t.Fatalf("expected a token for existing user %q, got %#v", existingUser.ID(), tokens.token)
	}
}
