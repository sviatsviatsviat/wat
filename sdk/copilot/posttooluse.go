package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the postToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the typed tool input (VS Code).
	ToolInput tools.Input `json:"-"`
	// ToolArgs is the typed tool input (camelCase).
	ToolArgs tools.Input `json:"-"`
	// ToolResult is the tool result (camelCase).
	ToolResult ToolResult `json:"toolResult"`
	// ToolResultSnake is the tool result (VS Code).
	ToolResultSnake ToolResult `json:"tool_result"`
}

// EventName returns the canonical hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

// NativeToolName returns the tool name from either wire format.
func (e PostToolUse) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input from either wire format.
func (e PostToolUse) Input() tools.Input {
	if e.ToolInput.HasRaw() {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ResultText returns the textual tool result from either wire format.
func (e PostToolUse) ResultText() string {
	if t := e.ToolResult.Text(); t != "" {
		return t
	}
	return e.ToolResultSnake.Text()
}

// ResultRaw returns the tool result JSON from either wire format.
func (e PostToolUse) ResultRaw() json.RawMessage {
	if raw := extractRawObjectField(e.DecodedRaw(), "toolResult", "tool_result"); raw != nil {
		return raw
	}
	if e.ToolResult.TextResultForLLM != "" || e.ToolResult.ResultType != "" {
		return marshalToolResultCamel(e.ToolResult)
	}
	if e.ToolResultSnake.TextResultSnake != "" || e.ToolResultSnake.ResultTypeSnake != "" {
		return marshalToolResultSnake(e.ToolResultSnake)
	}
	return nil
}

// PostToolOutput is the response for postToolUse events.
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
		out["modifiedResult"] = map[string]any{
			"resultType":       "success",
			"textResultForLlm": o.modifiedResult,
		}
	}
	if o.additionalContext != "" {
		out["additionalContext"] = o.additionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventPostToolUse, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *PostToolUse, raw []byte) {
			name := e.NativeToolName()
			e.ToolInput = tools.NewInputFromPayload(name, raw, "tool_input")
			e.ToolArgs = tools.NewInputFromPayload(name, raw, "toolArgs")
		})
	})
}

// OnPostToolUse registers a PostToolUse handler.
func OnPostToolUse(fn func(context.Context, Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	return (&chain{}).PostToolUse(fn)
}

// PostToolUse registers another PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, Hook[PostToolUse], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return c
}
