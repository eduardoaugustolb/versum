package email

import "context"

type SendRequest struct {
	To             string
	Subject        string
	HTML           string
	Text           string
	IdempotencyKey string
}

type Sender interface {
	Send(ctx context.Context, req SendRequest) error
}
