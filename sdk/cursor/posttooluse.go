package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the postToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolOutput is the tool output text.
	ToolOutput string `json:"tool_output"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

// DurationMillis returns the execution duration in milliseconds.
func (e PostToolUse) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

func init() {
	registerDecoder(EventPostToolUse, decodeAs[PostToolUse])
}

// PostToolUse registers a postToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, PostToolUseHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}
