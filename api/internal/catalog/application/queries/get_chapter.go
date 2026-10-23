package queries

import (
	"context"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"strings"
)

type GetChapter struct{ repository ports.CatalogReader }

func NewGetChapter(repository ports.CatalogReader) *GetChapter {
	return &GetChapter{repository: repository}
}
func (q *GetChapter) Execute(ctx context.Context, bookID string, number int) (ports.ChapterDTO, error) {
	if strings.TrimSpace(bookID) == "" {
		return ports.ChapterDTO{}, domain.ErrInvalidBookID
	}
	if number <= 0 {
		return ports.ChapterDTO{}, domain.ErrInvalidChapter
	}
	return q.repository.FindChapter(ctx, bookID, number)
}
