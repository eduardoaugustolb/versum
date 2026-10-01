package postgres

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/jackc/pgx/v5"
)

func TestLoginTokenRepositoryHashesSecretForCreateAndLookup(t *testing.T) {
	now := time.Now().UTC()
	token, err := domain.NewLoginToken("token-1", "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	db := &lookupExecutor{}
	repo := NewLoginTokenRepository(db)
	if err := repo.Create(t.Context(), token, "secret"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindByToken(t.Context(), "secret"); !errors.Is(err, application.ErrLoginTokenNotFound) {
		t.Fatal(err)
	}
	want := sha256.Sum256([]byte("secret"))
	if !bytes.Equal(db.execArgs[0][1].([]byte), want[:]) || !bytes.Equal(db.queryArgs[0][0].([]byte), want[:]) {
		t.Fatal("creation and lookup must use the SHA-256 digest")
	}
}

func TestLoginTokenRepositoryRejectsEmptySecret(t *testing.T) {
	repo := NewLoginTokenRepository(&lookupExecutor{})
	if err := repo.Create(t.Context(), nil, ""); !errors.Is(err, application.ErrInvalidLoginTokenSecret) {
		t.Fatal(err)
	}
	if _, err := repo.FindByToken(t.Context(), ""); !errors.Is(err, application.ErrInvalidLoginTokenSecret) {
		t.Fatal(err)
	}
}

func TestLoginTokenRepositoryRejectsDuplicateConsumption(t *testing.T) {
	now := time.Now().UTC()
	token, err := domain.NewLoginToken("token-1", "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := token.Consume(now); err != nil {
		t.Fatal(err)
	}
	repo := NewLoginTokenRepository(userRepositoryExecutor{row: errorRow{err: pgx.ErrNoRows}})
	if err := repo.Save(t.Context(), token); !errors.Is(err, domain.ErrLoginTokenAlreadyConsumed) {
		t.Fatal(err)
	}
}
