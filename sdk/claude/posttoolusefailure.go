package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	codec.Register(EventPostToolUseFailure, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PostToolUseFailure, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PostToolUseFailureResults is the hook-scoped response builder supplied to On* handlers by registration.
type PostToolUseFailureResults interface {
	// Context returns recovery guidance for PostToolUseFailure events.
	Context(text string) PostToolUseOutput
	isPostToolUseFailureResults()
}

type postToolUseFailureResults struct{}

func (postToolUseFailureResults) isPostToolUseFailureResults() {}

// Context returns recovery guidance for PostToolUseFailure events.
func (postToolUseFailureResults) Context(text string) PostToolUseOutput {
	return postToolUseOutput{eventName: EventPostToolUseFailure, additionalContext: text}
}

// PostToolUseFailure registers a PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolUseFailureResults) (PostToolUseOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[PostToolUseFailure]) (PostToolUseOutput, error) {
		return fn(ctx, hook, postToolUseFailureResults{})
	}))
	return c
}
