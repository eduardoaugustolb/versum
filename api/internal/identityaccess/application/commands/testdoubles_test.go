package commands_test

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	outboxports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	outboxdomain "github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

var (
	_ identityports.UserRepository           = (*FakeUserRepository)(nil)
	_ identityports.LoginTokenRepository     = (*FakeLoginTokenRepository)(nil)
	_ identityports.SessionRepository        = (*FakeSessionRepository)(nil)
	_ outboxports.OutboxRepository           = (*FakeOutboxRepository)(nil)
	_ identityports.IdentityAccessUnitOfWork = (*FakeUnitOfWork)(nil)
	_ clock.Clock                            = (*FakeClock)(nil)
	_ id.Generator                           = (*FakeIDGenerator)(nil)
	_ identityports.LoginTokenGenerator      = (*FakeLoginTokenGenerator)(nil)
	_ identityports.LoginTokenHasher         = (*FakeLoginTokenHasher)(nil)
)

// FakeUserRepository é o mock universal de identityports.UserRepository.
// Configure o estado (ByID/ByEmail), injete erros por método ou use os hooks
// OnFindByEmail/OnCreate para comportamentos sequenciais.
type FakeUserRepository struct {
	ByID    map[string]*domain.User
	ByEmail map[string]*domain.User

	CreateErr      error
	FindByIDErr    error
	FindByEmailErr error

	OnCreate      func(context.Context, *domain.User) error
	OnFindByEmail func(context.Context, domain.Email) (*domain.User, error)

	Created          []*domain.User
	CreateCalls      int
	FindByIDCalls    []string
	FindByEmailCalls []domain.Email
}

func (f *FakeUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	f.CreateCalls++
	if f.OnCreate != nil {
		return f.OnCreate(ctx, user)
	}
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = append(f.Created, user)
	if f.ByID == nil {
		f.ByID = map[string]*domain.User{}
	}
	f.ByID[user.ID()] = user
	if f.ByEmail == nil {
		f.ByEmail = map[string]*domain.User{}
	}
	f.ByEmail[user.Email().String()] = user
	return nil
}

func (f *FakeUserRepository) FindUserByID(_ context.Context, id string) (*domain.User, error) {
	f.FindByIDCalls = append(f.FindByIDCalls, id)
	if f.FindByIDErr != nil {
		return nil, f.FindByIDErr
	}
	if user, ok := f.ByID[id]; ok {
		return user, nil
	}
	return nil, identityports.ErrUserNotFound
}

func (f *FakeUserRepository) FindUserByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	f.FindByEmailCalls = append(f.FindByEmailCalls, email)
	if f.OnFindByEmail != nil {
		return f.OnFindByEmail(ctx, email)
	}
	if f.FindByEmailErr != nil {
		return nil, f.FindByEmailErr
	}
	if user, ok := f.ByEmail[email.String()]; ok {
		return user, nil
	}
	return nil, identityports.ErrUserNotFound
}

// FakeLoginTokenRepository é o mock universal de identityports.LoginTokenRepository.
type FakeLoginTokenRepository struct {
	ByID   map[string]*domain.LoginToken
	ByHash map[string]*domain.LoginToken

	CreateErr     error
	FindByIDErr   error
	FindByHashErr error
	ConsumeErr    error

	OnCreate func(context.Context, *domain.LoginToken) error

	Created     *domain.LoginToken
	CreateCalls int
}

func (f *FakeLoginTokenRepository) CreateLoginToken(ctx context.Context, token *domain.LoginToken) error {
	f.CreateCalls++
	if f.OnCreate != nil {
		return f.OnCreate(ctx, token)
	}
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = token
	return nil
}

func (f *FakeLoginTokenRepository) FindLoginTokenByID(_ context.Context, id string) (*domain.LoginToken, error) {
	if f.FindByIDErr != nil {
		return nil, f.FindByIDErr
	}
	if token, ok := f.ByID[id]; ok {
		return token, nil
	}
	return nil, identityports.ErrLoginTokenNotFound
}

func (f *FakeLoginTokenRepository) FindLoginTokenByTokenHash(_ context.Context, tokenHash []byte) (*domain.LoginToken, error) {
	if f.FindByHashErr != nil {
		return nil, f.FindByHashErr
	}
	if token, ok := f.ByHash[string(tokenHash)]; ok {
		return token, nil
	}
	return nil, identityports.ErrLoginTokenNotFound
}

