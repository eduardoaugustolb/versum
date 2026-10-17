CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email_ciphertext BYTEA NOT NULL,
    email_lookup_hmac BYTEA NOT NULL UNIQUE,
    email_encryption_key_version SMALLINT NOT NULL CHECK (email_encryption_key_version > 0),
    email_lookup_key_version SMALLINT NOT NULL CHECK (email_lookup_key_version > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE login_tokens (
    id TEXT PRIMARY KEY,
    token_hash BYTEA NOT NULL UNIQUE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > created_at),
    consumed_at TIMESTAMPTZ
);

CREATE INDEX login_tokens_user_created_idx ON login_tokens (user_id, created_at);
CREATE INDEX login_tokens_expires_idx ON login_tokens (expires_at);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    secret_hash BYTEA NOT NULL UNIQUE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > created_at),
    last_used_at TIMESTAMPTZ,
    family_id TEXT NOT NULL CHECK (length(trim(family_id)) > 0),
    replaced_by_session_id TEXT REFERENCES sessions(id) DEFERRABLE INITIALLY DEFERRED,
    ip_address TEXT NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT '',
    CONSTRAINT sessions_replacement_state CHECK (replaced_by_session_id IS NULL OR (replaced_by_session_id <> id AND revoked_at IS NOT NULL)),
    CONSTRAINT sessions_user_agent_length CHECK (octet_length(user_agent) <= 1024)
);

CREATE INDEX sessions_user_created_idx ON sessions (user_id, created_at);
CREATE INDEX sessions_last_used_idx ON sessions (last_used_at);
CREATE INDEX sessions_expires_idx ON sessions (expires_at);

CREATE INDEX sessions_family_id_idx ON sessions (family_id);
