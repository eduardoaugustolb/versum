package postgres

const (
	PublishEventQuery = `
		INSERT INTO outbox_events (
			id,
			event_type,
			payload_ciphertext,
			payload_key_version
		) VALUES ($1, $2, $3, $4)
	`
)
