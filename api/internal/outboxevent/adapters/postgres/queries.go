package postgres

const (
	ClaimEventQuery = `
		WITH candidate AS (
			SELECT id FROM outbox_events
			WHERE processed_at IS NULL AND failed_at IS NULL
			  AND available_at <= $1 AND (leased_until IS NULL OR leased_until <= $1)
			ORDER BY available_at, created_at, id
			LIMIT 1 FOR UPDATE SKIP LOCKED
		)
		UPDATE outbox_events AS event
		SET lease_token = $2, leased_until = $3, attempts = attempts + 1
		FROM candidate WHERE event.id = candidate.id
		RETURNING event.id, event.event_type, COALESCE(event.user_id, ''),
		          event.payload_ciphertext, event.payload_key_version, event.attempts
	`
	CompleteEventQuery = `
		UPDATE outbox_events SET processed_at = $3, lease_token = NULL, leased_until = NULL,
		    payload_ciphertext = ''::bytea, last_error_redacted = NULL
		WHERE id = $1 AND lease_token = $2 AND leased_until > $3
		  AND processed_at IS NULL AND failed_at IS NULL RETURNING id
	`
	RetryEventQuery = `
		UPDATE outbox_events SET available_at = $4, lease_token = NULL, leased_until = NULL,
		    last_error_redacted = 'delivery_failed'
		WHERE id = $1 AND lease_token = $2 AND leased_until > $3
		  AND processed_at IS NULL AND failed_at IS NULL RETURNING id
	`
	FailEventQuery = `
		UPDATE outbox_events SET failed_at = $3, lease_token = NULL, leased_until = NULL,
		    payload_ciphertext = ''::bytea, last_error_redacted = 'delivery_failed'
		WHERE id = $1 AND lease_token = $2 AND leased_until > $3
		  AND processed_at IS NULL AND failed_at IS NULL RETURNING id
	`
	PublishEventQuery = `
		INSERT INTO outbox_events (
			id,
			event_type,
			payload_ciphertext,
			payload_key_version,
			user_id
		) VALUES ($1, $2, $3, $4, NULLIF($5, ''))
	`
)
