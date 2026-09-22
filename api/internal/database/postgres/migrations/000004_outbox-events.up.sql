CREATE TABLE outbox_events (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload_ciphertext BYTEA NOT NULL,
    payload_key_version SMALLINT NOT NULL
        CHECK (payload_key_version > 0),

    attempts INTEGER NOT NULL DEFAULT 0
        CHECK (attempts >= 0),
    available_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    lease_token TEXT,
    leased_until TIMESTAMPTZ,

    processed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    last_error_redacted TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        (lease_token IS NULL AND leased_until IS NULL)
        OR
        (lease_token IS NOT NULL AND leased_until IS NOT NULL)
    ),
    CHECK (NOT (processed_at IS NOT NULL AND failed_at IS NOT NULL))
  );

CREATE INDEX outbox_events_ready_idx
    ON outbox_events (available_at, leased_until)
    WHERE processed_at IS NULL AND failed_at IS NULL;

CREATE INDEX outbox_events_type_created_idx
    ON outbox_events (event_type, created_at);
