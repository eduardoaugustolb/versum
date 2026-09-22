package postgres

const (
	FindUserByIDQuery               = `SELECT id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version FROM users WHERE id = $1`
	FindUserByLookupCandidatesQuery = `
		WITH lookup_candidates AS (
			SELECT lookup_hmac, lookup_key_version, ordinality
			FROM unnest($1::bytea[], $2::smallint[]) WITH ORDINALITY
				AS candidates(lookup_hmac, lookup_key_version, ordinality)
		)
		SELECT u.id, u.email_ciphertext, u.email_lookup_hmac,
			u.email_encryption_key_version, u.email_lookup_key_version
		FROM users AS u
		JOIN lookup_candidates AS c
			ON u.email_lookup_hmac = c.lookup_hmac
			AND u.email_lookup_key_version = c.lookup_key_version
		ORDER BY c.ordinality
		LIMIT 1
	`
	UpdateUserLookupKeyQuery          = `UPDATE users SET email_lookup_hmac = $2, email_lookup_key_version = $3 WHERE id = $1 AND email_lookup_key_version = $4`
	CreateUserQuery                   = `INSERT INTO users (id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING id`
	CreateSessionQuery                = `INSERT INTO sessions (id, secret_hash, user_id, revoked_at, used_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`
	FindSessionByIDQuery              = `SELECT id, secret_hash, user_id, revoked_at, used_at, expires_at FROM sessions WHERE id = $1`
	FindSessionBySecretHashQuery      = `SELECT id, secret_hash, user_id, revoked_at, used_at, expires_at FROM sessions WHERE secret_hash = $1`
	ListSessionsByUserIDQuery         = `SELECT id, secret_hash, user_id, revoked_at, used_at, expires_at FROM sessions WHERE user_id = $1 ORDER BY created_at ASC, id ASC`
	RevokeSessionQuery                = `UPDATE sessions SET revoked_at = $2 WHERE id = $1`
	RevokeAllSessionsQuery            = `UPDATE sessions SET revoked_at = $2 WHERE user_id = $1 AND revoked_at IS NULL`
	CreateLoginTokenQuery             = `INSERT INTO login_tokens (id, token_hash, user_id, expires_at, consumed_at) VALUES ($1, $2, $3, $4, $5)`
	FindLoginTokenByIDQuery           = `SELECT id, token_hash, user_id, expires_at, consumed_at FROM login_tokens WHERE id = $1`
	FindLoginTokenByTokenHashQuery    = `SELECT id, token_hash, user_id, expires_at, consumed_at FROM login_tokens WHERE token_hash = $1`
	ConsumeLoginTokenByTokenHashQuery = `UPDATE login_tokens SET consumed_at = $2 WHERE token_hash = $1`
)
