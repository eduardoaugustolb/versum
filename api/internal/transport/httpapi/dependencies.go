package httpapi

import (
	"github.com/eduardoaugustolb/versum/api/internal/catalog"
	"github.com/eduardoaugustolb/versum/api/internal/health"
)

type Dependencies struct {
	Health  health.CheckHealth
	Catalog CatalogDependencies
}

type CatalogDependencies struct {
	ListBooks  *catalog.ListBooks
	GetChapter *catalog.GetChapter
}
