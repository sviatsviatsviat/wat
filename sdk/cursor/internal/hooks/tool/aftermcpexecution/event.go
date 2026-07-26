package aftermcpexecution

import (
	"context"

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
	// ToolName is the native tool name (typically MCP:<tool>).
	ToolName string `json:"tool_name"`
	// ToolInput is the tool arguments from tool_input, bound to ToolName after decode.
	ToolInput tools.Input `json:"-"`
	// ResultJSON is the MCP result JSON text.
	ResultJSON string `json:"result_json"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterMCPExecution }

// DurationMillis returns the execution duration in milliseconds.
func (e Event) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterMCPExecution, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// RegisterHandler registers an observe-only AfterMCPExecution handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
