package httpapi

import (
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/ports/httprouter"
)

type healthResponse struct {
	Status string `json:"status"`
}

func registerHealthRoutes(router httprouter.Router, useCase health.CheckHealth) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := useCase.Execute()
		writeJSON(w, http.StatusOK, healthResponse{Status: status.State})
	})
}
