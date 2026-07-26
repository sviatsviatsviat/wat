package workspaceopen

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the workspaceOpen hook event.
//
// workspaceOpen is an app lifecycle hook: it fires in the Cursor desktop app
// and CLI when a workspace opens or its folders change, and is skipped when
// there are zero workspace folders. It does not run in cloud agents.
type Event struct {
	event.Envelope
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.WorkspaceOpen }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.WorkspaceOpen, hookkit.EventDecoder[Event](c))
}
