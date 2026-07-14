package claude

// MessageDisplay is the MessageDisplay hook event.
type MessageDisplay struct {
	Envelope
	// TurnID is the turn identifier.
	TurnID string `json:"turn_id"`
	// MessageID is the message identifier.
	MessageID string `json:"message_id"`
	// Index is the message index in the turn.
	Index int `json:"index"`
	// Final is true when this is the final delta.
	Final bool `json:"final"`
	// Delta is the streamed message delta.
	Delta string `json:"delta"`
}

// EventName returns the hook event name.
func (MessageDisplay) EventName() string { return EventMessageDisplay }

func init() {
	registerDecoder(EventMessageDisplay, decodeAs[MessageDisplay])
}

// MessageDisplayOutput is the response for MessageDisplay events.
type MessageDisplayOutput struct {
	Common
	// DisplayContent overrides displayed content when set.
	DisplayContent *string
}

func (o MessageDisplayOutput) isZero() bool {
	return o.Common.isZero() && o.DisplayContent == nil
}
