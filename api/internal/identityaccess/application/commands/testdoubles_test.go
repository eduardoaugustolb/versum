package commands_test

import (
	"context"
	"sort"
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
	ByID    map[string]*domain.LoginToken
	ByToken map[string]*domain.LoginToken

	CreateErr      error
	FindByIDErr    error
	FindByTokenErr error
	SaveErr        error

	OnCreate func(context.Context, *domain.LoginToken, string) error

	CreatedSecret string
	Created       *domain.LoginToken
	CreateCalls   int
}

func (f *FakeLoginTokenRepository) Create(ctx context.Context, token *domain.LoginToken, secret string) error {
	f.CreateCalls++
	if f.OnCreate != nil {
		return f.OnCreate(ctx, token, secret)
	}
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = token
	f.CreatedSecret = secret
	if f.ByID == nil {
		f.ByID = map[string]*domain.LoginToken{}
	}
	f.ByID[token.ID()] = token
	if f.ByToken == nil {
		f.ByToken = map[string]*domain.LoginToken{}
	}
	f.ByToken[secret] = token
	return nil
}

func (f *FakeLoginTokenRepository) FindByID(_ context.Context, id string) (*domain.LoginToken, error) {
	if f.FindByIDErr != nil {
		return nil, f.FindByIDErr
	}
	if token, ok := f.ByID[id]; ok {
		return token, nil
	}
	return nil, identityports.ErrLoginTokenNotFound
}

func (f *FakeLoginTokenRepository) FindByToken(_ context.Context, secret string) (*domain.LoginToken, error) {
	if f.FindByTokenErr != nil {
		return nil, f.FindByTokenErr
	}
	if token, ok := f.ByToken[secret]; ok {
		return token, nil
	}
	return nil, identityports.ErrLoginTokenNotFound
}

func (f *FakeLoginTokenRepository) Save(_ context.Context, _ *domain.LoginToken) error {
	return f.SaveErr
}

// FakeSessionRepository é o mock universal de identityports.SessionRepository.
type FakeSessionFamilyRevocation struct {
	UserID    string
	FamilyID  string
	RevokedAt time.Time
}

type FakeSessionRepository struct {
	OnFindBySecret    func(context.Context, string) (*domain.Session, error)
	ListErr           error
	RevokeFamilyCalls []FakeSessionFamilyRevocation
	FindBySecretErr   error
	ByID              map[string]*domain.Session
	BySecret          map[string]*domain.Session
	CreatedSecrets    []string
	CreateErr         error
	FindByIDErr       error
	RevokeErr         error

	Created     []*domain.Session
	CreateCalls int
	RevokeCalls []string
}

func (f *FakeSessionRepository) CreateSession(_ context.Context, session *domain.Session, secret string) error {
	f.CreateCalls++
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = append(f.Created, session)
	f.CreatedSecrets = append(f.CreatedSecrets, secret)
	if f.ByID == nil {
		f.ByID = map[string]*domain.Session{}
	}
	f.ByID[session.ID()] = session
	if f.BySecret == nil {
		f.BySecret = map[string]*domain.Session{}
	}
	f.BySecret[secret] = session
	return nil
}

func (f *FakeSessionRepository) FindSessionByID(_ context.Context, id string) (*domain.Session, error) {
	if f.FindByIDErr != nil {
		return nil, f.FindByIDErr
	}
	if session, ok := f.ByID[id]; ok {
		return session, nil
	}
	return nil, identityports.ErrSessionNotFound
}

func (f *FakeSessionRepository) FindSessionBySecret(ctx context.Context, secret string) (*domain.Session, error) {
	if f.OnFindBySecret != nil {
		return f.OnFindBySecret(ctx, secret)
	}
	if f.FindBySecretErr != nil {
		return nil, f.FindBySecretErr
	}
	if session, ok := f.BySecret[secret]; ok {
		return session, nil
	}
	return nil, identityports.ErrSessionNotFound
}

func (f *FakeSessionRepository) ListSessionsByUserID(_ context.Context, userID string) ([]domain.Session, error) {
	if f.ListErr != nil {
		return nil, f.ListErr
	}
	sessions := []domain.Session{}
	for _, session := range f.ByID {
		if session.UserID() == userID {
			sessions = append(sessions, *session)
		}
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].ID() < sessions[j].ID() })
	return sessions, nil
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

func (r *FakeSessionRepository) LockSessionsByUserID(context.Context, string) error  { return nil }
func (r *FakeSessionRepository) SaveRotation(context.Context, *domain.Session) error { return nil }
func (r *FakeSessionRepository) RevokeSessionFamily(_ context.Context, userID, familyID string, revokedAt time.Time) error {
	r.RevokeFamilyCalls = append(r.RevokeFamilyCalls, FakeSessionFamilyRevocation{UserID: userID, FamilyID: familyID, RevokedAt: revokedAt})
	if r.RevokeErr != nil {
		return r.RevokeErr
	}
	if userID == "" {
		return domain.ErrInvalidSessionUserID
	}
	if familyID == "" {
		return domain.ErrInvalidSessionFamilyID
	}
	if revokedAt.IsZero() {
		return domain.ErrInvalidSessionRevokedAt
	}
	for _, session := range r.ByID {
		if session.UserID() == userID && session.FamilyID() == familyID {
			if err := session.Revoke(revokedAt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *FakeSessionRepository) SaveUsage(context.Context, *domain.Session) error { return nil }
