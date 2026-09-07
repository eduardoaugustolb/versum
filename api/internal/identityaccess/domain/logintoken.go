package domain

import "time"

type LoginToken struct {
	id         string
	tokenHash  []byte
	userId     string
	expiresAt  *time.Time
	consumedAt *time.Time
}

func NewLoginToken(id string, tokenHash []byte, userId string, expiresAt time.Time, now time.Time) (*LoginToken, error) {
	if !expiresAt.After(now) {
		return nil, ErrInvalidLoginTokenExpiresAt
	}
	logintoken := &LoginToken{
		id:        id,
		tokenHash: cloneBytes(tokenHash),
		userId:    userId,
		expiresAt: cloneTime(&expiresAt)}

	if err := validateLoginToken(logintoken); err != nil {
		return nil, err
	}

	return logintoken, nil
}

func RehydrateLoginToken(id string, tokenHash []byte, userId string, expiresAt time.Time, consumedAt *time.Time) (*LoginToken, error) {
	logintoken := &LoginToken{
		id:         id,
		tokenHash:  cloneBytes(tokenHash),
		userId:     userId,
		expiresAt:  cloneTime(&expiresAt),
		consumedAt: cloneTime(consumedAt),
	}

	if err := validateLoginToken(logintoken); err != nil {
		return nil, err
	}

	return logintoken, nil
}

func validateLoginToken(t *LoginToken) error {
	if t.id == "" {
		return ErrInvalidLoginTokenID
	}

	if len(t.tokenHash) == 0 {
		return ErrInvalidLoginTokenHash
	}

	if t.userId == "" {
		return ErrInvalidLoginTokenUserID
	}

	if t.expiresAt == nil || t.expiresAt.IsZero() {
		return ErrInvalidLoginTokenExpiresAt
	}

	if t.consumedAt != nil && t.consumedAt.IsZero() {
		return ErrInvalidLoginTokenConsumedAt
	}

	return nil
}

func (t *LoginToken) ID() string {
	return t.id
}

func (t *LoginToken) TokenHash() []byte {
	return cloneBytes(t.tokenHash)
}

func (t *LoginToken) UserID() string {
	return t.userId
}

func (t *LoginToken) ExpiresAt() time.Time {
	return *t.expiresAt
}

func (t *LoginToken) ConsumedAt() (time.Time, bool) {
	if t.consumedAt == nil {
		return time.Time{}, false
	}
	return *t.consumedAt, true
}

func (t *LoginToken) Consume(now time.Time) error {
	if t.IsConsumed() {
		return ErrLoginTokenAlreadyConsumed
	}
	if t.IsExpired(now) {
		return ErrLoginTokenExpired
	}
	consumedAt := now
	t.consumedAt = &consumedAt
	return nil

}

func (t *LoginToken) IsConsumed() bool {
	return t.consumedAt != nil
}

func (t *LoginToken) IsExpired(now time.Time) bool {
	return t.expiresAt == nil || !t.expiresAt.After(now)
}
