package beforeshellexecution

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeShellExecution hook event.
//
// Permission responses use allow, deny, or ask. Unlike preToolUse (where ask is
// accepted by the schema but not enforced) and subagentStart (where ask is
// treated as deny), Cursor enforces ask on this event and escalates to the user.
// Deny writes agent_message by default and exits with PermissionDenyExit (2);
// chain WithUserMessage for a client-facing message. Prefer Deny when blocking
// without prompting, and Ask when the host should request approval.
type Event struct {
	event.Envelope
	// Command is the shell command about to run.
	Command string `json:"command"`
	// Sandbox reports whether the command runs in a sandbox.
	Sandbox bool `json:"sandbox"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.BeforeShellExecution }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.BeforeShellExecution, hookkit.EventDecoder[Event](c))
}
