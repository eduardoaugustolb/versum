package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

// pgxConn is satisfied by both *pgxpool.Pool and pgx.Tx, so the same
// PgxExecutor works against a pooled connection or a transaction.
type pgxConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type PgxExecutor struct{ conn pgxConn }

func NewPgxExecutor(conn pgxConn) PgxExecutor {
	return PgxExecutor{conn: conn}
}

func (e PgxExecutor) QueryRow(ctx context.Context, sql string, args ...any) dbexec.Row {
	return e.conn.QueryRow(ctx, sql, args...)
}

func (e PgxExecutor) Query(ctx context.Context, sql string, args ...any) (dbexec.Rows, error) {
	return e.conn.Query(ctx, sql, args...)
}

func (e PgxExecutor) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := e.conn.Exec(ctx, sql, args...)
	return err
}

func (e PgxExecutor) CopyFrom(ctx context.Context, table string, columns []string, rows [][]any) (int64, error) {
	return e.conn.CopyFrom(ctx, pgx.Identifier{table}, columns, pgx.CopyFromRows(rows))
}

var _ dbexec.Executor = PgxExecutor{}
