package domain

type Event struct {
	id        string
	eventType EventType
	payload   []byte
}

func NewEvent(
	id string,
	rawEventType string,
	payload []byte,
) (*Event, error) {

	eventType, err := ParseEventType(rawEventType)
	if err != nil {
		return nil, err
	}

	e := &Event{
		id:        id,
		eventType: eventType,
		payload:   append([]byte{}, payload...),
	}
	if err := validateEvent(e); err != nil {
		return nil, err
	}
	return e, nil
}

func RehydrateEvent(
	id string,
	rawEventType string,
	payload []byte,
) (*Event, error) {
	eventType, err := ParseEventType(rawEventType)
	if err != nil {
		return nil, err
	}
	e := &Event{
		id:        id,
		eventType: eventType,
		payload:   append([]byte{}, payload...),
	}
	if err := validateEvent(e); err != nil {
		return nil, err
	}
	return e, nil
}

func validateEvent(e *Event) error {
	if e.id == "" {
		return ErrInvalidId
	}
	if e.eventType == "" {
		return ErrInvalidEventType
	}

	if len(e.payload) == 0 {
		return ErrInvalidPayload
	}
	return nil
}

func (e *Event) ID() string {
	return e.id
}
func (e *Event) EventType() EventType {
	return e.eventType
}
func (e *Event) Payload() []byte {
	return append([]byte{}, e.payload...)
}
