package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionDenied is the PermissionDenied hook event.
type PermissionDenied struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// Reason is the classifier denial reason.
	Reason string `json:"reason"`
}

// EventName returns the hook event name.
func (PermissionDenied) EventName() string { return EventPermissionDenied }

func init() {
	registerDecoder(EventPermissionDenied, func(raw []byte) (Event, error) {
		return decodeAsAndThen(raw, func(e *PermissionDenied, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PermissionDeniedOutput is the response for PermissionDenied events.
// Construct via PermissionDeniedResults builders and With* methods.
// A nil value is a no-op.
type PermissionDeniedOutput interface {
	isPermissionDeniedOutput()
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) PermissionDeniedOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) PermissionDeniedOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) PermissionDeniedOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) PermissionDeniedOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) PermissionDeniedOutput
}

type permissionDeniedOutput struct {
	common
	retry bool
}

func (permissionDeniedOutput) isPermissionDeniedOutput() {}
func (o permissionDeniedOutput) isZero() bool {
	return o.common.isZero() && !o.retry
}

// WithContinue sets whether Claude should continue the session.
func (o permissionDeniedOutput) WithContinue(v bool) PermissionDeniedOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o permissionDeniedOutput) WithStopReason(reason string) PermissionDeniedOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o permissionDeniedOutput) WithSuppressOutput(v bool) PermissionDeniedOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o permissionDeniedOutput) WithSystemMessage(msg string) PermissionDeniedOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o permissionDeniedOutput) WithTerminalSequence(seq string) PermissionDeniedOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// PermissionDeniedResults is the hook-scoped response builder supplied to On* handlers by registration.
type PermissionDeniedResults interface {
	// Retry returns a retry-requested PermissionDenied result.
	Retry() PermissionDeniedOutput
	isPermissionDeniedResults()
}

type permissionDeniedResults struct{}

func (permissionDeniedResults) isPermissionDeniedResults() {}

// Retry returns a retry-requested PermissionDenied result.
func (permissionDeniedResults) Retry() PermissionDeniedOutput {
	return permissionDeniedOutput{retry: true}
}

func (permissionDeniedOutput) allowedEvents() []string {
	return []string{EventPermissionDenied}
}

func (o permissionDeniedOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.retry {
		hso["retry"] = true
	}
}

// OnPermissionDenied registers a PermissionDenied handler.
func OnPermissionDenied(fn func(context.Context, Hook[PermissionDenied], PermissionDeniedResults) (PermissionDeniedOutput, error)) *chain {
	return (&chain{}).PermissionDenied(fn)
}

// PermissionDenied registers another PermissionDenied handler on the chain.
func (c *chain) PermissionDenied(fn func(context.Context, Hook[PermissionDenied], PermissionDeniedResults) (PermissionDeniedOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PermissionDenied) (PermissionDeniedOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionDeniedResults{})
	})
	return c
}
