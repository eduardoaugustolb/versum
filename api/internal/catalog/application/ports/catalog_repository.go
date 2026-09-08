package ports

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

// CatalogRepository persists and retrieves catalog aggregates.
type CatalogRepository interface {
	ReplaceBook(ctx context.Context, book domain.Book, verses []domain.Verse) error
	ListBooks(ctx context.Context) ([]domain.Book, error)
	FindChapter(ctx context.Context, bookID string, number int) (domain.Chapter, error)
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(context.Context, CatalogRepository) error) error
}
