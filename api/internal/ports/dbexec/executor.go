// Package dbexec is the execution port between a domain's SQL-backed
// repository and whatever concrete driver runs the query. Its types never
// reference a third-party driver package — that is the whole point: the
// abstraction must not depend on the detail it abstracts over.
package dbexec

import "context"

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Row
	Next() bool
	Err() error
	Close()
}

type Executor interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	Exec(ctx context.Context, sql string, args ...any) error

	// CopyFrom bulk-loads rows into table using a native streaming
	// mechanism (e.g. Postgres COPY) instead of one statement per row.
	// table and columns are plain identifiers — no driver type leaks into
	// the signature, same as the rest of this port.
	CopyFrom(ctx context.Context, table string, columns []string, rows [][]any) (int64, error)
}
