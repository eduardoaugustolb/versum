package application

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

// EventHandler processes a delivery. Implementations must tolerate duplicates:
// a process can stop after the side effect and before acknowledging the event.
// Use the event ID as the idempotency key when the destination supports it.
type EventHandler interface {
	Handle(context.Context, *domain.Event) error
}
