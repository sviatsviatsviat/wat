package notification

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the Notification hook event.
type Event struct {
	event.Envelope
	// Message is the notification message.
	Message string `json:"message"`
	// NotificationType is the notification category.
	NotificationType string `json:"notification_type"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.Notification }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Notification, hookkit.EventDecoder[Event](c))
}
