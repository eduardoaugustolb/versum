package events

import (
	"bytes"
	_ "embed"
	"html/template"
	texttemplate "text/template"
	"time"
)

//go:embed templates/magic_link.html
var magicLinkHTML string

//go:embed templates/magic_link.txt
var magicLinkText string

var magicLinkHTMLTemplate = template.Must(template.New("magic-link").Parse(magicLinkHTML))
var magicLinkTextTemplate = texttemplate.Must(texttemplate.New("magic-link").Parse(magicLinkText))

func renderMagicLinkEmail(link string, expiresAt time.Time) (string, string, error) {
	// Absolute expiry keeps the message identical across retries; never render a
	// countdown, random value or current timestamp under the same idempotency key.
	data := struct{ URL, ExpiresAt string }{
		URL:       link,
		ExpiresAt: expiresAt.UTC().Format("02/01/2006 às 15:04:05 UTC"),
	}
	var html, text bytes.Buffer
	if err := magicLinkHTMLTemplate.Execute(&html, data); err != nil {
		return "", "", err
	}
	if err := magicLinkTextTemplate.Execute(&text, data); err != nil {
		return "", "", err
	}
	return html.String(), text.String(), nil
}
