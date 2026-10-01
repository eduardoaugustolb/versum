package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type ConsumeMagicLink struct {
	unitOfWork      application.IdentityAccessUnitOfWork
	secretGenerator application.SessionSecretGenerator
	idGenerator     id.Generator
	clock           clock.Clock
	sessionTTL      time.Duration
}

func NewConsumeMagicLink(uow application.IdentityAccessUnitOfWork, secrets application.SessionSecretGenerator, ids id.Generator, clock clock.Clock, sessionTTL time.Duration) (*ConsumeMagicLink, error) {
	if sessionTTL <= 0 {
		return nil, application.ErrInvalidSessionTTL
	}
	return &ConsumeMagicLink{unitOfWork: uow, secretGenerator: secrets, idGenerator: ids, clock: clock, sessionTTL: sessionTTL}, nil
}

// Execute consumes a single-use link and issues an independent session family.
// The bearer credential is returned only after the transaction commits.
func (uc *ConsumeMagicLink) Execute(ctx context.Context, token string, client domain.SessionClient) (*IssuedSession, error) {
	if token == "" {
		return nil, application.ErrInvalidLoginTokenSecret
	}
	var result *IssuedSession
	err := uc.unitOfWork.WithinTransaction(ctx, func(repos application.IdentityAccessUnitOfWorkRepositories) error {
		loginToken, err := repos.LoginTokens.FindByToken(ctx, token)
		if err != nil {
			return fmt.Errorf("finding login token: %w", err)
		}
		now := uc.clock.Now().UTC()
		if err := loginToken.Consume(now); err != nil {
			return err
		}
		secret, err := uc.secretGenerator.GenerateSessionSecret()
		if err != nil {
			return fmt.Errorf("generating session secret: %w", err)
		}
		if secret == "" || secret == token {
			return application.ErrInvalidSessionSecret
		}
		sessionID := uc.idGenerator.Generate()
		familyID := uc.idGenerator.Generate()
		session, err := domain.NewSession(sessionID, loginToken.UserID(), familyID, client, now.Add(uc.sessionTTL), now)
		if err != nil {
			return err
		}
		if err := repos.LoginTokens.Save(ctx, loginToken); err != nil {
			return fmt.Errorf("consuming login token: %w", err)
		}
		if err := repos.Sessions.CreateSession(ctx, session, secret); err != nil {
			return fmt.Errorf("creating session: %w", err)
		}
		result = &IssuedSession{Session: session, Secret: secret}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
