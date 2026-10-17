package postgres_test

import (
	"context"
	"errors"
	catalog "github.com/eduardoaugustolb/versum/api/internal/catalog/adapters/postgres"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	db "github.com/eduardoaugustolb/versum/api/internal/database"
	"testing"
)

type transactionDB struct {
	db.Executor
	tx       *transaction
	beginErr error
}

func (d transactionDB) Begin(context.Context) (db.Transaction, error) { return d.tx, d.beginErr }

type transaction struct {
	db.Executor
	committed, rolledBack bool
	commitErr             error
}

func (tx *transaction) Commit(context.Context) error   { tx.committed = true; return tx.commitErr }
func (tx *transaction) Rollback(context.Context) error { tx.rolledBack = true; return nil }
func TestTransactionsAlwaysReleaseResources(t *testing.T) {
	failure := errors.New("failure")
	for _, catalogWide := range []bool{false, true} {
		for _, mode := range []string{"success", "callback error", "panic", "commit error", "begin error"} {
			t.Run(mode+map[bool]string{true: " catalog", false: " book"}[catalogWide], func(t *testing.T) {
				tx := &transaction{}
				db := transactionDB{tx: tx}
				if mode == "begin error" {
					db.beginErr = failure
				}
				if mode == "commit error" {
					tx.commitErr = failure
				}
				manager := catalog.NewTransactionManager(db)
				callback := func() error {
					if mode == "panic" {
						panic(failure)
					}
					if mode == "callback error" {
						return failure
					}
					return nil
				}
				var err error
				var recovered any
				func() {
					defer func() { recovered = recover() }()
					if catalogWide {
						err = manager.WithinCatalogTransaction(t.Context(), func(ports.CatalogRepository, ports.CatalogVersionRepository) error { return callback() })
					} else {
						err = manager.WithinTransaction(t.Context(), func(context.Context, ports.CatalogRepository) error { return callback() })
					}
				}()
				if mode == "panic" && recovered != failure {
					t.Fatal("panic swallowed")
				}
				if mode != "begin error" && !tx.rolledBack {
					t.Fatal("transaction not released")
				}
				if tx.committed != (mode == "success" || mode == "commit error") {
					t.Fatal("incorrect commit")
				}
				if (mode == "callback error" || mode == "begin error" || mode == "commit error") && !errors.Is(err, failure) {
					t.Fatal(err)
				}
			})
		}
	}
}
