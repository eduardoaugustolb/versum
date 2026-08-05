package catalog

import "context"

type CatalogWriter interface {
	ReplaceBook(ctx context.Context, book Book, verses []Verse) error
}

type PublishBook struct {
	writer CatalogWriter
}

func NewPublishBook(writer CatalogWriter) PublishBook {
	return PublishBook{writer: writer}
}

func (uc PublishBook) Execute(ctx context.Context, book Book, verses []Verse) error {
	return uc.writer.ReplaceBook(ctx, book, verses)
}
