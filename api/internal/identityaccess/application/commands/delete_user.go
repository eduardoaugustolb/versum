package commands

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
)

type DeleteUser struct {
	unitOfWork application.IdentityAccessUnitOfWork
	clock      clock.Clock
}

func NewDeleteUser(uow application.IdentityAccessUnitOfWork, clock clock.Clock) *DeleteUser {
	return &DeleteUser{uow, clock}
}
func (uc *DeleteUser) Execute(ctx context.Context, secret string) error {
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
		if !session.IsValidAt(now) {
			return application.ErrInvalidSessionSecret
		}
		return repos.Users.DeleteUser(ctx, session.UserID())
	})
}
