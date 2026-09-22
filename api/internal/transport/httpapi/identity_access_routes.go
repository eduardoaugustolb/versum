package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func registerIdentityAccessRoutes(router Router, deps IdentityAccessDependencies) {
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
		if deps.RequestMagicLink == nil {
			slog.Error("request magic link use case is not configured")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		err = deps.RequestMagicLink.Execute(r.Context(), rawEmail)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidEmail) {
				http.Error(w, "invalid email", http.StatusBadRequest)
				return
			}

			slog.Error("failed to request magic link", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
