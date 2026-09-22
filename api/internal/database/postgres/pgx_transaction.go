package postgres

import (
	"context"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/jackc/pgx/v5"
)

type pgxTransaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	pgxConn
}

type PgxTransaction struct {
	tx pgxTransaction
}

var _ dbexec.Transaction = PgxTransaction{}

func NewPgxTransaction(tx pgxTransaction) PgxTransaction {
	return PgxTransaction{tx: tx}
}

func (e PgxTransaction) QueryRow(ctx context.Context, sql string, args ...any) dbexec.Row {
	return e.tx.QueryRow(ctx, sql, args...)
}

func (e PgxTransaction) Query(ctx context.Context, sql string, args ...any) (dbexec.Rows, error) {
	return e.tx.Query(ctx, sql, args...)
}

func (e PgxTransaction) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := e.tx.Exec(ctx, sql, args...)
	return err
}

func (e PgxTransaction) CopyFrom(ctx context.Context, table string, columns []string, rows [][]any) (int64, error) {
	return e.tx.CopyFrom(ctx, pgx.Identifier{table}, columns, pgx.CopyFromRows(rows))
}

// Begin opens a nested savepoint.
func (e PgxTransaction) Begin(ctx context.Context) (dbexec.Transaction, error) {
	tx, err := e.tx.Begin(ctx)
	dbexecTx := NewPgxTransaction(tx)
	return dbexecTx, err
}

func (e PgxTransaction) Commit(ctx context.Context) error {
	return e.tx.Commit(ctx)
}

func (e PgxTransaction) Rollback(ctx context.Context) error {
	return e.tx.Rollback(ctx)
}
