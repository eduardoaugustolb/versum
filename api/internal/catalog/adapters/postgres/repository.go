package postgres

import (
	"context"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type CatalogDAO interface {
	ReplaceBook(context.Context, ports.BookDTO, []ports.VerseDTO) error
}
type Repository struct{ dao CatalogDAO }

func NewRepository(dao CatalogDAO) *Repository { return &Repository{dao: dao} }

var _ ports.CatalogRepository = (*Repository)(nil)

func (r *Repository) ReplaceBook(ctx context.Context, book domain.Book, verses []domain.Verse) error {
	record := ports.BookDTO{ID: book.ID(), Order: book.Order(), Name: book.Name(), Testament: string(book.Testament()), ChapterCount: book.ChapterCount()}
	records := make([]ports.VerseDTO, len(verses))
	for i, v := range verses {
		records[i] = ports.VerseDTO{BookID: v.BookID(), Chapter: v.Chapter(), Number: v.Number(), Text: v.Text(), Part: v.Part()}
	}
	return r.dao.ReplaceBook(ctx, record, records)
}
