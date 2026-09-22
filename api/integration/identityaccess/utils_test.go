package identityaccess_test

import (
	"context"
	"os"
	"strings"
	"testing"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/database/postgres"
	identityaccesspg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
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

func (testEmailProtector) LookupCandidates(_ context.Context, email domain.Email) ([]ports.LookupCandidate, error) {
	return []ports.LookupCandidate{{LookupKeyVersion: 1, LookupHMAC: []byte("lookup:" + email.String())}}, nil
}

type rotatingEmailProtector struct{}

func (rotatingEmailProtector) Protect(_ context.Context, email domain.Email) (ports.ProtectedEmail, error) {
	value := email.String()
	return ports.ProtectedEmail{
		Ciphertext:           []byte("ciphertext:" + value),
		LookupHMAC:           []byte("lookup:v2:" + value),
		EncryptionKeyVersion: 1,
		LookupKeyVersion:     2,
	}, nil
}

func (rotatingEmailProtector) Unprotect(_ context.Context, protected ports.ProtectedEmail) (domain.Email, error) {
	return domain.ParseEmail(strings.TrimPrefix(string(protected.Ciphertext), "ciphertext:"))
}

func (rotatingEmailProtector) LookupHMAC(_ context.Context, email domain.Email) ([]byte, error) {
	return []byte("lookup:v2:" + email.String()), nil
}

func (rotatingEmailProtector) LookupCandidates(_ context.Context, email domain.Email) ([]ports.LookupCandidate, error) {
	value := email.String()
	return []ports.LookupCandidate{
		{LookupKeyVersion: 2, LookupHMAC: []byte("lookup:v2:" + value)},
		{LookupKeyVersion: 1, LookupHMAC: []byte("lookup:v1:" + value)},
	}, nil
}

func createTestUser(ctx context.Context, t *testing.T, db dbexec.Executor, id, rawEmail string) *domain.User {
	t.Helper()
	user, err := domain.NewUser(id, rawEmail)
	if err != nil {
		t.Fatal(err)
	}
	if err := identityaccesspg.NewUserRepository(db, testEmailProtector{}).CreateUser(ctx, user); err != nil {
		t.Fatal(err)
	}
	return user
}
