package commands

import (
	"context"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
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
	book, verses, err := buildPublication(input)
	if err != nil {
		return err
	}
	return uc.transactions.WithinTransaction(ctx, func(txctx context.Context, repository ports.CatalogRepository) error {
		return repository.ReplaceBook(txctx, book, verses)
	})
}

func buildPublication(input PublishBookInput) (domain.Book, []domain.Verse, error) {
	book, err := domain.NewBook(domain.NewBookParams{ID: input.ID, Order: input.Order, Name: input.Name, Testament: input.Testament, ChapterCount: input.ChapterCount})
	if err != nil {
		return domain.Book{}, nil, err
	}
	verses := make([]domain.Verse, 0, len(input.Verses))
	for _, inputVerse := range input.Verses {
		verse, err := domain.NewVerse(domain.NewVerseParams{BookID: inputVerse.BookID, Chapter: inputVerse.Chapter, Number: inputVerse.Number, Text: inputVerse.Text, Part: inputVerse.Part})
		if err != nil {
			return domain.Book{}, nil, err
		}
		verses = append(verses, verse)
	}
	if err := domain.ValidateBookVerses(book, verses); err != nil {
		return domain.Book{}, nil, err
	}
	return book, verses, nil
}
