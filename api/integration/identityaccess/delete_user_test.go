package identityaccess_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	identitypg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

type observedDeletionWork struct {
	application.IdentityAccessUnitOfWork
	locking chan struct{}
}

func (u observedDeletionWork) WithinTransaction(ctx context.Context, fn func(application.IdentityAccessUnitOfWorkRepositories) error) error {
	return u.IdentityAccessUnitOfWork.WithinTransaction(ctx, func(repos application.IdentityAccessUnitOfWorkRepositories) error {
		repos.Sessions = observedDeletionSessions{repos.Sessions, u.locking}
		return fn(repos)
	})
}

type observedDeletionSessions struct {
	application.SessionRepository
	locking chan struct{}
}

func (r observedDeletionSessions) LockSessionsByUserID(ctx context.Context, user string) error {
	close(r.locking)
	return r.SessionRepository.LockSessionsByUserID(ctx, user)
}
func TestDeletionRejectsSessionRevokedWhileWaitingForLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	repo, db, pool, user := setupSessionRepository(ctx, t)
	now := time.Now().UTC()
	session, err := domain.NewSession("delete-session", user.ID(), "delete-family", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateSession(ctx, session, "delete-secret"); err != nil {
		t.Fatal(err)
	}
	blocker, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer blocker.Rollback(ctx)
	if _, err := blocker.Exec(ctx, "SELECT id FROM users WHERE id=$1 FOR UPDATE", user.ID()); err != nil {
		t.Fatal(err)
	}
	if _, err := blocker.Exec(ctx, "UPDATE sessions SET revoked_at=$2 WHERE id=$1", session.ID(), now); err != nil {
		t.Fatal(err)
	}
	locking := make(chan struct{})
	uow := observedDeletionWork{identitypg.NewUnitOfWork(db, testEmailProtector{}, nil), locking}
	result := make(chan error, 1)
	go func() { result <- commands.NewDeleteUser(uow, clock.SystemClock{}).Execute(ctx, "delete-secret") }()
	select {
	case <-locking:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if err := blocker.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-result:
		if !errors.Is(err, application.ErrInvalidSessionSecret) {
			t.Fatalf("expected rejected deletion, got %v", err)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if _, err := identitypg.NewUserRepository(db, testEmailProtector{}).FindUserByID(ctx, user.ID()); err != nil {
		t.Fatalf("user removed: %v", err)
	}
}
