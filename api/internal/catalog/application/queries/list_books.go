package queries

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type ListBooks struct{ repository ports.CatalogRepository }

func NewListBooks(repository ports.CatalogRepository) *ListBooks {
	return &ListBooks{repository: repository}
}
func (q *ListBooks) Execute(ctx context.Context) ([]domain.Book, error) {
	return q.repository.ListBooks(ctx)
}
