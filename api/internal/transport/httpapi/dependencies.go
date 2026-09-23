package httpapi

import (
	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	cataloghttp "github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/catalog"
	identityaccesshttp "github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/identityaccess"
)

type Dependencies struct {
	Health         health.CheckHealth
	Catalog        CatalogDependencies
	IdentityAccess IdentityAccessDependencies
	Cache          cache.Cache
}

// CatalogDependencies is retained at the composition boundary for callers.
type CatalogDependencies = cataloghttp.Dependencies

// IdentityAccessDependencies is retained at the composition boundary for callers.
type IdentityAccessDependencies = identityaccesshttp.Dependencies
