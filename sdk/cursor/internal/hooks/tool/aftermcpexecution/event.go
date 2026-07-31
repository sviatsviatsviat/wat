package aftermcpexecution

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterMCPExecution hook event.
//
// Cursor Hooks docs list input fields only (tool_name, tool_input, result_json,
// duration); there are no host-honored output fields. Handlers are observe-only.
// Cloud agents do not load beforeMCPExecution / afterMCPExecution hooks.
type Event struct {
	event.Envelope
	event.ToolFields
	event.DurationFields
	// ResultJSON is the MCP result JSON text.
	ResultJSON string `json:"result_json"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterMCPExecution }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterMCPExecution, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.BindToolInput(raw)
			e.CaptureDurationPresent(raw)
			return nil
		})
	})
}
