package domain

import "strings"

type EventType string

func ParseEventType(raw string) (EventType, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrInvalidEventType
	}

	return EventType(raw), nil
}

func (e EventType) String() string {
	return string(e)
}
