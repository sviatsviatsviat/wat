package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUseFailure is the PostToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// Error is the failure message.
	Error string `json:"error"`
	// IsInterrupt is true when the failure was caused by an interrupt.
	IsInterrupt bool `json:"is_interrupt"`
	// DurationMs is the tool execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

func init() {
	registerDecoder(EventPostToolUseFailure, func(raw []byte) (Event, error) {
		return decodeAsAndThen(raw, func(e *PostToolUseFailure, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
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
