package polling

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type EventProcessor interface {
	ProcessNext(context.Context) (bool, error)
}

// Run drains available work and waits between idle polls or failures. Processing
// is serial per worker; multiple processes coordinate through PostgreSQL leases.
func Run(ctx context.Context, processor EventProcessor, interval time.Duration) error {
	if processor == nil || interval <= 0 {
		return errors.New("invalid polling worker configuration")
	}
	for ctx.Err() == nil {
		processed, err := processor.ProcessNext(ctx)
		if ctx.Err() != nil {
			break
		}
		if err != nil {
			// Provider/database errors may contain secrets. Log only a stable code.
			slog.ErrorContext(ctx, "outbox processing failed", "code", "outbox_processing_failed")
		}
		if processed && err == nil {
			continue
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
	return nil
}
