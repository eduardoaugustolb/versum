package commands

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
)

type PublishBookInput struct {
	ID           string
	Order        int
	Name         string
	Testament    domain.Testament
	ChapterCount int
	Verses       []VerseInput
}
type VerseInput struct {
	BookID          string
	Chapter, Number int
	Text            string
	Part            int
}

type PublishBook struct{ transactions ports.TransactionManager }

func NewPublishBook(transactions ports.TransactionManager) PublishBook {
	return PublishBook{transactions: transactions}
}
func (uc PublishBook) Execute(ctx context.Context, input PublishBookInput) error {
	book, err := domain.NewBook(domain.NewBookParams{ID: input.ID, Order: input.Order, Name: input.Name, Testament: input.Testament, ChapterCount: input.ChapterCount})
	if err != nil {
		return err
	}
	verses := make([]domain.Verse, 0, len(input.Verses))
	for _, inputVerse := range input.Verses {
		verse, err := domain.NewVerse(domain.NewVerseParams{BookID: inputVerse.BookID, Chapter: inputVerse.Chapter, Number: inputVerse.Number, Text: inputVerse.Text, Part: inputVerse.Part})
		if err != nil {
			return err
		}
		verses = append(verses, verse)
	}
	return uc.transactions.WithinTransaction(ctx, func(txctx context.Context, writer ports.CatalogWriter) error {
		return writer.ReplaceBook(txctx, book, verses)
	})
}
