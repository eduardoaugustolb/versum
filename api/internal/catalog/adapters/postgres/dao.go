package postgres

import (
	"context"

	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
)

type DAO struct {
	db dbexec.Executor
}

func NewDAO(db dbexec.Executor) *DAO {
	return &DAO{db: db}
}

var _ ports.CatalogReader = (*DAO)(nil)

var verseColumns = []string{"book_id", "chapter", "number", "text", "part"}

// ReplaceBook publishes a book and its verses atomically from the caller's
// point of view: it deletes the book's existing verses, upserts the book,
// then loads the new verses via CopyFrom. Deleting verses before touching
// the book avoids an FK violation, and replacing (not upserting) verses
// means a verse removed or renumbered in a corpus revision doesn't linger
// as an orphan row. CopyFrom streams every row in one pass instead of one
// Exec per verse — a book like Psalms (~2,461 verses) would otherwise cost
// ~2,461 network round trips, and the full corpus is ~35,624 verses. The
// caller (the publication use case) is expected to run this inside a transaction
// so a partial failure doesn't leave the book half-published.
func (r *DAO) ReplaceBook(ctx context.Context, book ports.BookDTO, verses []ports.VerseDTO) error {
	if err := r.db.Exec(ctx, DeleteBookVersesQuery, book.ID); err != nil {
		return err
	}
	if err := r.db.Exec(ctx, UpsertBookQuery, book.ID, book.Order, book.Name, book.Testament, book.ChapterCount); err != nil {
		return err
	}

	rows := make([][]any, len(verses))
	for i, verse := range verses {
		rows[i] = []any{verse.BookID, verse.Chapter, verse.Number, verse.Text, verse.Part}
	}

	_, err := r.db.CopyFrom(ctx, "verses", verseColumns, rows)
	return err
}

func (r *DAO) ListBooks(ctx context.Context) ([]ports.BookDTO, error) {
	books := []ports.BookDTO{}
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
		book := ports.BookDTO{ID: id, Order: order, Name: name, Testament: testament, ChapterCount: chapterCount}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *DAO) FindChapter(ctx context.Context, bookID string, number int) (ports.ChapterDTO, error) {
	verses := []ports.VerseDTO{}
	bookName := ""
	rows, err := r.db.Query(ctx, FindChapterVersesQuery, bookID, number)
	if err != nil {
		return ports.ChapterDTO{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var verseBookID, text string
		var chapter, verseNumber, part int
		if err := rows.Scan(&verseBookID, &chapter, &verseNumber, &text, &part, &bookName); err != nil {
			return ports.ChapterDTO{}, err
		}
		verse := ports.VerseDTO{BookID: verseBookID, Chapter: chapter, Number: verseNumber, Text: text, Part: part}
		verses = append(verses, verse)
	}

	if err := rows.Err(); err != nil {
		return ports.ChapterDTO{}, err
	}

	if len(verses) == 0 {
		return ports.ChapterDTO{}, ports.ErrChapterNotFound
	}

	return ports.ChapterDTO{BookID: bookID, BookName: bookName, Number: number, Verses: verses}, nil
}

func (d *DAO) RecordVersion(ctx context.Context, hash string) error {
	return d.db.Exec(ctx, UpsertCatalogVersionQuery, hash)
}
