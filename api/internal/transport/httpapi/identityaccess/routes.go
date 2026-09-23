// Package identityaccess adapts identity and access use cases to HTTP endpoints.
package identityaccess

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/middleware"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
)

// Dependencies are the identity and access use cases required by this HTTP adapter.
type Dependencies struct {
	RequestMagicLink *commands.RequestMagicLink
	Cache            cache.Cache
}

// RegisterRoutes registers identity and access endpoints.
func RegisterRoutes(router routing.Router, deps Dependencies) {
	if deps.Cache == nil {
		slog.Warn("identityaccess cache not configured, magic-link requests will fail")
	}
	handler := handler{
		deps:           deps,
		emailRateLimit: middleware.NewRateLimit(func(_ *http.Request) string { return "identityaccess:magic-link:email" }, time.Minute, 10, deps.Cache, false),
	}
	ipRateLimit := middleware.NewRateLimit(func(_ *http.Request) string { return "identityaccess:magic-link:ip" }, time.Minute, 10, deps.Cache, false)
	router.Route("/auth", func(router routing.Router) {
		router.Use(ipRateLimit.Execute)
		router.Post("/magic-link", handler.requestMagicLink)
	})
}

type handler struct {
	deps           Dependencies
	emailRateLimit *middleware.RateLimit
}

func (h handler) requestMagicLink(w http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}
	if request.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	var body struct {
		Email string `json:"email"`
	}
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	email, err := domain.ParseEmail(body.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	allowed, err := h.emailRateLimit.Allow(request.Context(), rateLimitEmailKey(email))
	if err != nil {
		slog.ErrorContext(request.Context(), "failed to allow magic link email rate limit", "error", err)
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	if !allowed {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	if h.deps.RequestMagicLink == nil {
		slog.Error("request magic link use case is not configured")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := h.deps.RequestMagicLink.Execute(request.Context(), email.String()); err != nil {
		if errors.Is(err, domain.ErrInvalidEmail) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}
		slog.ErrorContext(request.Context(), "failed to request magic link", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func rateLimitEmailKey(email domain.Email) string {
	digest := sha256.Sum256([]byte(email.String()))
	return "identityaccess:magic-link:email:" + hex.EncodeToString(digest[:])
}
