package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the PostToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolResponse is the tool response JSON.
	ToolResponse json.RawMessage `json:"tool_response"`
	// DurationMs is the tool execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

func init() {
	registerDecoder(EventPostToolUse, decodeAs[PostToolUse])
}

// PostToolUseOutput is the response for PostToolUse and PostToolUseFailure events.
type PostToolUseOutput struct {
	Common
	// Block prompts Claude with Reason when true.
	Block bool
	// Reason is the block reason.
	Reason string
	// AdditionalContext injects model context.
	AdditionalContext string
	// UpdatedToolOutput replaces the tool result when set.
	UpdatedToolOutput any
}

func (o PostToolUseOutput) isZero() bool {
	return o.Common.isZero() && !o.Block && o.Reason == "" &&
		o.AdditionalContext == "" && o.UpdatedToolOutput == nil
}

// PostToolUseResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolUseResults interface {
	// Context returns a context-injection-only PostToolUse result.
	Context(text string) PostToolUseOutput
	// Block returns a block result with an agent-facing reason.
	Block(reason string) PostToolUseOutput
	isPostToolUseResults()
}

type postToolUseResults struct{}

func (postToolUseResults) isPostToolUseResults() {}

// Context returns a context-injection-only PostToolUse result.
func (postToolUseResults) Context(text string) PostToolUseOutput {
	return PostToolUseOutput{AdditionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (postToolUseResults) Block(reason string) PostToolUseOutput {
	return PostToolUseOutput{Block: true, Reason: reason}
}

// PostToolUse registers a PostToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, PostToolUseHook, PostToolUseResults) (PostToolUseOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolUseOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolUseResults{})
	})
	return &Chain{}
}
