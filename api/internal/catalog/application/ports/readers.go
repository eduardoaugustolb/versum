package ports

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type BookReader interface {
	ListBooks(ctx context.Context) ([]domain.Book, error)
}

type ChapterReader interface {
	FindChapter(ctx context.Context, bookID string, number int) (domain.Chapter, error)
}

type CatalogWriter interface {
	ReplaceBook(ctx context.Context, book domain.Book, verses []domain.Verse) error
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(context.Context, CatalogWriter) error) error
}
