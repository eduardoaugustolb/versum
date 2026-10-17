package commands

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
)

type RevokeSessionFamily struct {
	unitOfWork application.IdentityAccessUnitOfWork
	clock      clock.Clock
}

func NewRevokeSessionFamily(unitOfWork application.IdentityAccessUnitOfWork, clock clock.Clock) *RevokeSessionFamily {
	return &RevokeSessionFamily{unitOfWork: unitOfWork, clock: clock}
}

func (uc *RevokeSessionFamily) Execute(ctx context.Context, secret string) error {
	return uc.unitOfWork.WithinTransaction(ctx, func(repos application.IdentityAccessUnitOfWorkRepositories) error {
		session, err := repos.Sessions.FindSessionBySecret(ctx, secret)
		if err != nil {
			return err
		}

		if session == nil {
			return application.ErrSessionNotFound
		}

		if err := repos.Sessions.LockSessionsByUserID(ctx, session.UserID()); err != nil {
			return err
		}
		session, err = repos.Sessions.FindSessionBySecret(ctx, secret)
		if err != nil {
			return err
		}
		if session == nil {
			return application.ErrSessionNotFound
		}
		now := uc.clock.Now().UTC()
		if session.IsExpired(now) {
			return application.ErrSessionExpired
		}

		if err := repos.Sessions.RevokeSessionFamily(ctx, session.UserID(), session.FamilyID(), now); err != nil {
			return err
		}

		return nil
	})
}
