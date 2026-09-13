package ports

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type OutboxRepository interface {
	Publish(ctx context.Context, event *domain.Event) error
}
