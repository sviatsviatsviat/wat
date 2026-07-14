package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUseFailure is the PostToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// Error is the failure message.
	Error string `json:"error"`
}

// EventName returns the hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

func init() {
	registerDecoder(EventPostToolUseFailure, decodeAs[PostToolUseFailure])
}

// PostToolUseFailureResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolUseFailureResults interface {
	// Context returns recovery guidance for PostToolUseFailure events.
	Context(text string) PostToolUseOutput
	isPostToolUseFailureResults()
}

type postToolUseFailureResults struct{}

func (postToolUseFailureResults) isPostToolUseFailureResults() {}

// Context returns recovery guidance for PostToolUseFailure events.
func (postToolUseFailureResults) Context(text string) PostToolUseOutput {
	return PostToolUseOutput{AdditionalContext: text}
}

// PostToolUseFailure registers a PostToolUseFailure handler.
func (c *Chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailureHook, PostToolUseFailureResults) (PostToolUseOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolUseOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolUseFailureResults{})
	})
	return &Chain{}
}
