package identityaccess_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	identityaccessPg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupUserRepository(ctx context.Context, t *testing.T) (*identityaccessPg.UserRepository, dbexec.Executor, *pgxpool.Pool, *domain.User) {
	dbExec, pool, err := setupPostgresDBExecutor(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	repo := identityaccessPg.NewUserRepository(dbExec, testEmailProtector{})

	if err := dbExec.Exec(ctx, "DELETE FROM users WHERE id = $1", "test"); err != nil {
		t.Fatal(err)
	}

	persistedUser := createTestUser(ctx, t, dbExec, "test", "test@example.com")
	return repo, dbExec, pool, persistedUser
}

func TestFindUserByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	repo, dbExec, pool, persistedUser := setupUserRepository(ctx, t)
	defer pool.Close()

	defer dbExec.Exec(ctx, "DELETE FROM users WHERE id = $1", "test")

	tests := []struct {
		name          string
		userID        string
		expectedError error
		expectedUser  domain.User
	}{
		{
			name:          "find user by id",
			userID:        "test",
			expectedError: nil,
			expectedUser:  *persistedUser,
		},
		{
			name:          "find user by id not found",
			userID:        "test-not-found",
			expectedError: application.ErrUserNotFound,
			expectedUser:  domain.User{},
		},
		{
			name:          "find user by empty id",
			userID:        "",
			expectedError: domain.ErrInvalidUserID,
			expectedUser:  domain.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := repo.FindUserByID(ctx, test.userID)
			if err != test.expectedError {
				t.Errorf("expected error: %v, got: %v", test.expectedError, err)
			}
			if test.expectedError != nil {
				if user != nil {
					t.Fatalf("expected nil user, got %+v", user)
				}
				return
			}

			if user == nil {
				t.Fatal("expected user, got nil")
			}
			if !reflect.DeepEqual(*user, test.expectedUser) {
				t.Errorf("expected user: %v, got: %v", test.expectedUser, user)
			}
		})
	}
}

func TestFindUserByEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	repo, dbExec, pool, persistedUser := setupUserRepository(ctx, t)
	defer pool.Close()

	defer dbExec.Exec(ctx, "DELETE FROM users WHERE id = $1", persistedUser.ID())

	tests := []struct {
		name          string
		email         domain.Email
		expectedError error
		expectedUser  domain.User
	}{
		{
			name:          "find user by email",
			email:         persistedUser.Email(),
			expectedError: nil,
			expectedUser:  *persistedUser,
		},
		{
			name: "find user by email not found",
			email: func() domain.Email {
				email, err := domain.ParseEmail("not-found@example.com")
				if err != nil {
					t.Fatal(err)
				}
				return email
			}(),
			expectedError: application.ErrUserNotFound,
			expectedUser:  domain.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := repo.FindUserByEmail(ctx, test.email)
			if err != test.expectedError {
				t.Errorf("expected error: %v, got: %v", test.expectedError, err)
			}
			if test.expectedError != nil {
				if user != nil {
					t.Fatalf("expected nil user, got %+v", user)
				}
				return
			}

			if user == nil {
				t.Fatal("expected user, got nil")
			}
			if !reflect.DeepEqual(*user, test.expectedUser) {
				t.Errorf("expected user: %v, got: %v", test.expectedUser, user)
			}
		})
	}
}
