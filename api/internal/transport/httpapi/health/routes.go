// Package health adapts health application use cases to HTTP endpoints.
package health

import (
	"net/http"

	application "github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/response"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
)

// RegisterRoutes registers health endpoints.
func RegisterRoutes(router routing.Router, useCase application.CheckHealth) {
	router.Get("/health", handler{useCase: useCase}.check)
}

type handler struct {
	useCase application.CheckHealth
}

func (h handler) check(w http.ResponseWriter, request *http.Request) {
	status := h.useCase.Execute()
	response.WriteJSON(w, request, http.StatusOK, healthResponse{Status: status.State})
}

type healthResponse struct {
	Status string `json:"status"`
}
