package postgres

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

	DeleteBookVersesQuery = `DELETE FROM verses WHERE book_id = $1`

	UpsertBookQuery = `
		INSERT INTO books (id, "order", name, testament, chapter_count)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			"order" = EXCLUDED."order",
			name = EXCLUDED.name,
			testament = EXCLUDED.testament,
			chapter_count = EXCLUDED.chapter_count`

	UpsertCatalogVersionQuery = `
		INSERT INTO catalog_version (id, corpus_sha256, published_at)
		VALUES (1, $1, now())
		ON CONFLICT (id) DO UPDATE SET
			corpus_sha256 = EXCLUDED.corpus_sha256,
			published_at = EXCLUDED.published_at
	`
)
