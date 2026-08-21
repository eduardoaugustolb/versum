package queries

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type GetChapter struct{ reader ports.ChapterReader }

func NewGetChapter(reader ports.ChapterReader) *GetChapter { return &GetChapter{reader: reader} }
func (q *GetChapter) Execute(ctx context.Context, bookID string, number int) (domain.Chapter, error) {
	return q.reader.FindChapter(ctx, bookID, number)
}
