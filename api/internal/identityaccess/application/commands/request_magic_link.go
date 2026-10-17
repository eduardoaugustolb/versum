package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/events"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxDomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type RequestMagicLink struct {
	limiter         application.EmailRequestLimiter
	unitOfwork      application.IdentityAccessUnitOfWork
	tokenGenerator  application.LoginTokenGenerator
	idGenerator     id.Generator
	clock           clock.Clock
	magicLinkPolicy policy.MagicLinkPolicy
}

func NewRequestMagicLink(
	unitOfWork application.IdentityAccessUnitOfWork,
	tokenGenerator application.LoginTokenGenerator,
	idGenerator id.Generator,
	clock clock.Clock,
	magicLinkPolicy policy.MagicLinkPolicy,
	limiter application.EmailRequestLimiter,
) (*RequestMagicLink, error) {
	if err := policy.ValidateMagicLinkPolicy(magicLinkPolicy); err != nil {
		return nil, err
	}
	if limiter == nil {
		return nil, application.ErrEmailRateLimitUnavailable
	}
	return &RequestMagicLink{
		limiter:         limiter,
		unitOfwork:      unitOfWork,
		tokenGenerator:  tokenGenerator,
		idGenerator:     idGenerator,
		clock:           clock,
		magicLinkPolicy: magicLinkPolicy,
	}, nil
}

func (uc *RequestMagicLink) Execute(ctx context.Context, email domain.Email) error {
	if _, err := domain.ParseEmail(email.String()); err != nil {
		return err
	}
	allowed, err := uc.limiter.Allow(ctx, email)
	if err != nil {
		return fmt.Errorf("%w: %w", application.ErrEmailRateLimitUnavailable, err)
	}
	if !allowed {
		return application.ErrEmailRateLimitExceeded
	}
	return uc.unitOfwork.WithinTransaction(ctx, func(repositories application.IdentityAccessUnitOfWorkRepositories) error {
		return uc.execute(ctx, email, repositories)
	})
}

func (uc *RequestMagicLink) execute(ctx context.Context, email domain.Email, repositories application.IdentityAccessUnitOfWorkRepositories) error {
	userID := uc.idGenerator.Generate()
	user, err := domain.NewUser(userID, email.String())
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

	id := uc.idGenerator.Generate()
	now := uc.clock.Now().UTC()

	loginToken, err := domain.NewLoginToken(
		id,
		user.ID(),
		now.Add(uc.magicLinkPolicy.TTL),
		now,
	)

	if err != nil {
		return fmt.Errorf("creating login token: %w", err)
	}

	if err := repositories.LoginTokens.Create(ctx, loginToken, loginTokenRaw); err != nil {
		return fmt.Errorf("saving login token: %w", err)
	}

	eventID := uc.idGenerator.Generate()
	eventType := events.MagicLinkRequested

	payload, err := json.Marshal(struct {
		LoginTokenID string `json:"login_token_id"`
	}{LoginTokenID: loginToken.ID()})
	if err != nil {
		return fmt.Errorf("encoding event payload: %w", err)
	}

	event, err := outboxDomain.NewEvent(eventID, eventType, user.ID(), payload)

	if err != nil {
		return fmt.Errorf("creating event: %w", err)
	}

	err = repositories.Outbox.Publish(ctx, event)
	if err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}
