package permissionrequest

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the PermissionRequest hook event.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PermissionRequest }

// NativeToolName returns the tool name.
func (e Event) NativeToolName() string {
	return e.ToolName
}

// Input returns tool input.
func (e Event) Input() tools.Input {
	return e.ToolInput
}

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e Event) ShellCommand() string {
	if !tools.IsShellToolName(e.NativeToolName()) {
		return ""
	}
	return hookkit.ExtractShellCommand(e.Input().Raw())
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.PermissionRequest, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, payload []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.NativeToolName(), payload, "tool_input")
		})
	})
}
