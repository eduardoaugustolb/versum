package domain_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

func TestNewEvent(t *testing.T) {
	payload := []byte("ciphertext")
	event, err := domain.NewEvent("event-1", "identityaccess.magic_link_requested", payload)
	if err != nil {
		t.Fatal(err)
	}
	if event.ID() != "event-1" {
		t.Fatalf("unexpected id %q", event.ID())
	}
	if event.EventType() != "identityaccess.magic_link_requested" {
		t.Fatalf("unexpected event type %q", event.EventType())
	}
	if !bytes.Equal(event.Payload(), payload) {
		t.Fatalf("unexpected payload %q", event.Payload())
	}
}

func TestNewEventRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name, id, eventType string
		payload             []byte
		want                error
	}{
		{"missing id", "", "event.type", []byte("payload"), domain.ErrInvalidId},
		{"missing event type", "event-1", "", []byte("payload"), domain.ErrInvalidEventType},
		{"empty payload", "event-1", "event.type", nil, domain.ErrInvalidPayload},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewEvent(tt.id, tt.eventType, tt.payload)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestEventDoesNotExposePayloadBackingArray(t *testing.T) {
	payload := []byte("ciphertext")
	event, err := domain.NewEvent("event-1", "event.type", payload)
	if err != nil {
		t.Fatal(err)
	}
	payload[0] = 'X'
	if bytes.Equal(event.Payload(), payload) {
		t.Fatal("event retained caller payload backing array")
	}
	returned := event.Payload()
	returned[0] = 'Y'
	if bytes.Equal(event.Payload(), returned) {
		t.Fatal("event exposed payload backing array")
	}
}
