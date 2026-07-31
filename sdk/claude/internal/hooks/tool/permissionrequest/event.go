package permissionrequest

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the PermissionRequest hook event.
type Event struct {
	event.Envelope
	event.ToolFields
	// PermissionSuggestions are optional "always allow" options from the
	// permission dialog. Echo one or more as updatedPermissions on allow to
	// apply them without prompting the user.
	PermissionSuggestions []event.PermissionUpdate `json:"permission_suggestions,omitempty"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PermissionRequest }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PermissionRequest, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.BindToolInput(raw)
			return nil
		})
	})
}
