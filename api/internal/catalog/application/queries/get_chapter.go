package queries

import (
	"context"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"strings"
)

type GetChapter struct{ repository ports.CatalogRepository }

func NewGetChapter(repository ports.CatalogRepository) *GetChapter {
	return &GetChapter{repository: repository}
}
func (q *GetChapter) Execute(ctx context.Context, bookID string, number int) (domain.Chapter, error) {
	if strings.TrimSpace(bookID) == "" {
		return domain.Chapter{}, domain.ErrInvalidBookID
	}
	if number <= 0 {
		return domain.Chapter{}, domain.ErrInvalidChapter
	}
	return q.repository.FindChapter(ctx, bookID, number)
}
