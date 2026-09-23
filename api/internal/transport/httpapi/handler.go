package httpapi

import (
	"net/http"

	cataloghttp "github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/catalog"
	healthhttp "github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/health"
	identityaccesshttp "github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/identityaccess"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/middleware"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing/chiadapter"
)

// NewHandler composes the inbound HTTP adapter and returns the application
// handler consumed by net/http.
func NewHandler(deps Dependencies) http.Handler {
	router := chiadapter.New()

	middleware.ApplyDefault(router, deps.Cache)

	healthhttp.RegisterRoutes(router, deps.Health)
	cataloghttp.RegisterRoutes(router, deps.Catalog)
	if deps.IdentityAccess.Cache == nil {
		deps.IdentityAccess.Cache = deps.Cache
	}
	identityaccesshttp.RegisterRoutes(router, deps.IdentityAccess)

	return router
}

// NewRouter is kept for backward compatibility. Prefer NewHandler.
func NewRouter(deps Dependencies) http.Handler {
	return NewHandler(deps)
}
