package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityAccessPorts "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/clock"
	"github.com/eduardoaugustolb/versum/api/internal/ports/id"
)

type RequestMagicLink struct {
	userRepo       identityAccessPorts.UserRepository
	loginTokenRepo identityAccessPorts.LoginTokenRepository
	tokenGenerator identityAccessPorts.LoginTokenGenerator
	tokenHasher    identityAccessPorts.LoginTokenHasher
	idGenerator    id.IDGenerator
	clock          clock.Clock
	magicLinkTTL   time.Duration
}

func NewRequestMagicLink(
	userRepo identityAccessPorts.UserRepository,
	loginTokenRepo identityAccessPorts.LoginTokenRepository,
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
		userRepo:       userRepo,
		loginTokenRepo: loginTokenRepo,
		tokenGenerator: tokenGenerator,
		tokenHasher:    tokenHasher,
		idGenerator:    idGenerator,
		clock:          clock,
		magicLinkTTL:   magicLinkTTL,
	}, nil
}

func (uc *RequestMagicLink) Execute(ctx context.Context, email domain.Email) error {
	user, err := uc.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("finding user by email: %w", err)
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

	if err := uc.loginTokenRepo.CreateLoginToken(ctx, loginToken); err != nil {
		return fmt.Errorf("saving login token: %w", err)
	}

	return nil
}
