package postgres

const (
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
