package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUseFailure is the postToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
	// Error is the failure payload (string or object).
	Error json.RawMessage `json:"error"`
}

// EventName returns the canonical hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

// NativeToolName returns the tool name from either wire format.
func (e PostToolUseFailure) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input JSON from either wire format.
func (e PostToolUseFailure) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ErrorMessage returns the failure message from the error field.
func (e PostToolUseFailure) ErrorMessage() string {
	if len(e.Error) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(e.Error, &s) == nil {
		return s
	}
	var detail ErrorDetail
	if json.Unmarshal(e.Error, &detail) == nil {
		return detail.Message
	}
	return string(e.Error)
}

// PostToolFailureOutput is the response for postToolUseFailure events.
type PostToolFailureOutput struct {
	// Context is recovery guidance written as raw stdout text.
	Context string
}

func (o PostToolFailureOutput) isZero() bool {
	return o.Context == ""
}

// PostToolFailureResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolFailureResults interface {
	// Context returns recovery guidance for postToolUseFailure events.
	Context(text string) PostToolFailureOutput
	isPostToolFailureResults()
}

type postToolFailureResults struct{}

func (postToolFailureResults) isPostToolFailureResults() {}

// Context returns recovery guidance for postToolUseFailure events.
func (postToolFailureResults) Context(text string) PostToolFailureOutput {
	return PostToolFailureOutput{Context: text}
}

func encodePostToolFailure(o PostToolFailureOutput) ([]byte, int, error) {
	if o.Context == "" {
		return nil, 0, nil
	}
	return []byte(o.Context), WarnExit, nil
}

func init() {
	registerDecoder(EventPostToolUseFailure, decodeAs[PostToolUseFailure])
}

// PostToolUseFailure registers a PostToolUseFailure handler.
func (c *Chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailureHook, PostToolFailureResults) (PostToolFailureOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolFailureOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolFailureResults{})
	})
	return &Chain{}
}
