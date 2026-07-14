package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the postToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
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

// Input returns tool input JSON from either wire format.
func (e PostToolUse) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
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
type PostToolOutput struct {
	// ModifiedResult replaces the tool result text when set.
	ModifiedResult string
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PostToolOutput) isZero() bool {
	return o.ModifiedResult == "" && o.AdditionalContext == ""
}

// PostToolResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolOutput
	isPostToolResults()
}

type postToolResults struct{}

func (postToolResults) isPostToolResults() {}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolOutput {
	return PostToolOutput{AdditionalContext: text}
}

func encodePostTool(o PostToolOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.ModifiedResult != "" {
		out["modifiedResult"] = map[string]any{
			"resultType":       "success",
			"textResultForLlm": o.ModifiedResult,
		}
	}
	if o.AdditionalContext != "" {
		out["additionalContext"] = o.AdditionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventPostToolUse, decodeAs[PostToolUse])
}

// PostToolUse registers a PostToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, PostToolUseHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}
