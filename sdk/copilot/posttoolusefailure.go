package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUseFailure is the PostToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
	// Error is the failure payload (string or object).
	Error json.RawMessage `json:"error"`
}

// EventName returns the canonical hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

// NativeToolName returns the tool name.
func (e PostToolUseFailure) NativeToolName() string {
	return e.ToolName
}

// Input returns tool input.
func (e PostToolUseFailure) Input() tools.Input {
	return e.ToolInput
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

// PostToolFailureOutput is the response for PostToolUseFailure events.
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

// PostToolFailureResults is the hook-scoped response builder supplied to On* handlers by registration.
type PostToolFailureResults interface {
	// Context returns recovery guidance for PostToolUseFailure events.
	Context(text string) PostToolFailureOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PostToolFailureOutput
	isPostToolFailureResults()
}

type postToolFailureResults struct{}

func (postToolFailureResults) isPostToolFailureResults() {}

// Context returns recovery guidance for PostToolUseFailure events.
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
	codec.Register(EventPostToolUseFailure, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PostToolUseFailure, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.NativeToolName(), raw, "tool_input")
		})
	})
}

// OnPostToolUseFailure registers a PostToolUseFailure handler.
func OnPostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolFailureResults) (PostToolFailureOutput, error)) *chain {
	return (&chain{}).PostToolUseFailure(fn)
}

// PostToolUseFailure registers another PostToolUseFailure handler on the chain.
func (c *chain) PostToolUseFailure(fn func(context.Context, run.Hook[PostToolUseFailure], PostToolFailureResults) (PostToolFailureOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolFailureOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), postToolFailureResults{})
	})
	return c
}
