package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AfterMCPExecution is the afterMCPExecution hook event.
type AfterMCPExecution struct {
	Envelope
	// ToolName is the MCP tool name (MCP: prefix).
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
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
	registerDecoder(EventAfterMCPExecution, decodeAs[AfterMCPExecution])
}

// AfterMCPExecution registers an afterMCPExecution handler.
func (c *Chain) AfterMCPExecution(fn func(context.Context, AfterMCPExecutionHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterMCPExecution) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}
