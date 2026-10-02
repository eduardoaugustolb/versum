package commands

import (
	"context"
	"fmt"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type RotateSession struct {
	unitOfWork      application.IdentityAccessUnitOfWork
	secretGenerator application.SessionSecretGenerator
	idGenerator     id.Generator
	clock           clock.Clock
}

func NewRotateSession(uow application.IdentityAccessUnitOfWork, secrets application.SessionSecretGenerator, ids id.Generator, clock clock.Clock) *RotateSession {
	return &RotateSession{unitOfWork: uow, secretGenerator: secrets, idGenerator: ids, clock: clock}
}

func (uc *RotateSession) Execute(ctx context.Context, secret string, ipAddress, userAgent string) (*IssuedSession, error) {
	var result *IssuedSession
	var reused bool
	err := uc.unitOfWork.WithinTransaction(ctx, func(repos application.IdentityAccessUnitOfWorkRepositories) error {
		current, err := repos.Sessions.FindSessionBySecret(ctx, secret)
		if err != nil {
			return fmt.Errorf("finding session: %w", err)
		}
		if err := repos.Sessions.LockSessionsByUserID(ctx, current.UserID()); err != nil {
			return err
		}
		// Re-read after the lock: another transaction may have rotated the credential.
		current, err = repos.Sessions.FindSessionBySecret(ctx, secret)
		if err != nil {
			return err
		}
		now := uc.clock.Now().UTC()
		if _, replaced := current.ReplacedBySessionID(); replaced {
			if err := repos.Sessions.RevokeSessionFamily(ctx, current.UserID(), current.FamilyID(), now); err != nil {
				return err
			}
			reused = true
			return nil // Commit the revocation before reporting reuse to the caller.
		}
		if !current.IsValidAt(now) {
			return domain.ErrInvalidSessionState
		}
		nextSecret, err := uc.secretGenerator.GenerateSessionSecret()
		if err != nil {
			return fmt.Errorf("generating session secret: %w", err)
		}
		if nextSecret == "" || nextSecret == secret {
			return application.ErrInvalidSessionSecret
		}
		next, err := domain.NewSession(uc.idGenerator.Generate(), current.UserID(), current.FamilyID(), ipAddress, userAgent, current.ExpiresAt(), now)
		if err != nil {
			return err
		}
		if err := current.ReplaceWith(next, now); err != nil {
			return err
		}
		if err := repos.Sessions.CreateSession(ctx, next, nextSecret); err != nil {
			return err
		}
		if err := repos.Sessions.SaveRotation(ctx, current); err != nil {
			return err
		}
		result = &IssuedSession{Session: next, Secret: nextSecret}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if reused {
		return nil, domain.ErrSessionReused
	}
	return result, nil
}
