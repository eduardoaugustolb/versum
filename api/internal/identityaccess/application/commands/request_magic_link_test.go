package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/id"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

type userRepository struct{ user *domain.User }

func (r userRepository) CreateUser(context.Context, *domain.User) error { return nil }
func (r userRepository) FindUserByID(context.Context, string) (*domain.User, error) {
	return r.user, nil
}
func (r userRepository) FindUserByEmail(context.Context, domain.Email) (*domain.User, error) {
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

func (idGenerator) Generate() id.UUID { return id.UUID("token-1") }

func TestRequestMagicLinkUsesClockAndConfiguredTTL(t *testing.T) {
	email, err := domain.ParseEmail("ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	user, err := domain.NewUser("user-1", email)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
	tokens := &loginTokenRepository{}
	useCase, err := commands.NewRequestMagicLink(
		userRepository{user: user},
		tokens,
		tokenGenerator{},
		tokenHasher{},
		idGenerator{},
		fixedClock{now: now},
		15*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := useCase.Execute(t.Context(), email); err != nil {
		t.Fatal(err)
	}
	if tokens.token == nil {
		t.Fatal("expected a persisted login token")
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
	_, err := commands.NewRequestMagicLink(nil, nil, nil, nil, nil, fixedClock{}, 0)
	if !errors.Is(err, application.ErrInvalidMagicLinkTTL) {
		t.Fatalf("expected invalid TTL error, got %v", err)
	}
}
