package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the PostToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
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
	codec.Register(EventPostToolUse, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PostToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PostToolUseOutput is the response for PostToolUse and PostToolUseFailure events.
// Construct via PostToolUseResults / PostToolUseFailureResults builders and With* methods.
// A nil value is a no-op.
type PostToolUseOutput interface {
	Output
	isPostToolUseOutput()
	// WithUpdatedToolOutput replaces the tool result when set.
	WithUpdatedToolOutput(output any) PostToolUseOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) PostToolUseOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) PostToolUseOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) PostToolUseOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) PostToolUseOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) PostToolUseOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) PostToolUseOutput
}

type postToolUseOutput struct {
	common
	block             bool
	reason            string
	additionalContext string
	updatedToolOutput any
}

func (postToolUseOutput) isClaudeOutput() {}

func (postToolUseOutput) isPostToolUseOutput() {}
func (o postToolUseOutput) isZero() bool {
	return o.common.isZero() && !o.block && o.reason == "" &&
		o.additionalContext == "" && o.updatedToolOutput == nil
}

// WithUpdatedToolOutput replaces the tool result when set.
func (o postToolUseOutput) WithUpdatedToolOutput(output any) PostToolUseOutput {
	o.updatedToolOutput = output
	return o
}

// WithAdditionalContext injects model context.
func (o postToolUseOutput) WithAdditionalContext(text string) PostToolUseOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o postToolUseOutput) WithContinue(v bool) PostToolUseOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o postToolUseOutput) WithStopReason(reason string) PostToolUseOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o postToolUseOutput) WithSuppressOutput(v bool) PostToolUseOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o postToolUseOutput) WithSystemMessage(msg string) PostToolUseOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o postToolUseOutput) WithTerminalSequence(seq string) PostToolUseOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// PostToolUseResults is the hook-scoped response builder supplied to On* handlers by registration.
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
	return postToolUseOutput{additionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (postToolUseResults) Block(reason string) PostToolUseOutput {
	return postToolUseOutput{block: true, reason: reason}
}

func (postToolUseOutput) allowedEvents() []string {
	return []string{EventPostToolUse, EventPostToolUseFailure}
}

func (o postToolUseOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.block {
		top["decision"] = "block"
		if o.reason != "" {
			top["reason"] = o.reason
		}
	}
	if o.updatedToolOutput != nil {
		hso["updatedToolOutput"] = o.updatedToolOutput
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// OnPostToolUse registers a PostToolUse handler.
func OnPostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolUseResults) (PostToolUseOutput, error)) *chain {
	return (&chain{}).PostToolUse(fn)
}

// PostToolUse registers another PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolUseResults) (PostToolUseOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolUseOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), postToolUseResults{})
	})
	return c
}
