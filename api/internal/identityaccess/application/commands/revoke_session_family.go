package commands

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
)

type RevokeSessionFamily struct {
	sessionRepo application.SessionRepository
	clock       clock.Clock
}

func NewRevokeSessionFamily(sessionRepo application.SessionRepository, clock clock.Clock) *RevokeSessionFamily {
	return &RevokeSessionFamily{sessionRepo: sessionRepo, clock: clock}
}

func (uc *RevokeSessionFamily) Execute(ctx context.Context, secret string) error {
	session, err := uc.sessionRepo.FindSessionBySecret(ctx, secret)
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

	if err := uc.sessionRepo.RevokeSessionFamily(ctx, session.UserID(), session.FamilyID(), now); err != nil {
		return err
	}

	return nil
}
