package notification

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the notification hook event.
type Event struct {
	event.Envelope
	// Message is the notification message.
	Message string `json:"message"`
	// Title is the notification title.
	Title string `json:"title"`
	// NotificationType is the notification category (VS Code).
	NotificationType string `json:"notification_type"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.Notification }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Notification, hookkit.EventDecoder[Event](c))
}
