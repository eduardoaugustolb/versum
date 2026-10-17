package commands_test

import (
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestDeleteUser_RemovesUserData(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	userRepo := &FakeUserRepository{}
	user, err := domain.NewUser("userID", "email@example.com")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	sessionRepo := &FakeSessionRepository{}
	session, err := domain.NewSession("ID", "userID", "familyID", "127.0.0.1", "userAgent", now.Add(time.Hour), now)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	userRepo.CreateUser(t.Context(), user)
	sessionRepo.CreateSession(t.Context(), session, "secret")

	deleteUser := commands.NewDeleteUser(&FakeUnitOfWork{Inner: application.IdentityAccessUnitOfWorkRepositories{Users: userRepo, Sessions: sessionRepo}}, FakeClock{NowTime: now})
	err = deleteUser.Execute(t.Context(), "secret")
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	if len(userRepo.DeleteCalls) != 1 || userRepo.DeleteCalls[0] != "userID" {
		t.Fatalf("expected delete call with user id, got %v", userRepo.DeleteCalls)
	}

	user, err = userRepo.FindUserByID(t.Context(), "userID")
	if !errors.Is(err, application.ErrUserNotFound) || user != nil {
		t.Fatalf("expected deleted user to be absent, got %v, %v", user, err)
	}

	email, err := domain.ParseEmail("email@example.com")
	if err != nil {
		t.Fatal(err)
	}
	user, err = userRepo.FindUserByEmail(t.Context(), email)
	if !errors.Is(err, application.ErrUserNotFound) || user != nil {
		t.Fatalf("expected deleted email to be absent, got %v, %v", user, err)
	}
}

func TestDeleteUser_ReturnsErrorWhenUserNotFound(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	userRepo := &FakeUserRepository{}

	sessionRepo := &FakeSessionRepository{}
	session, err := domain.NewSession("ID", "userID", "familyID", "127.0.0.1", "userAgent", now.Add(time.Hour), now)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	sessionRepo.CreateSession(t.Context(), session, "secret")

	deleteUser := commands.NewDeleteUser(&FakeUnitOfWork{Inner: application.IdentityAccessUnitOfWorkRepositories{Users: userRepo, Sessions: sessionRepo}}, FakeClock{NowTime: now})
	err = deleteUser.Execute(t.Context(), "secret")
	if err == nil || !errors.Is(err, application.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteUser_ReturnsErrorWhenSessionNotFound(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	userRepo := &FakeUserRepository{}

	sessionRepo := &FakeSessionRepository{}

	deleteUser := commands.NewDeleteUser(&FakeUnitOfWork{Inner: application.IdentityAccessUnitOfWorkRepositories{Users: userRepo, Sessions: sessionRepo}}, FakeClock{NowTime: now})
	err := deleteUser.Execute(t.Context(), "secret")
	if err == nil || !errors.Is(err, application.ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestDeleteUser_ReturnsErrorWhenExpiredSession(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	userRepo := &FakeUserRepository{}

	sessionRepo := &FakeSessionRepository{}
	session, err := domain.NewSession("ID", "userID", "familyID", "127.0.0.1", "userAgent", now.Add(-time.Hour), now.Add(-2*time.Hour))
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	sessionRepo.CreateSession(t.Context(), session, "secret")

	deleteUser := commands.NewDeleteUser(&FakeUnitOfWork{Inner: application.IdentityAccessUnitOfWorkRepositories{Users: userRepo, Sessions: sessionRepo}}, FakeClock{NowTime: now})
	err = deleteUser.Execute(t.Context(), "secret")
	if err == nil || !errors.Is(err, application.ErrSessionExpired) {
		t.Fatalf("expected ErrSessionExpired, got %v", err)
	}
}

func TestDeleteUser_ReturnsErrorWhenRevoked(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	userRepo := &FakeUserRepository{}

	sessionRepo := &FakeSessionRepository{}
	session, err := domain.NewSession("ID", "userID", "familyID", "127.0.0.1", "userAgent", now.Add(time.Hour), now)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	session.Revoke(now)
	sessionRepo.CreateSession(t.Context(), session, "secret")

	deleteUser := commands.NewDeleteUser(&FakeUnitOfWork{Inner: application.IdentityAccessUnitOfWorkRepositories{Users: userRepo, Sessions: sessionRepo}}, FakeClock{NowTime: now})
	err = deleteUser.Execute(t.Context(), "secret")
	if err == nil || !errors.Is(err, application.ErrInvalidSessionSecret) {
		t.Fatalf("expected ErrInvalidSessionSecret, got %v", err)
	}
}

func TestDeleteUserRechecksAuthorizationAfterLock(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	active, _ := domain.NewSession("session", "user", "family", "", "", now.Add(time.Hour), now)
	revoked, _ := domain.RehydrateSession("session", "user", "family", "", "", &now, nil, now.Add(time.Hour), "")
	for _, test := range []struct {
		name               string
		lockErr, rereadErr error
		after              *domain.Session
		want               error
	}{
		{name: "revoked while waiting", after: revoked, want: application.ErrInvalidSessionSecret},
		{name: "lock failure", lockErr: errors.New("lock failed")},
		{name: "reread failure", rereadErr: application.ErrSessionNotFound, want: application.ErrSessionNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			users := &FakeUserRepository{}
			sessions := &rotationSessions{FakeSessionRepository: FakeSessionRepository{BySecret: map[string]*domain.Session{"secret": active}}, lockErr: test.lockErr, rereadErr: test.rereadErr, afterLock: test.after}
			uow := &rotationUow{repos: application.IdentityAccessUnitOfWorkRepositories{Users: users, Sessions: sessions}}
			err := commands.NewDeleteUser(uow, FakeClock{NowTime: now}).Execute(t.Context(), "secret")
			want := test.want
			if test.lockErr != nil {
				want = test.lockErr
			}
			if !errors.Is(err, want) {
				t.Fatalf("want %v, got %v", want, err)
			}
			if len(users.DeleteCalls) != 0 || uow.committed {
				t.Fatal("unauthorized deletion committed")
			}
		})
	}
}
