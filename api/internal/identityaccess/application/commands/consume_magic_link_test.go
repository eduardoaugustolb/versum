package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type consumeTokens struct {
	FakeLoginTokenRepository
	saves int
}

func (r *consumeTokens) Save(ctx context.Context, token *domain.LoginToken) error {
	r.saves++
	return r.FakeLoginTokenRepository.Save(ctx, token)
}

func TestConsumeMagicLinkCreatesSessionAndConsumesToken(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.FixedZone("BRT", -3*3600))
	token, err := domain.NewLoginToken("link", "user", now.Add(time.Minute), now)
	if err != nil {
		t.Fatal(err)
	}
	tokens := &consumeTokens{FakeLoginTokenRepository: FakeLoginTokenRepository{ByToken: map[string]*domain.LoginToken{"link-secret": token}}}
	sessions := &FakeSessionRepository{}
	uow := &rotationUow{repos: application.IdentityAccessUnitOfWorkRepositories{LoginTokens: tokens, Sessions: sessions}}
	uc, err := commands.NewConsumeMagicLink(uow, sessionSecrets{secret: "session-secret"}, &FakeIDGenerator{IDs: []string{"session", "family"}}, FakeClock{NowTime: now}, 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	got, err := uc.Execute(t.Context(), "link-secret", "192.0.2.1", "browser")
	if err != nil {
		t.Fatal(err)
	}
	if !uow.committed || tokens.saves != 1 || sessions.CreateCalls != 1 {
		t.Fatal("must consume once and issue a session in one transaction")
	}
	if got.Secret != "session-secret" || got.Session.ID() != "session" || got.Session.FamilyID() != "family" || got.Session.UserID() != "user" || got.Session.IPAddress() != "192.0.2.1" || got.Session.UserAgent() != "browser" {
		t.Fatal("incorrect issued session")
	}
	if !got.Session.ExpiresAt().Equal(now.Add(24*time.Hour)) || got.Session.ExpiresAt().Location() != time.UTC {
		t.Fatal("incorrect session expiration")
	}
	if !got.Session.IsValidAt(now) || sessions.CreatedSecrets[0] != got.Secret {
		t.Fatal("session credential not persisted")
	}
	if at, ok := token.ConsumedAt(); !ok || !at.Equal(now) {
		t.Fatal("token not consumed")
	}
	if _, err := uc.Execute(t.Context(), "link-secret", "", ""); !errors.Is(err, domain.ErrLoginTokenAlreadyConsumed) {
		t.Fatal(err)
	}
	if sessions.CreateCalls != 1 {
		t.Fatal("duplicate consumption issued a session")
	}
}

func TestConsumeMagicLinkRejectsInvalidTTL(t *testing.T) {
	for _, ttl := range []time.Duration{0, -time.Second} {
		uc, err := commands.NewConsumeMagicLink(nil, nil, nil, nil, ttl)
		if uc != nil || !errors.Is(err, application.ErrInvalidSessionTTL) {
			t.Fatal(err)
		}
	}
}

func TestConsumeMagicLinkFailures(t *testing.T) {
	now := time.Now().UTC()
	failure := errors.New("dependency failure")
	for _, name := range []string{"empty token", "not found", "lookup", "consumed", "expired", "expiry boundary", "generation", "empty secret", "same secret", "invalid session id", "invalid family", "invalid ip", "invalid agent", "save", "create", "begin", "commit"} {
		t.Run(name, func(t *testing.T) {
			token, _ := domain.NewLoginToken("link", "user", now.Add(time.Minute), now)
			tokens := &consumeTokens{FakeLoginTokenRepository: FakeLoginTokenRepository{ByToken: map[string]*domain.LoginToken{"link-secret": token}}}
			sessions := &FakeSessionRepository{}
			uow := &rotationUow{repos: application.IdentityAccessUnitOfWorkRepositories{LoginTokens: tokens, Sessions: sessions}}
			secrets := sessionSecrets{secret: "session-secret"}
			ids := &FakeIDGenerator{IDs: []string{"session", "family"}}
			at := now
			input := "link-secret"
			ipAddress, userAgent := "", ""
			want := failure
			switch name {
			case "empty token":
				input = ""
				want = application.ErrInvalidLoginTokenSecret
			case "not found":
				tokens.ByToken = nil
				want = application.ErrLoginTokenNotFound
			case "lookup":
				tokens.FindByTokenErr = failure
			case "consumed":
				token.Consume(now)
				want = domain.ErrLoginTokenAlreadyConsumed
			case "expired":
				at = now.Add(time.Hour)
				want = domain.ErrLoginTokenExpired
			case "expiry boundary":
				at = token.ExpiresAt()
				want = domain.ErrLoginTokenExpired
			case "generation":
				secrets.err = failure
			case "empty secret":
				secrets.secret = ""
				want = application.ErrInvalidSessionSecret
			case "same secret":
				secrets.secret = input
				want = application.ErrInvalidSessionSecret
			case "invalid session id":
				ids.IDs = []string{"", "family"}
				want = domain.ErrInvalidSessionID
			case "invalid family":
				ids.IDs = []string{"session", ""}
				want = domain.ErrInvalidSessionFamilyID
			case "invalid ip":
				ipAddress = "invalid"
				want = domain.ErrInvalidSessionIPAddress
			case "invalid agent":
				userAgent = "browser\nheader"
				want = domain.ErrInvalidSessionUserAgent
			case "save":
				tokens.SaveErr = failure
			case "create":
				sessions.CreateErr = failure
			case "begin":
				uow.beginErr = failure
			case "commit":
				uow.commitErr = failure
			}
			uc, err := commands.NewConsumeMagicLink(uow, secrets, ids, FakeClock{NowTime: at}, time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			got, err := uc.Execute(t.Context(), input, ipAddress, userAgent)
			if got != nil || !errors.Is(err, want) {
				t.Fatalf("want %v and no credential, got %+v, %v", want, got, err)
			}
			if uow.committed {
				t.Fatal("failed operation must not commit")
			}
			if name != "create" && name != "commit" && sessions.CreateCalls != 0 {
				t.Fatal("must not issue a session after failure")
			}
		})
	}
}
