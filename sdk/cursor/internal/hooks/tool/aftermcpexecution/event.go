package aftermcpexecution

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/tools"
)

// Event is the afterMCPExecution hook event.
//
// Cursor Hooks docs list input fields only (tool_name, tool_input, result_json,
// duration); there are no host-honored output fields. Handlers are observe-only.
// Cloud agents do not load beforeMCPExecution / afterMCPExecution hooks.
type Event struct {
	event.Envelope
	event.DurationFields
	// ToolName is the native tool name (typically MCP:<tool>).
	ToolName string `json:"tool_name"`
	// ToolInput is the tool arguments from tool_input, bound to ToolName after decode.
	ToolInput tools.Input `json:"-"`
	// ResultJSON is the MCP result JSON text.
	ResultJSON string `json:"result_json"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterMCPExecution }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterMCPExecution, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
			e.CaptureDurationPresent(raw)
		})
	})
}
