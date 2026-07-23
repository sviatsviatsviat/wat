package permissionrequest

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the PermissionRequest hook event.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PermissionRequest }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.PermissionRequest, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}
