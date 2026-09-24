package httpapi_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	redisadapter "github.com/eduardoaugustolb/versum/api/internal/cache/redis"
	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	identityapplication "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identitycommands "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	identitydomain "github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxdomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
	redisclient "github.com/redis/go-redis/v9"
)

type identityUserRepository struct {
	user    *identitydomain.User
	created *identitydomain.User
}

func (r *identityUserRepository) CreateUser(_ context.Context, user *identitydomain.User) error {
	r.created = user
	return nil
}
func (r *identityUserRepository) FindUserByID(context.Context, string) (*identitydomain.User, error) {
	return r.user, nil
}
func (r *identityUserRepository) FindUserByEmail(context.Context, identitydomain.Email) (*identitydomain.User, error) {
	if r.user == nil {
		return nil, identityapplication.ErrUserNotFound
	}
	return r.user, nil
}

type identityLoginTokenRepository struct{ token *identitydomain.LoginToken }

func (r *identityLoginTokenRepository) CreateLoginToken(_ context.Context, token *identitydomain.LoginToken) error {
	r.token = token
	return nil
}
func (r *identityLoginTokenRepository) FindLoginTokenByID(context.Context, string) (*identitydomain.LoginToken, error) {
	return nil, identityapplication.ErrLoginTokenNotFound
}
func (r *identityLoginTokenRepository) FindLoginTokenByTokenHash(context.Context, []byte) (*identitydomain.LoginToken, error) {
	return nil, identityapplication.ErrLoginTokenNotFound
}
func (r *identityLoginTokenRepository) ConsumeLoginTokenByTokenHash(context.Context, []byte, *time.Time) error {
	return nil
}

type identitySessionRepository struct{}

func (identitySessionRepository) CreateSession(context.Context, *identitydomain.Session) error {
	return nil
}
func (identitySessionRepository) FindSessionByID(context.Context, string) (*identitydomain.Session, error) {
	return nil, nil
}
func (identitySessionRepository) FindSessionBySecretHash(context.Context, []byte) (*identitydomain.Session, error) {
	return nil, nil
}
func (identitySessionRepository) ListSessionsByUserID(context.Context, string) ([]identitydomain.Session, error) {
	return nil, nil
}
func (identitySessionRepository) RevokeSession(context.Context, string, *time.Time) error { return nil }
func (identitySessionRepository) RevokeAllSessions(context.Context, string, *time.Time) error {
	return nil
}

type identityOutboxRepository struct {
	event *outboxdomain.Event
}

func (r *identityOutboxRepository) Publish(_ context.Context, event *outboxdomain.Event) error {
	r.event = event
	return nil
}

type identityUnitOfWork struct {
	repositories identityports.IdentityAccessUnitOfWorkRepositories
	called       bool
}

func (u *identityUnitOfWork) WithinTransaction(ctx context.Context, fn func(identityports.IdentityAccessUnitOfWorkRepositories) error) error {
	u.called = true
	return fn(u.repositories)
}

type identityTokenGenerator struct{}

func (identityTokenGenerator) GenerateLoginToken() (string, error) { return "raw-token", nil }

type identityTokenHasher struct{}

func (identityTokenHasher) Hash(value string) ([]byte, error) {
	if value != "raw-token" {
		return nil, errors.New("unexpected token")
	}
	return []byte("hashed-token"), nil
}

type identityIDGenerator struct{ next int }

func (g *identityIDGenerator) Generate() string {
	g.next++
	return "id-" + string(rune('0'+g.next))
}

type identityClock struct{}

func (identityClock) Now() time.Time { return time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC) }

var _ clock.Clock = identityClock{}

func newIdentityAccessHandler(t *testing.T) (*httptest.ResponseRecorder, http.Handler) {
	t.Helper()
	redisServer := miniredis.RunT(t)
	redisClient := redisclient.NewClient(&redisclient.Options{Addr: redisServer.Addr()})
	t.Cleanup(func() { _ = redisClient.Close() })

	users := &identityUserRepository{}
	tokens := &identityLoginTokenRepository{}
	outbox := &identityOutboxRepository{}
	uow := &identityUnitOfWork{repositories: identityports.IdentityAccessUnitOfWorkRepositories{
		Users: users, LoginTokens: tokens, Sessions: identitySessionRepository{}, Outbox: outbox,
	}}
	useCase, err := identitycommands.NewRequestMagicLink(uow, identityTokenGenerator{}, identityTokenHasher{}, &identityIDGenerator{}, identityClock{}, policy.DefaultMagicLinkPolicy)
	if err != nil {
		t.Fatal(err)
	}
	handler := httpapi.NewHandler(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Cache:  redisadapter.NewRedisCache(redisClient),
		IdentityAccess: httpapi.IdentityAccessDependencies{
			RequestMagicLink: useCase,
			Cache:            redisadapter.NewRedisCache(redisClient),
		},
	})
	return httptest.NewRecorder(), handler
}

func TestRequestMagicLinkEndpointRateLimitsRequestsFromSameAddress(t *testing.T) {
	_, handler := newIdentityAccessHandler(t)

	for requestNumber := 1; requestNumber <= 10; requestNumber++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"ana@example.com"}`))
		request.Header.Set("Content-Type", "application/json")
		request.RemoteAddr = "198.51.100.1:54321"
		handler.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d: %s", requestNumber, recorder.Code, recorder.Body.String())
		}
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"ana@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = "198.51.100.1:54321"
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestRequestMagicLinkEndpointRateLimitsSameEmailAcrossAddresses(t *testing.T) {
	_, handler := newIdentityAccessHandler(t)

	for requestNumber := 1; requestNumber <= 10; requestNumber++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":" ANA@EXAMPLE.COM "}`))
		request.Header.Set("Content-Type", "application/json")
		request.RemoteAddr = fmt.Sprintf("198.51.100.%d:54321", requestNumber)
		handler.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d: %s", requestNumber, recorder.Code, recorder.Body.String())
		}
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"ana@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = "198.51.100.11:54321"
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestRequestMagicLinkEndpoint(t *testing.T) {
	recorder, handler := newIdentityAccessHandler(t)
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"ana@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestRequestMagicLinkEndpointRejectsInvalidJSON(t *testing.T) {
	recorder, handler := newIdentityAccessHandler(t)
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestRequestMagicLinkEndpointRejectsInvalidEmail(t *testing.T) {
	recorder, handler := newIdentityAccessHandler(t)
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"invalid"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestRequestMagicLinkEndpointRejectsInvalidContentType(t *testing.T) {
	recorder, handler := newIdentityAccessHandler(t)
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"email@email.com"}`))
	request.Header.Set("Content-Type", "application/xml")
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}
