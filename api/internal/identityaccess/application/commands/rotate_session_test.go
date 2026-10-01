package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type rotationSessions struct {
	FakeSessionRepository
	lockErr, saveErr, revokeFamilyErr      error
	finds, locks, saves, familyRevocations int
	user, family                           string
	rereadErr                              error
	afterLock                              *domain.Session
}

func (r *rotationSessions) FindSessionBySecret(ctx context.Context, secret string) (*domain.Session, error) {
	r.finds++
	if r.finds == 2 {
		if r.rereadErr != nil {
			return nil, r.rereadErr
		}
		if r.afterLock != nil {
			return r.afterLock, nil
		}
	}
	return r.FakeSessionRepository.FindSessionBySecret(ctx, secret)
}
func (r *rotationSessions) LockSessionsByUserID(_ context.Context, user string) error {
	r.locks++
	r.user = user
	return r.lockErr
}
func (r *rotationSessions) SaveRotation(_ context.Context, _ *domain.Session) error {
	r.saves++
	return r.saveErr
}
func (r *rotationSessions) RevokeSessionFamily(_ context.Context, user, family string, _ time.Time) error {
	r.familyRevocations++
	r.user = user
	r.family = family
	return r.revokeFamilyErr
}

type sessionSecrets struct {
	secret string
	err    error
}

func (g sessionSecrets) GenerateSessionSecret() (string, error) { return g.secret, g.err }

type rotationUow struct {
	repos               application.IdentityAccessUnitOfWorkRepositories
	beginErr, commitErr error
	committed           bool
}

func (u *rotationUow) WithinTransaction(_ context.Context, fn func(application.IdentityAccessUnitOfWorkRepositories) error) error {
	if u.beginErr != nil {
		return u.beginErr
	}
	if err := fn(u.repos); err != nil {
		return err
	}
	if u.commitErr != nil {
		return u.commitErr
	}
	u.committed = true
	return nil
}

func TestRotateSessionSuccessAndReuseCommitsFamilyRevocation(t *testing.T) {
	now := time.Now().UTC()
	old, _ := domain.NewSession("old", "user", "family", domain.SessionClient{}, now.Add(time.Hour), now)
	repo := &rotationSessions{FakeSessionRepository: FakeSessionRepository{BySecret: map[string]*domain.Session{"old-secret": old}}}
	uow := &rotationUow{repos: application.IdentityAccessUnitOfWorkRepositories{Sessions: repo}}
	uc := commands.NewRotateSession(uow, sessionSecrets{secret: "new-secret"}, &FakeIDGenerator{IDs: []string{"next"}}, FakeClock{NowTime: now})
	got, err := uc.Execute(t.Context(), "old-secret", domain.SessionClient{IPAddress: "192.0.2.1", UserAgent: "browser"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Secret != "new-secret" || got.Session.FamilyID() != old.FamilyID() || got.Session.UserID() != old.UserID() || !got.Session.ExpiresAt().Equal(old.ExpiresAt()) || got.Session.IPAddress() != "192.0.2.1" {
		t.Fatal("incorrect successor")
	}
	if repo.finds != 2 || repo.locks != 1 || repo.saves != 1 || !uow.committed {
		t.Fatal("rotation must lock, re-read and commit")
	}
	uow.committed = false
	got, err = uc.Execute(t.Context(), "old-secret", domain.SessionClient{})
	if got != nil || !errors.Is(err, domain.ErrSessionReused) {
		t.Fatalf("expected reuse, got %v", err)
	}
	if !uow.committed || repo.familyRevocations != 1 || repo.family != "family" || repo.user != "user" {
		t.Fatal("reuse revocation must commit before returning the error")
	}
	if repo.CreateCalls != 1 {
		t.Fatal("reuse must not issue another credential")
	}
}

func TestRotateSessionFailurePaths(t *testing.T) {
	failure := errors.New("dependency failure")
	now := time.Now().UTC()
	for _, name := range []string{"not found", "lookup", "lock", "reread", "expired", "revoked", "generation", "empty secret", "same secret", "invalid client", "same id", "create", "save", "begin", "commit", "family revoke"} {
		t.Run(name, func(t *testing.T) {
			old, _ := domain.NewSession("old", "user", "family", domain.SessionClient{}, now.Add(time.Hour), now)
			repo := &rotationSessions{FakeSessionRepository: FakeSessionRepository{BySecret: map[string]*domain.Session{"old-secret": old}}}
			uow := &rotationUow{repos: application.IdentityAccessUnitOfWorkRepositories{Sessions: repo}}
			generator := sessionSecrets{secret: "new-secret"}
			ids := &FakeIDGenerator{IDs: []string{"next"}}
			client := domain.SessionClient{}
			at := now
			want := failure
			switch name {
			case "not found":
				repo.BySecret = nil
				want = application.ErrSessionNotFound
			case "lookup":
				repo.FindBySecretErr = failure
			case "lock":
				repo.lockErr = failure
			case "reread":
				repo.rereadErr = failure
			case "expired":
				at = old.ExpiresAt()
				want = domain.ErrInvalidSessionState
			case "revoked":
				old.Revoke(now)
				want = domain.ErrInvalidSessionState
			case "generation":
				generator.err = failure
			case "empty secret":
				generator.secret = ""
				want = application.ErrInvalidSessionSecret
			case "same secret":
				generator.secret = "old-secret"
				want = application.ErrInvalidSessionSecret
			case "invalid client":
				client.IPAddress = "invalid"
				want = domain.ErrInvalidSessionIPAddress
			case "same id":
				ids.IDs = []string{"old"}
				want = domain.ErrInvalidSessionReplacement
			case "create":
				repo.CreateErr = failure
			case "save":
				repo.saveErr = failure
			case "begin":
				uow.beginErr = failure
			case "commit":
				uow.commitErr = failure
			case "family revoke":
				next, _ := domain.NewSession("next", "user", "family", domain.SessionClient{}, old.ExpiresAt(), now)
				old.ReplaceWith(next, now)
				repo.revokeFamilyErr = failure
			}
			uc := commands.NewRotateSession(uow, generator, ids, FakeClock{NowTime: at})
			got, err := uc.Execute(t.Context(), "old-secret", client)
			if got != nil || !errors.Is(err, want) {
				t.Fatalf("want %v and no credential, got %+v, %v", want, got, err)
			}
			if uow.committed {
				t.Fatal("failed transaction must not commit")
			}
		})
	}
}

func TestRotateSessionRechecksStateAfterLock(t *testing.T) {
	now := time.Now().UTC()
	before, _ := domain.NewSession("old", "user", "family", domain.SessionClient{}, now.Add(time.Hour), now)
	after, _ := domain.RehydrateSession("old", "user", "family", domain.SessionClient{}, &now, &now, before.ExpiresAt(), "other-successor")
	repo := &rotationSessions{FakeSessionRepository: FakeSessionRepository{BySecret: map[string]*domain.Session{"secret": before}}, afterLock: after}
	uow := &rotationUow{repos: application.IdentityAccessUnitOfWorkRepositories{Sessions: repo}}
	uc := commands.NewRotateSession(uow, sessionSecrets{secret: "new"}, &FakeIDGenerator{}, FakeClock{NowTime: now})
	if _, err := uc.Execute(t.Context(), "secret", domain.SessionClient{}); !errors.Is(err, domain.ErrSessionReused) {
		t.Fatal(err)
	}
	if !uow.committed || repo.CreateCalls != 0 || repo.familyRevocations != 1 {
		t.Fatal("must detect concurrent rotation after acquiring the lock")
	}
}
