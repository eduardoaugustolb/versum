package application

import "context"

type CatalogUnitOfWork interface {
	WithinCatalogTransaction(context.Context, func(CatalogRepository, CatalogVersionRepository) error) error
}
