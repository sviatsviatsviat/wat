package pretooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/tools"
)

// Event is the preToolUse hook event.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
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
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}
