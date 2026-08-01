package pretooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
)

// Event is the PreToolUse hook event.
type Event struct {
	event.Envelope
	event.ToolFields
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PreToolUse }

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e Event) ShellCommand() string {
	if !tools.IsShellToolName(e.NativeToolName()) {
		return ""
	}
	return hookkit.ExtractShellCommand(e.Input().Raw())
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PreToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEventErr(c, raw, func(e *Event, payload []byte) error {
			return e.BindToolInput(payload)
		})
	})
}
