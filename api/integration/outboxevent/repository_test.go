package outboxevent_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	adapterpostgres "github.com/eduardoaugustolb/versum/api/internal/database/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	outboxeventpg "github.com/eduardoaugustolb/versum/api/internal/outboxevent/adapters/postgres"
	ports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
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
	return bytes.TrimPrefix(protected.Ciphertext, []byte("protected:")), nil
}

func TestConcurrentDeliveryClaimsLeaseRecoveryAndPayloadErasure(t *testing.T) {
	ctx := t.Context()
	publisher, db, _ := setupOutboxRepository(ctx, t)
	delivery := outboxeventpg.NewEventDeliveryRepository(outboxeventpg.NewDAO(db), testPayloadProtector{}, id.UUIDGenerator{})
	// Integration suites share tables and must run serially against a disposable DB.
	if err := db.Exec(ctx, "DELETE FROM outbox_events"); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"delivery-a", "delivery-b"} {
		event, _ := domain.NewEvent(id, "delivery.test", "", []byte("secret"))
		if err := publisher.Publish(ctx, event); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now().UTC().Add(time.Second)
	claims := make(chan ports.ClaimedEvent, 2)
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Go(func() { claim, err := delivery.Claim(ctx, now, now.Add(time.Minute)); claims <- claim; errs <- err })
	}
	wg.Wait()
	close(claims)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var acquired []ports.ClaimedEvent
	for claim := range claims {
		acquired = append(acquired, claim)
		if claim.Attempt != 1 || string(claim.Event.Payload()) != "secret" {
			t.Fatal("incorrect claim")
		}
	}
	if acquired[0].Event.ID() == acquired[1].Event.ID() {
		t.Fatal("two workers claimed the same event")
	}
	if _, err := delivery.Claim(ctx, now, now.Add(time.Minute)); !errors.Is(err, ports.ErrNoEventAvailable) {
		t.Fatalf("active leases claimed: %v", err)
	}
	first := acquired[0]
	if err := delivery.Complete(ctx, first.Event.ID(), "wrong-lease", now); !errors.Is(err, ports.ErrEventLeaseLost) {
		t.Fatalf("wrong lease accepted: %v", err)
	}
	if err := delivery.Complete(ctx, first.Event.ID(), first.LeaseToken, now); err != nil {
		t.Fatal(err)
	}
	second := acquired[1]
	later := now.Add(time.Minute)
	recovered, err := delivery.Claim(ctx, later, later.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Event.ID() != second.Event.ID() || recovered.Attempt != 2 || recovered.LeaseToken == second.LeaseToken {
		t.Fatal("expired lease not recovered correctly")
	}
	if err := delivery.Complete(ctx, second.Event.ID(), second.LeaseToken, later); !errors.Is(err, ports.ErrEventLeaseLost) {
		t.Fatalf("stale owner accepted: %v", err)
	}
	if err := delivery.Retry(ctx, recovered.Event.ID(), recovered.LeaseToken, later, later.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := delivery.Claim(ctx, later, later.Add(time.Minute)); !errors.Is(err, ports.ErrNoEventAvailable) {
		t.Fatalf("retried too early: %v", err)
	}
	last, err := delivery.Claim(ctx, later.Add(time.Second), later.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if last.Attempt != 3 {
		t.Fatal("retry attempt not counted")
	}
	if err := delivery.Fail(ctx, last.Event.ID(), last.LeaseToken, later.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	var retained int
	if err := db.QueryRow(ctx, "SELECT count(*) FROM outbox_events WHERE octet_length(payload_ciphertext) > 0").Scan(&retained); err != nil {
		t.Fatal(err)
	}
	if retained != 0 {
		t.Fatal("terminal deliveries retained secret payloads")
	}
	if _, err := delivery.Claim(ctx, later.Add(2*time.Minute), later.Add(3*time.Minute)); !errors.Is(err, ports.ErrNoEventAvailable) {
		t.Fatalf("terminal event reclaimed: %v", err)
	}
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
	return outboxeventpg.NewOutboxRepository(outboxeventpg.NewDAO(db), testPayloadProtector{}), db, pool
}

func TestRepositoryPublishesProtectedEvent(t *testing.T) {
	ctx := t.Context()
	repository, db, _ := setupOutboxRepository(ctx, t)

	event, err := domain.NewEvent(
		"outbox-event-1",
		"identityaccess.magic_link_requested", "",
		[]byte("raw-secret-payload"),
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

	if id != event.ID() || eventType != event.EventType().String() {
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
	if storedAvailableAt.IsZero() || storedCreatedAt.IsZero() {
		t.Fatal("expected database timestamps")
	}
}

func TestDeletingUserCascadesOwnedEventsAndPreservesUnownedEvents(t *testing.T) {
	ctx := t.Context()
	_, db, _ := setupOutboxRepository(ctx, t)
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if err := tx.Exec(ctx, `INSERT INTO users (id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version) VALUES ($1, $2, $3, 1, 1)`, "outbox-cascade-user", []byte("email"), []byte("outbox-cascade-lookup")); err != nil {
		t.Fatal(err)
	}
	repo := outboxeventpg.NewOutboxRepository(outboxeventpg.NewDAO(tx), testPayloadProtector{})
	for _, owner := range []string{"outbox-cascade-user", ""} {
		event, err := domain.NewEvent("cascade-event-"+owner, "event.type", owner, []byte("payload"))
		if err != nil {
			t.Fatal(err)
		}
		if err := repo.Publish(ctx, event); err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, "outbox-cascade-user"); err != nil {
		t.Fatal(err)
	}
	var ownedCount, unownedCount int
	if err := tx.QueryRow(ctx, `SELECT count(*) FILTER (WHERE id = 'cascade-event-outbox-cascade-user'), count(*) FILTER (WHERE id = 'cascade-event-' AND user_id IS NULL) FROM outbox_events`).Scan(&ownedCount, &unownedCount); err != nil {
		t.Fatal(err)
	}
	if ownedCount != 0 || unownedCount != 1 {
		t.Fatalf("expected owned event deleted and unowned event preserved; got %d and %d", ownedCount, unownedCount)
	}
}
