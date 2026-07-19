package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUseFailure is the postToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ErrorMessage is the failure message.
	ErrorMessage string `json:"error_message"`
	// FailureType is the failure category.
	FailureType string `json:"failure_type"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

// DurationMillis returns the execution duration in milliseconds.
func (e PostToolUseFailure) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

func init() {
	registerDecoder(EventPostToolUseFailure, func(raw []byte) (Event, error) {
		return decodeAsAndThen(raw, func(e *PostToolUseFailure, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// OnPostToolUseFailure registers a postToolUseFailure handler.
func OnPostToolUseFailure(fn func(context.Context, Hook[PostToolUseFailure], PostToolResults) (PostToolOutput, error)) *chain {
	return (&chain{}).PostToolUseFailure(fn)
}

// PostToolUseFailure registers another PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, Hook[PostToolUseFailure], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return c
}
