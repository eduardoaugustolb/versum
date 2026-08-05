package catalog

/*
 * book_id TEXT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
 chapter SMALLINT NOT NULL,
 number SMALLINT NOT NULL,
 text TEXT NOT NULL,
 part SMALLINT NOT NULL,
 PRIMARY KEY (book_id, chapter, number, part)
*/

const (
	ListBooksQuery = `SELECT id, "order", name, testament, chapter_count FROM books ORDER BY "order"`

	FindChapterVersesQuery = `
		SELECT v.book_id, v.chapter, v.number, v.text, v.part, b.name
		FROM verses v
		JOIN books b ON b.id = v.book_id
		WHERE v.book_id = $1 AND v.chapter = $2
		ORDER BY v.number, v.part`
)
