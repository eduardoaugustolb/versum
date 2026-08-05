package catalog

import (
	"context"
	"errors"
)

var ErrChapterNotFound = errors.New("chapter not found")

type Chapter struct {
	BookID   string  `json:"book_id"`
	BookName string  `json:"book_name"`
	Number   int     `json:"number"`
	Verses   []Verse `json:"verses"`
}

type ChapterRepository interface {
	FindChapter(ctx context.Context, bookID string, number int) (Chapter, error)
}
