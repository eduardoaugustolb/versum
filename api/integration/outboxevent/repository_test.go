package outboxevent_test

import (
	"context"
	"os"
	"testing"
	"time"

	adapterpostgres "github.com/eduardoaugustolb/versum/api/internal/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
	outboxeventpg "github.com/eduardoaugustolb/versum/api/internal/outboxevent/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
	"github.com/jackc/pgx/v5/pgxpool"
)

type testPayloadProtector struct{}

func (testPayloadProtector) Protect(_ context.Context, payload []byte) (ports.ProtectedPayload, error) {
	return ports.ProtectedPayload{
		Ciphertext: append([]byte("protected:"), payload...),
		KeyVersion: 7,
	}, nil
}

func (testPayloadProtector) Unprotect(_ context.Context, protected ports.ProtectedPayload) ([]byte, error) {
	return protected.Ciphertext, nil
}

func setupOutboxRepository(ctx context.Context, t *testing.T) (*outboxeventpg.OutboxRepository, dbexec.Executor, *pgxpool.Pool) {
	t.Helper()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	db := adapterpostgres.NewPgxExecutor(pool)
	return outboxeventpg.NewOutboxRepository(db, testPayloadProtector{}), db, pool
}

func TestRepositoryPublishesProtectedEvent(t *testing.T) {
	ctx := t.Context()
	repository, db, _ := setupOutboxRepository(ctx, t)

	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	availableAt := createdAt.Add(time.Minute)
	event, err := domain.NewEvent(
		"outbox-event-1",
		"identityaccess.magic_link_requested",
		[]byte("raw-secret-payload"),
		availableAt,
		createdAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(ctx, "DELETE FROM outbox_events WHERE id = $1", event.ID()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Exec(context.Background(), "DELETE FROM outbox_events WHERE id = $1", event.ID())
	})

	if err := repository.Publish(ctx, event); err != nil {
		t.Fatal(err)
	}

	var (
		id, eventType                      string
		payloadCiphertext                  []byte
		payloadKeyVersion, attempts        int
		storedAvailableAt, storedCreatedAt time.Time
		leaseToken, lastErrorRedacted      *string
		leasedUntil, processedAt, failedAt *time.Time
	)
	err = db.QueryRow(ctx, `
		SELECT id, event_type, payload_ciphertext, payload_key_version,
		       attempts, available_at, lease_token, leased_until,
		       processed_at, failed_at, last_error_redacted, created_at
		FROM outbox_events
		WHERE id = $1
	`, event.ID()).Scan(
		&id,
		&eventType,
		&payloadCiphertext,
		&payloadKeyVersion,
		&attempts,
		&storedAvailableAt,
		&leaseToken,
		&leasedUntil,
		&processedAt,
		&failedAt,
		&lastErrorRedacted,
		&storedCreatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	if id != event.ID() || eventType != event.EventType() {
		t.Fatalf("unexpected persisted event identity: id=%q type=%q", id, eventType)
	}
	if string(payloadCiphertext) != "protected:raw-secret-payload" {
		t.Fatalf("expected protected payload, got %q", payloadCiphertext)
	}
	if string(payloadCiphertext) == string(event.Payload()) {
		t.Fatal("adapter persisted the raw payload")
	}
	if payloadKeyVersion != 7 {
		t.Fatalf("expected payload key version 7, got %d", payloadKeyVersion)
	}
	if attempts != 0 || leaseToken != nil || leasedUntil != nil || processedAt != nil || failedAt != nil || lastErrorRedacted != nil {
		t.Fatal("new persisted event must be pending and unleased")
	}
	if !storedAvailableAt.Equal(event.AvailableAt()) || !storedCreatedAt.Equal(event.CreatedAt()) {
		t.Fatalf("unexpected persisted timestamps: available_at=%s created_at=%s", storedAvailableAt, storedCreatedAt)
	}
}
