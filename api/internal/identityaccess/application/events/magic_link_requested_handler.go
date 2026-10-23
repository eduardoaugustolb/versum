package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/email"
	identityaccessApplication "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type MagicLinkRequestedHandler struct {
	emailSender    email.Sender
	loginTokenRepo identityaccessApplication.LoginTokenRepository
	userRepo       identityaccessApplication.UserRepository
	clock          clock.Clock
	accessURL      url.URL
}

var _ application.EventHandler = (*MagicLinkRequestedHandler)(nil)

var (
	ErrInvalidMagicLinkHandler = errors.New("invalid magic link handler configuration")
	ErrInvalidMagicLinkEvent   = errors.New("invalid magic link event")
)

// accessURL is the trusted web confirmation page, never a URL supplied by a request.
// Opening that page must not consume the token automatically (email scanners follow links).
func NewMagicLinkRequestedHandler(
	emailSender email.Sender,
	loginTokenRepo identityaccessApplication.LoginTokenRepository,
	userRepo identityaccessApplication.UserRepository,
	clock clock.Clock,
	accessURL string,
) (*MagicLinkRequestedHandler, error) {
	u, err := url.Parse(accessURL)
	if err != nil || u == nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.ForceQuery || emailSender == nil || loginTokenRepo == nil || userRepo == nil || clock == nil {
		return nil, ErrInvalidMagicLinkHandler
	}
	return &MagicLinkRequestedHandler{emailSender: emailSender, loginTokenRepo: loginTokenRepo, userRepo: userRepo, clock: clock, accessURL: *u}, nil
}

func (h *MagicLinkRequestedHandler) Handle(ctx context.Context, event *domain.Event) error {
	if event == nil || event.EventType() != MagicLinkRequested || event.UserID() == "" {
		return ErrInvalidMagicLinkEvent
	}
	var payload MagicLinkRequestedPayload
	if err := json.Unmarshal(event.Payload(), &payload); err != nil {
		return ErrInvalidMagicLinkEvent
	}
	if strings.TrimSpace(payload.LoginTokenID) == "" || strings.TrimSpace(payload.Token) == "" {
		return ErrInvalidMagicLinkEvent
	}

	// Lookup by secret also proves that the payload token matches the stored hash.
	loginToken, err := h.loginTokenRepo.FindByToken(ctx, payload.Token)
	if errors.Is(err, identityaccessApplication.ErrLoginTokenNotFound) {
		return nil // Deleted tokens must not be delivered or retried.
	}
	if err != nil {
		return fmt.Errorf("fetching login token: %w", err)
	}
	if loginToken == nil || loginToken.ID() != payload.LoginTokenID || loginToken.UserID() != event.UserID() {
		return ErrInvalidMagicLinkEvent
	}
	if loginToken.IsConsumed() || loginToken.IsExpired(h.clock.Now()) {
		return nil
	}

	user, err := h.userRepo.FindUserByID(ctx, loginToken.UserID())

	if errors.Is(err, identityaccessApplication.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("fetching user: %w", err)
	}
	if user == nil || user.ID() != event.UserID() {
		return ErrInvalidMagicLinkEvent
	}

	link := h.accessURL
	query := url.Values{"token": []string{payload.Token}}
	link.RawQuery = query.Encode()
	html, text, err := renderMagicLinkEmail(link.String(), loginToken.ExpiresAt())
	if err != nil {
		return fmt.Errorf("rendering magic link email: %w", err)
	}
	// Check again after dependencies/rendering: do not start an expired delivery.
	if loginToken.IsExpired(h.clock.Now()) {
		return nil
	}
	return h.emailSender.Send(ctx, email.SendRequest{
		To:             user.Email().String(),
		Subject:        "Seu link de acesso ao Versum",
		HTML:           html,
		Text:           text,
		IdempotencyKey: "magic-link/" + event.ID(),
	})
}
