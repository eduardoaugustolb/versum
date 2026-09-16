package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/ports/httprouter"
)

func registerIdentityAccessRoutes(router httprouter.Router, deps IdentityAccessDependencies) {
	router.Post("/auth/magic-link", func(w http.ResponseWriter, r *http.Request) {
		rawEmail := r.FormValue("email")
		email, err := domain.ParseEmail(rawEmail)
		if err != nil {
			slog.Error("invalid email", "error", err)
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}
		tx, err := deps.DbExecutor.Begin(r.Context())
		if err != nil {
			slog.Error("failed to begin transaction", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		unitOfWork := postgres.NewUnitOfWork(tx, deps.EmailProtector, deps.OutboxPayloadProtector)
		err = unitOfWork.WithinTransaction(r.Context(), func(repositories ports.IdentityAccessUnitOfWorkRepositories) error {
			uc, err := commands.NewRequestMagicLink(repositories.Users, repositories.LoginTokens, deps.LoginTokenGenerator, deps.LoginTokenHasher, deps.IdGenerator, deps.Clock, deps.MagicLinkPolicy.TTL)
			if err != nil {
				return err
			}
			return uc.Execute(r.Context(), email)
		})
		if err != nil {
			slog.Error("failed to request magic link", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		slog.Info("magic link requested", "email", email)
		w.WriteHeader(http.StatusOK)
	})
}
