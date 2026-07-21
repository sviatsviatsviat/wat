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
	codec.Register(EventAfterMCPExecution, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *AfterMCPExecution, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// AfterMCPExecution registers a AfterMCPExecution handler on the chain.
func (c *chain) AfterMCPExecution(fn func(context.Context, run.Hook[AfterMCPExecution], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[AfterMCPExecution]) (PostToolOutput, error) {
		return fn(ctx, hook, postToolResults{})
	}))
	return c
}
