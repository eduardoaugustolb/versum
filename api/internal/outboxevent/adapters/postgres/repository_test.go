package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	ports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type outboxExecutor struct{ err error }

func (e outboxExecutor) Exec(context.Context, string, ...any) error               { return e.err }
func (outboxExecutor) QueryRow(context.Context, string, ...any) dbexec.Row        { return nil }
func (outboxExecutor) Query(context.Context, string, ...any) (dbexec.Rows, error) { return nil, nil }
func (outboxExecutor) CopyFrom(context.Context, string, []string, [][]any) (int64, error) {
	return 0, nil
}
func (outboxExecutor) Begin(context.Context) (dbexec.Transaction, error) { return nil, nil }

type outboxPayloadProtector struct{}

func (outboxPayloadProtector) Protect(context.Context, []byte) (ports.ProtectedPayload, error) {
	return ports.ProtectedPayload{Ciphertext: []byte("ciphertext"), KeyVersion: 1}, nil
}

func (outboxPayloadProtector) Unprotect(context.Context, ports.ProtectedPayload) ([]byte, error) {
	return nil, nil
}

func TestOutboxRepositoryTranslatesUniqueViolationOnPublish(t *testing.T) {
	event, err := domain.NewEvent("event-1", "identityaccess.magic_link_requested", []byte("payload"))
	if err != nil {
		t.Fatal(err)
	}
	repo := NewOutboxRepository(outboxExecutor{err: &pgconn.PgError{Code: "23505"}}, outboxPayloadProtector{})

	if err := repo.Publish(t.Context(), event); !errors.Is(err, application.ErrOutboxEventAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}
