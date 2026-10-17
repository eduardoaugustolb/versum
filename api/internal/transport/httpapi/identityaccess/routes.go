// Package identityaccess adapts identity and access use cases to HTTP endpoints.
package identityaccess

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/middleware"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
)

// Dependencies are the identity and access use cases required by this HTTP adapter.
type MagicLinkRequester interface {
	Execute(context.Context, domain.Email) error
}

type Dependencies struct {
	RequestMagicLink MagicLinkRequester
	Cache            cache.Cache
}

// RegisterRoutes registers identity and access endpoints.
func RegisterRoutes(router routing.Router, deps Dependencies) {
	if deps.Cache == nil {
		slog.Warn("identityaccess cache not configured, magic-link requests will fail")
	}
	handler := handler{
		deps: deps,
	}
	ipRateLimit := middleware.NewRateLimit(func(_ *http.Request) string { return "identityaccess:magic-link:ip" }, time.Minute, 10, deps.Cache, false)
	router.Route("/auth", func(router routing.Router) {
		router.Use(ipRateLimit.Execute)
		router.Post("/magic-link", handler.requestMagicLink)
	})
}

type handler struct {
	deps Dependencies
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
		// neutral response
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if h.deps.RequestMagicLink == nil {
		slog.Error("request magic link use case is not configured")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := h.deps.RequestMagicLink.Execute(request.Context(), email); err != nil {
		if errors.Is(err, application.ErrEmailRateLimitExceeded) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		if errors.Is(err, application.ErrEmailRateLimitUnavailable) {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		slog.ErrorContext(request.Context(), "failed to request magic link", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
