package identityaccess_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	identityaccesspg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupPostgresDBExecutor(ctx context.Context, t *testing.T) (dbexec.Executor, *pgxpool.Pool, error) {
	t.Helper()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, err
	}

	dbExec := postgres.NewPgxExecutor(pool)
	return dbExec, pool, nil
}

type testEmailProtector struct{}

func (testEmailProtector) Protect(_ context.Context, email domain.Email) (ports.ProtectedEmail, error) {
	value := email.String()
	return ports.ProtectedEmail{
		Ciphertext:           []byte("ciphertext:" + value),
		LookupHMAC:           []byte("lookup:" + value),
		EncryptionKeyVersion: 1,
		LookupKeyVersion:     1,
	}, nil
}

func (testEmailProtector) Unprotect(_ context.Context, protected ports.ProtectedEmail) (domain.Email, error) {
	return domain.ParseEmail(strings.TrimPrefix(string(protected.Ciphertext), "ciphertext:"))
}

func (testEmailProtector) LookupHMAC(_ context.Context, email domain.Email) ([]byte, error) {
	return []byte("lookup:" + email.String()), nil
}

func createTestUser(ctx context.Context, t *testing.T, db dbexec.Executor, id, rawEmail string) *domain.User {
	t.Helper()
	email, err := domain.ParseEmail(rawEmail)
	if err != nil {
		t.Fatal(err)
	}
	user, err := domain.NewUser(id, email)
	if err != nil {
		t.Fatal(err)
	}
	if err := identityaccesspg.NewUserRepository(db, testEmailProtector{}).CreateUser(ctx, user); err != nil {
		t.Fatal(err)
	}
	return user
}