func (f *FakeLoginTokenRepository) ConsumeLoginTokenByTokenHash(_ context.Context, _ []byte, _ *time.Time) error {
	return f.ConsumeErr
}

// FakeSessionRepository é o mock universal de identityports.SessionRepository.
type FakeSessionRepository struct {
	CreateErr   error
	FindByIDErr error
	RevokeErr   error

	Created     []*domain.Session
	CreateCalls int
	RevokeCalls []string
}

func (f *FakeSessionRepository) CreateSession(_ context.Context, session *domain.Session) error {
	f.CreateCalls++
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = append(f.Created, session)
	return nil
}

func (f *FakeSessionRepository) FindSessionByID(_ context.Context, _ string) (*domain.Session, error) {
	if f.FindByIDErr != nil {
		return nil, f.FindByIDErr
	}
	return nil, identityports.ErrSessionNotFound
}

func (f *FakeSessionRepository) FindSessionBySecretHash(_ context.Context, _ []byte) (*domain.Session, error) {
	return nil, identityports.ErrSessionNotFound
}

func (f *FakeSessionRepository) ListSessionsByUserID(_ context.Context, _ string) ([]domain.Session, error) {
	return nil, nil
}

func (f *FakeSessionRepository) RevokeSession(_ context.Context, sessionID string, _ *time.Time) error {
	f.RevokeCalls = append(f.RevokeCalls, sessionID)
	return f.RevokeErr
}

func (f *FakeSessionRepository) RevokeAllSessions(_ context.Context, _ string, _ *time.Time) error {
	return f.RevokeErr
}

// FakeOutboxRepository é o mock universal de outboxports.OutboxRepository.
type FakeOutboxRepository struct {
	PublishErr error
	OnPublish  func(context.Context, *outboxdomain.Event) error

	Published []*outboxdomain.Event
	Calls     int
}

func (f *FakeOutboxRepository) Publish(ctx context.Context, event *outboxdomain.Event) error {
	f.Calls++
	if f.OnPublish != nil {
		return f.OnPublish(ctx, event)
	}
	if f.PublishErr != nil {
		return f.PublishErr
	}
	f.Published = append(f.Published, event)
	return nil
}

// FakeUnitOfWork é o mock universal de identityports.IdentityAccessUnitOfWork.
// Encaminha para Inner por padrão; injete TxErr ou OnTransaction para falhas.
type FakeUnitOfWork struct {
	Inner         identityports.IdentityAccessUnitOfWorkRepositories
	TxErr         error
	OnTransaction func(identityports.IdentityAccessUnitOfWorkRepositories) error

	Calls int
}

func (f *FakeUnitOfWork) WithinTransaction(_ context.Context, fn func(identityports.IdentityAccessUnitOfWorkRepositories) error) error {
	f.Calls++
	if f.TxErr != nil {
		return f.TxErr
	}
	if f.OnTransaction != nil {
		return fn(f.Inner)
	}
	return fn(f.Inner)
}

// FakeClock é o mock universal de clock.Clock.
type FakeClock struct {
	NowTime time.Time
}

func (c FakeClock) Now() time.Time { return c.NowTime }

// FakeIDGenerator é o mock universal de id.Generator.
// Retorna os IDs da fila em ordem; esgotada a fila, repete o último
// (ou "id-1" quando a fila está vazia).
type FakeIDGenerator struct {
	IDs   []string
	calls int
}

func (f *FakeIDGenerator) Generate() string {
	if len(f.IDs) == 0 {
		return "id-1"
	}
	idx := f.calls
	if idx >= len(f.IDs) {
		idx = len(f.IDs) - 1
	}
	f.calls++
	return f.IDs[idx]
}

// FakeLoginTokenGenerator é o mock universal de identityports.LoginTokenGenerator.
type FakeLoginTokenGenerator struct {
	Token string
	Err   error
}

func (f FakeLoginTokenGenerator) GenerateLoginToken() (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	if f.Token != "" {
		return f.Token, nil
	}
	return "raw-token", nil
}

// FakeLoginTokenHasher é o mock universal de identityports.LoginTokenHasher.
type FakeLoginTokenHasher struct {
	Sum []byte
	Err error
}

func (f FakeLoginTokenHasher) Hash(_ string) ([]byte, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	if f.Sum != nil {
		return f.Sum, nil
	}
	return []byte("hashed-token"), nil
}
