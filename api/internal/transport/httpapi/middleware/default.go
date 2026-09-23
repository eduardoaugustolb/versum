// Package middleware contains HTTP-wide cross-cutting policies.
package middleware

import (
	"net/http"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// ApplyDefault registers infrastructure middleware shared by every endpoint.
func ApplyDefault(router routing.Router, cacheStore cache.Cache) {
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
	router.Use(NewRateLimit(func(_ *http.Request) string { return "global" }, time.Minute, 50, cacheStore, true).Execute)
}
