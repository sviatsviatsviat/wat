package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the PostToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
	// ToolResult is the tool result.
	ToolResult ToolResult `json:"tool_result"`
}

// EventName returns the canonical hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

// NativeToolName returns the tool name.
func (e PostToolUse) NativeToolName() string {
	return e.ToolName
}

// Input returns tool input.
func (e PostToolUse) Input() tools.Input {
	return e.ToolInput
}

// ResultText returns the textual tool result.
func (e PostToolUse) ResultText() string {
	return e.ToolResult.Text()
}

// ResultRaw returns the tool result JSON.
func (e PostToolUse) ResultRaw() json.RawMessage {
	if e.ToolResult.TextResultForLLM != "" || e.ToolResult.ResultType != "" {
		return marshalToolResult(e.ToolResult)
	}
	return nil
}

// PostToolOutput is the response for PostToolUse events.
// Construct via PostToolResults builders and With* methods. A nil value is a no-op.
type PostToolOutput interface {
	isPostToolOutput()
	// WithModifiedResult replaces the tool result text when set.
	WithModifiedResult(result string) PostToolOutput
}

type postToolOutput struct {
	modifiedResult    string
	additionalContext string
}

func (postToolOutput) isPostToolOutput() {}

func (o postToolOutput) isZero() bool {
	return o.modifiedResult == "" && o.additionalContext == ""
}

// WithModifiedResult replaces the tool result text when set.
func (o postToolOutput) WithModifiedResult(result string) PostToolOutput {
	o.modifiedResult = result
	return o
}

// PostToolResults is the hook-scoped response builder supplied to On* handlers by registration.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PostToolOutput
	isPostToolResults()
}

type postToolResults struct{}

func (postToolResults) isPostToolResults() {}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolOutput {
	return postToolOutput{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (postToolResults) Noop() PostToolOutput {
	return postToolOutput{}
}

func (postToolOutput) allowedEvents() []string {
	return []string{EventPostToolUse}
}

func (o postToolOutput) encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.modifiedResult != "" {
		out["modified_result"] = map[string]any{
			"result_type":         "success",
			"text_result_for_llm": o.modifiedResult,
		}
	}
	if o.additionalContext != "" {
		out["additional_context"] = o.additionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	codec.Register(EventPostToolUse, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PostToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.NativeToolName(), raw, "tool_input")
		})
	})
}

// OnPostToolUse registers a PostToolUse handler.
func OnPostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	return (&chain{}).PostToolUse(fn)
}

// PostToolUse registers another PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return c
}
