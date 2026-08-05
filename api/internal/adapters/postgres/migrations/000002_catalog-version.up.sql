CREATE TABLE catalog_version (
    id            SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    corpus_sha256 TEXT NOT NULL CHECK (length(trim(corpus_sha256)) > 0),
    published_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
