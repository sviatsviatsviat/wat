package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the postToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
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
	registerDecoder(EventPostToolUse, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *PostToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PostToolUse registers a postToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return c
}
