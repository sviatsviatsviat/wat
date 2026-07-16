package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUseFailure is the postToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the typed tool input (VS Code).
	ToolInput tools.Input `json:"-"`
	// ToolArgs is the typed tool input (camelCase).
	ToolArgs tools.Input `json:"-"`
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

// Input returns tool input from either wire format.
func (e PostToolUseFailure) Input() tools.Input {
	if e.ToolInput.HasRaw() {
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
// Construct via PostToolFailureResults builders. A nil value is a no-op.
type PostToolFailureOutput interface {
	isPostToolFailureOutput()
}

type postToolFailureOutput struct {
	context string
}

func (postToolFailureOutput) isPostToolFailureOutput() {}

func (o postToolFailureOutput) isZero() bool {
	return o.context == ""
}

// PostToolFailureResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolFailureResults interface {
	// Context returns recovery guidance for postToolUseFailure events.
	Context(text string) PostToolFailureOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PostToolFailureOutput
	isPostToolFailureResults()
}

type postToolFailureResults struct{}

func (postToolFailureResults) isPostToolFailureResults() {}

// Context returns recovery guidance for postToolUseFailure events.
func (postToolFailureResults) Context(text string) PostToolFailureOutput {
	return postToolFailureOutput{context: text}
}

// Noop returns an empty response (silent stdout).
func (postToolFailureResults) Noop() PostToolFailureOutput {
	return postToolFailureOutput{}
}

func (postToolFailureOutput) allowedEvents() []string {
	return []string{EventPostToolUseFailure}
}

func (o postToolFailureOutput) encode() ([]byte, int, error) {
	if o.context == "" {
		return nil, 0, nil
	}
	return []byte(o.context), WarnExit, nil
}

func init() {
	registerDecoder(EventPostToolUseFailure, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *PostToolUseFailure, raw []byte) {
			name := e.NativeToolName()
			e.ToolInput = tools.NewInputFromPayload(name, raw, "tool_input")
			e.ToolArgs = tools.NewInputFromPayload(name, raw, "toolArgs")
		})
	})
}

// PostToolUseFailure registers a PostToolUseFailure handler.
func (c *Chain) PostToolUseFailure(fn func(context.Context, Hook[PostToolUseFailure], PostToolFailureResults) (PostToolFailureOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolFailureOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolFailureResults{})
	})
	return &Chain{}
}
