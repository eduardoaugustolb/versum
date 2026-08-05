package catalog

import "context"

type BookTestament string

const (
	TestamentOld = BookTestament("old")
	TestamentNew = BookTestament("new")
)

type Book struct {
	ID           string        `json:"id"`
	Order        int           `json:"order"`
	Name         string        `json:"name"`
	Testament    BookTestament `json:"testament"`
	ChapterCount int           `json:"chapter_count"`
}

type BookRepository interface {
	ListBooks(ctx context.Context) ([]Book, error)
}
