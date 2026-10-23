package resend

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/email"
	"github.com/resend/resend-go/v4"
)

type ResendSender struct {
	from   string
	client *resend.Client
}

var _ email.Sender = (*ResendSender)(nil)

func NewResendSender(apiKey, from string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (s *ResendSender) Send(ctx context.Context, req email.SendRequest) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{req.To},
		Subject: req.Subject,
		Html:    req.HTML,
		Text:    req.Text,
	}
	options := &resend.SendEmailOptions{
		IdempotencyKey: req.IdempotencyKey,
	}

	if _, err := s.client.Emails.SendWithOptions(ctx, params, options); err != nil {
		return err
	}

	return nil
}
