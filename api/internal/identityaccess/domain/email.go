package domain

import (
	"regexp"
	"strings"
)

type Email struct {
	value string
}

func ParseEmail(raw string) (Email, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Email{}, ErrInvalidEmail
	}

	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+` +
			`(?:\.[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*` +
			`@` +
			`[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?` +
			`(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`,
	)
	isValid := emailRegex.MatchString(trimmed)

	if !isValid {
		return Email{}, ErrInvalidEmail
	}
	capitalized := strings.ToLower(trimmed)

	return Email{value: capitalized}, nil
}

func (e Email) String() string {
	return e.value
}
