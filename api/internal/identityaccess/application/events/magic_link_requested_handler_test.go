package events

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/email"
	identity "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outbox "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type tokenReader struct {
	identity.LoginTokenRepository
	token  *domain.LoginToken
	err    error
	secret string
}

func (r *tokenReader) FindByToken(_ context.Context, secret string) (*domain.LoginToken, error) {
	r.secret = secret
	return r.token, r.err
}

type userReader struct {
	identity.UserRepository
	user *domain.User
	err  error
}

func (r *userReader) FindUserByID(_ context.Context, _ string) (*domain.User, error) {
	return r.user, r.err
}

type senderSpy struct {
	requests []email.SendRequest
	err      error
}

func (s *senderSpy) Send(_ context.Context, req email.SendRequest) error {
	s.requests = append(s.requests, req)
	return s.err
}

type deliveryClock struct{ now time.Time }

func (c *deliveryClock) Now() time.Time { return c.now }

func TestMagicLinkHandlerDeliveryAndRetry(t *testing.T) {
	now := time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
	token, _ := domain.NewLoginToken("token-1", "user-1", now.Add(15*time.Minute), now)
	user, _ := domain.NewUser("user-1", "ana@example.com")
	tokens := &tokenReader{token: token}
	sendErr := errors.New("provider unavailable")
	sender := &senderSpy{err: sendErr}
	clock := &deliveryClock{now: now}
	handler, err := NewMagicLinkRequestedHandler(sender, tokens, &userReader{user: user}, clock, "https://app.example.com/access/confirm")
	if err != nil {
		t.Fatal(err)
	}
	secret := `abc+/?&"<script>`
	payload, _ := json.Marshal(MagicLinkRequestedPayload{LoginTokenID: token.ID(), Token: secret})
	event, _ := outbox.NewEvent("event-1", MagicLinkRequested, user.ID(), payload)
	if err := handler.Handle(t.Context(), event); !errors.Is(err, sendErr) {
		t.Fatalf("send failure must reach processor: %v", err)
	}
	sender.err = nil
	clock.now = now.Add(time.Minute)
	if err := handler.Handle(t.Context(), event); err != nil {
		t.Fatal(err)
	}
	if len(sender.requests) != 2 || !reflect.DeepEqual(sender.requests[0], sender.requests[1]) {
		t.Fatal("retry changed request or idempotency key")
	}
	req := sender.requests[0]
	if req.IdempotencyKey != "magic-link/event-1" || req.To != user.Email().String() || tokens.secret != secret {
		t.Fatal("unexpected delivery identity")
	}
	link := "https://app.example.com/access/confirm?" + url.Values{"token": []string{secret}}.Encode()
	if !strings.Contains(req.Text, link) || !strings.Contains(req.HTML, link) || strings.Contains(req.HTML, "<script>") || strings.Contains(req.HTML, "#ZgotmplZ") {
		t.Fatal("incorrect or unsafe link rendering")
	}
	if !strings.Contains(req.HTML, "Acessar minha conta") || !strings.Contains(req.Text, "22/10/2026 às 12:15:00 UTC") {
		t.Fatal("missing access instructions or absolute expiry")
	}
	other, _ := outbox.NewEvent("event-2", MagicLinkRequested, user.ID(), payload)
	if err := handler.Handle(t.Context(), other); err != nil {
		t.Fatal(err)
	}
	if sender.requests[2].IdempotencyKey == req.IdempotencyKey {
		t.Fatal("different events share an idempotency key")
	}
}

func TestMagicLinkHandlerRejectsOrSkipsUndeliverableEvents(t *testing.T) {
	now := time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
	lookupErr := errors.New("database unavailable")
	for _, tt := range []struct {
		name, owner, payload    string
		consumed, expired       bool
		tokenErr, userErr, want error
	}{
		{name: "empty payload fields", payload: `{}`, want: ErrInvalidMagicLinkEvent},
		{name: "malformed json", payload: `{`, want: ErrInvalidMagicLinkEvent},
		{name: "wrong token id", payload: `{"login_token_id":"other","token":"secret"}`, want: ErrInvalidMagicLinkEvent},
		{name: "wrong owner", owner: "other", want: ErrInvalidMagicLinkEvent},
		{name: "missing owner", owner: "-", want: ErrInvalidMagicLinkEvent},
		{name: "consumed", consumed: true},
		{name: "expired at boundary", expired: true},
		{name: "deleted token", tokenErr: identity.ErrLoginTokenNotFound},
		{name: "deleted user", userErr: identity.ErrUserNotFound},
		{name: "token lookup unavailable", tokenErr: lookupErr, want: lookupErr},
		{name: "user lookup unavailable", userErr: lookupErr, want: lookupErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			expiresAt := now.Add(time.Minute)
			if tt.expired {
				expiresAt = now
			}
			var consumedAt *time.Time
			if tt.consumed {
				at := now.Add(-time.Second)
				consumedAt = &at
			}
			token, err := domain.RehydrateLoginToken("token-1", "user-1", expiresAt, consumedAt)
			if err != nil {
				t.Fatal(err)
			}
			user, _ := domain.NewUser("user-1", "ana@example.com")
			sender := &senderSpy{}
			handler, err := NewMagicLinkRequestedHandler(sender, &tokenReader{token: token, err: tt.tokenErr}, &userReader{user: user, err: tt.userErr}, &deliveryClock{now: now}, "https://app.example.com/access/confirm")
			if err != nil {
				t.Fatal(err)
			}
			owner := tt.owner
			if owner == "" {
				owner = "user-1"
			}
			if owner == "-" {
				owner = ""
			}
			payload := tt.payload
			if payload == "" {
				payload = `{"login_token_id":"token-1","token":"secret"}`
			}
			event, _ := outbox.NewEvent("event-1", MagicLinkRequested, owner, []byte(payload))
			err = handler.Handle(t.Context(), event)
			if !errors.Is(err, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, err)
			}
			if len(sender.requests) != 0 {
				t.Fatal("undeliverable event sent an email")
			}
		})
	}
}

func TestMagicLinkHandlerRequiresTrustedURL(t *testing.T) {
	for _, raw := range []string{"", "javascript:alert(1)", "http://app.example.com", "https://user:password@app.example.com", "https://app.example.com?token=old", "https://app.example.com#fragment"} {
		_, err := NewMagicLinkRequestedHandler(&senderSpy{}, &tokenReader{}, &userReader{}, &deliveryClock{}, raw)
		if !errors.Is(err, ErrInvalidMagicLinkHandler) {
			t.Fatalf("accepted invalid URL %q", raw)
		}
	}
}
