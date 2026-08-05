package catalog

import "context"

type GetChapter struct {
	repo ChapterRepository
}

func NewGetChapter(repo ChapterRepository) *GetChapter {
	return &GetChapter{repo: repo}
}
func (uc *GetChapter) Execute(ctx context.Context, bookID string, number int) (Chapter, error) {
	return uc.repo.FindChapter(ctx, bookID, number)
}
