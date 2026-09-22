CREATE TABLE books (
    id TEXT PRIMARY KEY,
    "order" SMALLINT NOT NULL UNIQUE CHECK ("order" > 0),
    name VARCHAR(255) NOT NULL UNIQUE CHECK (length(trim(name)) > 0),
    testament VARCHAR(3) NOT NULL CHECK (testament IN ('old', 'new')),
    chapter_count SMALLINT NOT NULL CHECK (chapter_count > 0)
);

CREATE TABLE verses (
    book_id TEXT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    chapter SMALLINT NOT NULL CHECK (chapter > 0),
    number SMALLINT NOT NULL CHECK (number > 0),
    text TEXT NOT NULL CHECK (length(trim(text)) > 0),
    part SMALLINT NOT NULL CHECK (part > 0),
    PRIMARY KEY (book_id, chapter, number, part)
);
