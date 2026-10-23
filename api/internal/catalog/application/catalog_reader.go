package application

import "context"

type BookDTO struct {
	ID           string
	Order        int
	Name         string
	Testament    string
	ChapterCount int
}
type VerseDTO struct {
	BookID  string
	Chapter int
	Number  int
	Text    string
	Part    int
}
type ChapterDTO struct {
	BookID   string
	BookName string
	Number   int
	Verses   []VerseDTO
}

// CatalogReader is the DAO port for catalog projections. It never hydrates aggregates.
type CatalogReader interface {
	ListBooks(context.Context) ([]BookDTO, error)
	FindChapter(context.Context, string, int) (ChapterDTO, error)
}
