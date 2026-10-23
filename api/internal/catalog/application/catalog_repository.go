package application

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

// CatalogRepository persists catalog aggregates. Read projections use CatalogReader.
type CatalogRepository interface {
	ReplaceBook(ctx context.Context, book domain.Book, verses []domain.Verse) error
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(context.Context, CatalogRepository) error) error
}
