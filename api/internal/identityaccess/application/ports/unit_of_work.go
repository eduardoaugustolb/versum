package ports

import (
	"context"

	outboxPorts "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application/ports"
)

type IdentityAccessUnitOfWorkRepositories struct {
	Users       UserRepository
	LoginTokens LoginTokenRepository
	Sessions    SessionRepository
	Outbox      outboxPorts.OutboxRepository
}

type IdentityAccessUnitOfWork interface {
	WithinTransaction(
		ctx context.Context,
		fn func(IdentityAccessUnitOfWorkRepositories) error,
	) error
}
