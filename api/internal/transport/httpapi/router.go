package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eduardoaugustolb/versum/api/internal/health"
)

type healthResponse struct {
	Status string `json:"status"`
}

func NewRouter(useCase health.CheckHealth) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := useCase.Execute()

		body, err := json.Marshal(healthResponse{Status: status.State})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	return router
}
