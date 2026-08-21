package httpapi

import (
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/health"
)

type Dependencies struct {
	Health  health.CheckHealth
	Catalog CatalogDependencies
}

type CatalogDependencies struct {
	ListBooks  *queries.ListBooks
	GetChapter *queries.GetChapter
}
