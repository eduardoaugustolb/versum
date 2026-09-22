package identityaccess_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	identityaccessPg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func setupUserRepository(ctx context.Context, t *testing.T) (*identityaccessPg.UserRepository, *domain.User) {
	t.Helper()

	dbExec, pool, err := setupPostgresDBExecutor(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	repo := identityaccessPg.NewUserRepository(dbExec, testEmailProtector{})

	if err := dbExec.Exec(ctx, "DELETE FROM users WHERE id = $1", "test"); err != nil {
		t.Fatal(err)
	}

	persistedUser := createTestUser(ctx, t, dbExec, "test", "test@example.com")

	t.Cleanup(func() { // executa primeiro, pois Cleanup é LIFO
		if err := dbExec.Exec(
			context.Background(),
			"DELETE FROM users WHERE id = $1",
			"test",
		); err != nil {
			t.Error(err)
		}
	})

	return repo, persistedUser
}

func TestFindUserByID(t *testing.T) {
	ctx := context.Background()
	repo, persistedUser := setupUserRepository(ctx, t)

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
	ctx := context.Background()
	repo, persistedUser := setupUserRepository(ctx, t)

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

func TestFindUserByEmailMigratesLookupKeyVersion(t *testing.T) {
	ctx := context.Background()
	dbExec, pool, err := setupPostgresDBExecutor(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	const userID = "lookup-rotation-user"
	const rawEmail = "lookup-rotation@example.com"
	if err := dbExec.Exec(ctx, "DELETE FROM users WHERE id = $1", userID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := dbExec.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID); err != nil {
			t.Error(err)
		}
	})

	if err := dbExec.Exec(
		ctx,
		`INSERT INTO users (id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version)
		 VALUES ($1, $2, $3, $4, $5)`,
		userID,
		[]byte("ciphertext:"+rawEmail),
		[]byte("lookup:v1:"+rawEmail),
		1,
		1,
	); err != nil {
		t.Fatal(err)
	}

	repo := identityaccessPg.NewUserRepository(dbExec, rotatingEmailProtector{})
	email, err := domain.ParseEmail(rawEmail)
	if err != nil {
		t.Fatal(err)
	}
	user, err := repo.FindUserByEmail(ctx, email)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID() != userID {
		t.Fatalf("expected user %q, got %q", userID, user.ID())
	}

	var lookupHMAC []byte
	var lookupKeyVersion int
	if err := dbExec.QueryRow(ctx, "SELECT email_lookup_hmac, email_lookup_key_version FROM users WHERE id = $1", userID).Scan(&lookupHMAC, &lookupKeyVersion); err != nil {
		t.Fatal(err)
	}
	if string(lookupHMAC) != "lookup:v2:"+rawEmail || lookupKeyVersion != 2 {
		t.Fatalf("expected lookup migration to version 2, got hmac=%q version=%d", lookupHMAC, lookupKeyVersion)
	}
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	repo, persistedUser := setupUserRepository(ctx, t)
	anotherEmail, err := domain.ParseEmail("anotheruser@test.com")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name          string
		id            string
		email         domain.Email
		expectedError error
	}{
		{
			name:          "try create existed user",
			email:         persistedUser.Email(),
			id:            persistedUser.ID(),
			expectedError: application.ErrUserAlreadyExists,
		},
		{
			name:  "create user",
			email: anotherEmail,
			id:    "another-id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := domain.NewUser(tc.id, tc.email.String())
			if err != nil {
				t.Fatal(err)
			}
			err = repo.CreateUser(ctx, user)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error matching: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
