package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxDomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type RequestMagicLink struct {
	unitOfwork      application.IdentityAccessUnitOfWork
	tokenGenerator  application.LoginTokenGenerator
	tokenHasher     application.LoginTokenHasher
	idGenerator     id.Generator
	clock           clock.Clock
	magicLinkPolicy policy.MagicLinkPolicy
}

func NewRequestMagicLink(
	unitOfWork application.IdentityAccessUnitOfWork,
	tokenGenerator application.LoginTokenGenerator,
	tokenHasher application.LoginTokenHasher,
	idGenerator id.Generator,
	clock clock.Clock,
	magicLinkPolicy policy.MagicLinkPolicy,
) (*RequestMagicLink, error) {
	if magicLinkPolicy.TTL <= 0 {
		return nil, application.ErrInvalidMagicLinkTTL
	}
	return &RequestMagicLink{
		unitOfwork:      unitOfWork,
		tokenGenerator:  tokenGenerator,
		tokenHasher:     tokenHasher,
		idGenerator:     idGenerator,
		clock:           clock,
		magicLinkPolicy: magicLinkPolicy,
	}, nil
}

func (uc *RequestMagicLink) Execute(ctx context.Context, rawEmail string) error {
	return uc.unitOfwork.WithinTransaction(ctx, func(repositories application.IdentityAccessUnitOfWorkRepositories) error {
		return uc.execute(ctx, rawEmail, repositories)
	})
}

func (uc *RequestMagicLink) execute(ctx context.Context, rawEmail string, repositories application.IdentityAccessUnitOfWorkRepositories) error {
	userID := uc.idGenerator.Generate()
	user, err := domain.NewUser(userID, rawEmail)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	existingUser, err := repositories.Users.FindUserByEmail(ctx, user.Email())
	if err == nil {
		user = existingUser
	} else if !errors.Is(err, application.ErrUserNotFound) {
		return fmt.Errorf("finding user: %w", err)
	} else if err := repositories.Users.CreateUser(ctx, user); err != nil {
		if !errors.Is(err, application.ErrUserAlreadyExists) {
			return fmt.Errorf("creating user: %w", err)
		}
		user, err = repositories.Users.FindUserByEmail(ctx, user.Email())
		if err != nil {
			return fmt.Errorf("finding user: %w", err)
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
		id,
		loginTokenHash,
		user.ID(),
		now.Add(uc.magicLinkPolicy.TTL),
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
	}{Email: user.Email().String(), Token: loginTokenRaw})
	if err != nil {
		return fmt.Errorf("encoding event payload: %w", err)
	}

	event, err := outboxDomain.NewEvent(eventID, eventType.String(), payload)

	if err != nil {
		return fmt.Errorf("creating event: %w", err)
	}

	err = repositories.Outbox.Publish(ctx, event)
	if err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}
