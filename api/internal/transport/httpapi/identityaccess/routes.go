// Package identityaccess adapts identity and access use cases to HTTP endpoints.
package identityaccess

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
)

// Dependencies are the identity and access use cases required by this HTTP adapter.
type Dependencies struct{ RequestMagicLink *commands.RequestMagicLink }

// RegisterRoutes registers identity and access endpoints.
func RegisterRoutes(router routing.Router, deps Dependencies) {
	handler := handler{deps: deps}
	router.Route("/auth", func(router routing.Router) {
		router.Post("/magic-link", handler.requestMagicLink)
	})
}

type handler struct{ deps Dependencies }

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
	if h.deps.RequestMagicLink == nil {
		slog.Error("request magic link use case is not configured")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := h.deps.RequestMagicLink.Execute(request.Context(), body.Email); err != nil {
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
