package permissiondenied

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the PermissionDenied hook event.
type Event struct {
	event.Envelope
	event.ToolFields
	// Reason is the classifier denial reason.
	Reason string `json:"reason"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PermissionDenied }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PermissionDenied, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.BindToolInput(raw)
			return nil
		})
	})
}
