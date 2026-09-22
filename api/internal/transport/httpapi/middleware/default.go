// Package middleware contains HTTP-wide cross-cutting policies.
package middleware

import (
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// ApplyDefault registers infrastructure middleware shared by every endpoint.
func ApplyDefault(router routing.Router) {
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
}
