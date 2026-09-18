package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/httprouter"
)

func registerIdentityAccessRoutes(router httprouter.Router, deps IdentityAccessDependencies) {
	router.Post("/auth/magic-link", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		if r.Body == nil {
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		}

		bodyDecoder := json.NewDecoder(r.Body)
		bodyDecoder.DisallowUnknownFields()
		var body struct {
			Email string `json:"email"`
		}

		err := bodyDecoder.Decode(&body)

		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		rawEmail := body.Email
		email, err := domain.ParseEmail(rawEmail)
		if err != nil {
			slog.Error("invalid email", "error", err)
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}
		if deps.RequestMagicLink == nil {
			slog.Error("request magic link use case is not configured")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		err = deps.RequestMagicLink.Execute(r.Context(), email)
		if err != nil {
			slog.Error("failed to request magic link", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		slog.Info("magic link requested", "email", email)
		w.WriteHeader(http.StatusOK)
	})
}
