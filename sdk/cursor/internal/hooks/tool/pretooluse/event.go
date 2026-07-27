package pretooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the preToolUse hook event.
type Event struct {
	event.Envelope
	event.ToolFields
	// AgentMessage is the pre-call narrative from the agent when present.
	AgentMessage string `json:"agent_message"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PreToolUse }

// ShellCommand extracts the shell command when the tool is Shell.
func (e Event) ShellCommand() string {
	if e.ToolName != "Shell" && e.ToolName != "shell" {
		return ""
	}
	return hookkit.ExtractShellCommand(e.ToolInput.Raw())
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PreToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.BindToolInput(raw)
		})
	})
}
