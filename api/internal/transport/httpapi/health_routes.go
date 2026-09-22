package httpapi

import (
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/health"
)

type healthResponse struct {
	Status string `json:"status"`
}

func registerHealthRoutes(router Router, useCase health.CheckHealth) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := useCase.Execute()
		writeJSON(w, r, http.StatusOK, healthResponse{Status: status.State})
	})
}
