package queries

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type GetChapter struct{ repository ports.CatalogRepository }

func NewGetChapter(repository ports.CatalogRepository) *GetChapter {
	return &GetChapter{repository: repository}
}
func (q *GetChapter) Execute(ctx context.Context, bookID string, number int) (domain.Chapter, error) {
	return q.repository.FindChapter(ctx, bookID, number)
}
