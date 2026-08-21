package postgres

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

type Repository struct {
	db dbexec.Executor
}

func NewRepository(db dbexec.Executor) *Repository {
	return &Repository{db: db}
}

var verseColumns = []string{"book_id", "chapter", "number", "text", "part"}

// ReplaceBook publishes a book and its verses atomically from the caller's
// point of view: it deletes the book's existing verses, upserts the book,
// then loads the new verses via CopyFrom. Deleting verses before touching
// the book avoids an FK violation, and replacing (not upserting) verses
// means a verse removed or renumbered in a corpus revision doesn't linger
// as an orphan row. CopyFrom streams every row in one pass instead of one
// Exec per verse — a book like Psalms (~2,461 verses) would otherwise cost
// ~2,461 network round trips, and the full corpus is ~35,624 verses. The
// caller (the seed command) is expected to run this inside a transaction
// so a partial failure doesn't leave the book half-published.
func (r *Repository) ReplaceBook(ctx context.Context, book domain.Book, verses []domain.Verse) error {
	if err := r.db.Exec(ctx, DeleteBookVersesQuery, book.ID()); err != nil {
		return err
	}
	if err := r.db.Exec(ctx, UpsertBookQuery, book.ID(), book.Order(), book.Name(), book.Testament(), book.ChapterCount()); err != nil {
		return err
	}

	rows := make([][]any, len(verses))
	for i, verse := range verses {
		rows[i] = []any{book.ID(), verse.Chapter(), verse.Number(), verse.Text(), verse.Part()}
	}

	_, err := r.db.CopyFrom(ctx, "verses", verseColumns, rows)
	return err
}

func (r *Repository) ListBooks(ctx context.Context) ([]domain.Book, error) {
	books := []domain.Book{}
	rows, err := r.db.Query(ctx, ListBooksQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, testament string
		var order, chapterCount int
		if err := rows.Scan(&id, &order, &name, &testament, &chapterCount); err != nil {
			return nil, err
		}
		book, err := domain.NewBook(domain.NewBookParams{ID: id, Order: order, Name: name, Testament: domain.Testament(testament), ChapterCount: chapterCount})
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) FindChapter(ctx context.Context, bookID string, number int) (domain.Chapter, error) {
	verses := []domain.Verse{}
	bookName := ""
	rows, err := r.db.Query(ctx, FindChapterVersesQuery, bookID, number)
	if err != nil {
		return domain.Chapter{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var verseBookID, text string
		var chapter, verseNumber, part int
		if err := rows.Scan(&verseBookID, &chapter, &verseNumber, &text, &part, &bookName); err != nil {
			return domain.Chapter{}, err
		}
		verse, err := domain.NewVerse(domain.NewVerseParams{BookID: verseBookID, Chapter: chapter, Number: verseNumber, Text: text, Part: part})
		if err != nil {
			return domain.Chapter{}, err
		}
		verses = append(verses, verse)
	}

	if err := rows.Err(); err != nil {
		return domain.Chapter{}, err
	}

	if len(verses) == 0 {
		return domain.Chapter{}, domain.ErrChapterNotFound
	}

	return domain.NewChapter(bookID, bookName, number, verses), nil
}
