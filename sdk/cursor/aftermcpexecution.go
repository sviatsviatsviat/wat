package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AfterMCPExecution is the afterMCPExecution hook event.
type AfterMCPExecution struct {
	Envelope
	hookkit.RawPayload
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
func (AfterMCPExecution) EventName() string { return EventAfterMCPExecution }

// DurationMillis returns the execution duration in milliseconds.
func (e AfterMCPExecution) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

func init() {
	registerDecoder(EventAfterMCPExecution, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *AfterMCPExecution, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// OnAfterMCPExecution registers an afterMCPExecution handler.
func OnAfterMCPExecution(fn func(context.Context, Hook[AfterMCPExecution], PostToolResults) (PostToolOutput, error)) *chain {
	return (&chain{}).AfterMCPExecution(fn)
}

// AfterMCPExecution registers another AfterMCPExecution handler on the chain.
func (c *chain) AfterMCPExecution(fn func(context.Context, Hook[AfterMCPExecution], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterMCPExecution) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return c
}
