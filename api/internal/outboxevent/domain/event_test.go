package domain_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

func TestNewEvent(t *testing.T) {
	createdAt := time.Date(2026, time.September, 13, 10, 0, 0, 0, time.UTC)
	availableAt := createdAt.Add(time.Minute)
	payload := []byte("ciphertext")

	event, err := domain.NewEvent(
		"event-1",
		"identityaccess.magic_link_requested",
		payload,
		availableAt,
		createdAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	if event.ID() != "event-1" {
		t.Fatalf("expected id %q, got %q", "event-1", event.ID())
	}
	if event.EventType() != "identityaccess.magic_link_requested" {
		t.Fatalf("unexpected event type %q", event.EventType())
	}
	if !bytes.Equal(event.Payload(), payload) {
		t.Fatalf("unexpected payload %q", event.Payload())
	}
	if !event.AvailableAt().Equal(availableAt) {
		t.Fatalf("expected available at %s, got %s", availableAt, event.AvailableAt())
	}
	if !event.CreatedAt().Equal(createdAt) {
		t.Fatalf("expected created at %s, got %s", createdAt, event.CreatedAt())
	}
	if event.Attempts() != 0 {
		t.Fatalf("expected zero attempts, got %d", event.Attempts())
	}
	if event.LeaseToken() != "" || !event.LeasedUntil().IsZero() {
		t.Fatal("new event must not have a lease")
	}
	if !event.ProcessedAt().IsZero() || !event.FailedAt().IsZero() {
		t.Fatal("new event must be pending")
	}
}

func TestNewEventRejectsInvalidData(t *testing.T) {
	now := time.Date(2026, time.September, 13, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		id          string
		eventType   string
		payload     []byte
		availableAt time.Time
		createdAt   time.Time
	}{
		{
			name:        "missing id",
			eventType:   "identityaccess.magic_link_requested",
			payload:     []byte("payload"),
			availableAt: now,
			createdAt:   now,
		},
		{
			name:        "missing event type",
			id:          "event-1",
			payload:     []byte("payload"),
			availableAt: now,
			createdAt:   now,
		},
		{
			name:        "empty payload",
			id:          "event-1",
			eventType:   "identityaccess.magic_link_requested",
			availableAt: now,
			createdAt:   now,
		},
		{
			name:        "available before creation",
			id:          "event-1",
			eventType:   "identityaccess.magic_link_requested",
			payload:     []byte("payload"),
			availableAt: now.Add(-time.Nanosecond),
			createdAt:   now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := domain.NewEvent(
				tt.id,
				tt.eventType,
				tt.payload,
				tt.availableAt,
				tt.createdAt,
			)
			if err == nil {
				t.Fatalf("expected validation error, got event %+v", event)
			}
		})
	}
}

func TestEventDoesNotExposePayloadBackingArray(t *testing.T) {
	now := time.Date(2026, time.September, 13, 10, 0, 0, 0, time.UTC)
	payload := []byte("ciphertext")
	event, err := domain.NewEvent("event-1", "event.type", payload, now, now)
	if err != nil {
		t.Fatal(err)
	}

	payload[0] = 'X'
	if bytes.Equal(event.Payload(), payload) {
		t.Fatal("event retained the caller payload backing array")
	}

	returned := event.Payload()
	returned[0] = 'Y'
	if bytes.Equal(event.Payload(), returned) {
		t.Fatal("event exposed its payload backing array")
	}
}
