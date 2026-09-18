package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityAccessPorts "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxDomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/clock"
	"github.com/eduardoaugustolb/versum/api/internal/ports/id"
)

type RequestMagicLink struct {
	unitOfwork identityAccessPorts.IdentityAccessUnitOfWork
	tokenGenerator identityAccessPorts.LoginTokenGenerator
	tokenHasher    identityAccessPorts.LoginTokenHasher
	idGenerator    id.IDGenerator
	clock          clock.Clock
	magicLinkTTL   time.Duration
}

func NewRequestMagicLink(
	unitOfWork identityAccessPorts.IdentityAccessUnitOfWork,
	tokenGenerator identityAccessPorts.LoginTokenGenerator,
	tokenHasher identityAccessPorts.LoginTokenHasher,
	idGenerator id.IDGenerator,
	clock clock.Clock,
	magicLinkTTL time.Duration,
) (*RequestMagicLink, error) {
	if magicLinkTTL <= 0 {
		return nil, application.ErrInvalidMagicLinkTTL
	}
	return &RequestMagicLink{
		unitOfwork: unitOfWork,
		tokenGenerator: tokenGenerator,
		tokenHasher:    tokenHasher,
		idGenerator:    idGenerator,
		clock:          clock,
		magicLinkTTL:   magicLinkTTL,
	}, nil
}

func (uc *RequestMagicLink) Execute(ctx context.Context, email domain.Email) error {
	return uc.unitOfwork.WithinTransaction(ctx, func(repositories identityAccessPorts.IdentityAccessUnitOfWorkRepositories) error {
		return uc.execute(ctx, email, repositories)
	})
}

func (uc *RequestMagicLink) execute(ctx context.Context, email domain.Email, repositories identityAccessPorts.IdentityAccessUnitOfWorkRepositories) error {
	user, err := repositories.Users.FindUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, application.ErrUserNotFound) {
		return fmt.Errorf("finding user by email: %w", err)
	}

	if errors.Is(err, application.ErrUserNotFound) {
		userID := string(uc.idGenerator.Generate())
		user, err = domain.NewUser(userID, email.String())
		if err != nil {
			return fmt.Errorf("creating user: %w", err)
		}
		if err := repositories.Users.CreateUser(ctx, user); err != nil {
			return fmt.Errorf("saving user: %w", err)
		}
	}

	loginTokenRaw, err := uc.tokenGenerator.GenerateLoginToken()
	if err != nil {
		return fmt.Errorf("generating login token: %w", err)
	}
	loginTokenHash, err := uc.tokenHasher.Hash(loginTokenRaw)
	if err != nil {
		return fmt.Errorf("hashing login token: %w", err)
	}

	id := uc.idGenerator.Generate()
	now := uc.clock.Now().UTC()

	loginToken, err := domain.NewLoginToken(
		string(id),
		loginTokenHash,
		user.ID(),
		now.Add(uc.magicLinkTTL),
		now,
	)

	if err != nil {
		return fmt.Errorf("creating login token: %w", err)
	}

	if err := repositories.LoginTokens.CreateLoginToken(ctx, loginToken); err != nil {
		return fmt.Errorf("saving login token: %w", err)
	}

	eventID := uc.idGenerator.Generate()
	eventType := outboxDomain.EventTypeMagicLinkRequested

	payload, err := json.Marshal(struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{Email: email.String(), Token: loginTokenRaw})
	if err != nil {
		return fmt.Errorf("encoding event payload: %w", err)
	}

	event, err := outboxDomain.NewEvent(eventID.String(), eventType.String(), payload)

	if err != nil {
		return fmt.Errorf("creating event: %w", err)
	}

	err = repositories.Outbox.Publish(ctx, event)
	if err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}
