package queries

import (
	"context"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
)

type ListBooks struct{ repository ports.CatalogReader }

func NewListBooks(repository ports.CatalogReader) *ListBooks {
	return &ListBooks{repository: repository}
}
func (q *ListBooks) Execute(ctx context.Context) ([]ports.BookDTO, error) {
	return q.repository.ListBooks(ctx)
}
