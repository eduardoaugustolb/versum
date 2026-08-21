package queries

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type ListBooks struct{ reader ports.BookReader }

func NewListBooks(reader ports.BookReader) *ListBooks { return &ListBooks{reader: reader} }
func (q *ListBooks) Execute(ctx context.Context) ([]domain.Book, error) {
	return q.reader.ListBooks(ctx)
}
