package httpapi

import (
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
)

type Dependencies struct {
	Health         health.CheckHealth
	Catalog        CatalogDependencies
	IdentityAccess IdentityAccessDependencies
}

type CatalogDependencies struct {
	ListBooks  *queries.ListBooks
	GetChapter *queries.GetChapter
}

type IdentityAccessDependencies struct {
	RequestMagicLink *commands.RequestMagicLink
}
