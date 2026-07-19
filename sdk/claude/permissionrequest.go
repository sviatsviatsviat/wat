package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionRequest is the PermissionRequest hook event.
type PermissionRequest struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the hook event name.
func (PermissionRequest) EventName() string { return EventPermissionRequest }

func init() {
	codec.Register(EventPermissionRequest, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PermissionRequest, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PermissionRequestOutput is the response for PermissionRequest events.
// Construct via PermissionRequestResults builders and With* methods.
// A nil value is a no-op.
type PermissionRequestOutput interface {
	isPermissionRequestOutput()
	// WithUpdatedInput replaces tool arguments when set.
	WithUpdatedInput(input map[string]any) PermissionRequestOutput
	// WithInterrupt stops the session when true.
	WithInterrupt(v bool) PermissionRequestOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) PermissionRequestOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) PermissionRequestOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) PermissionRequestOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) PermissionRequestOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) PermissionRequestOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) PermissionRequestOutput
}

type permissionRequestOutput struct {
	common
	behavior          string
	updatedInput      map[string]any
	message           string
	interrupt         bool
	additionalContext string
}

func (permissionRequestOutput) isPermissionRequestOutput() {}
func (o permissionRequestOutput) isZero() bool {
	return o.common.isZero() && o.behavior == "" && o.updatedInput == nil &&
		o.message == "" && !o.interrupt && o.additionalContext == ""
}

// WithUpdatedInput replaces tool arguments when set.
func (o permissionRequestOutput) WithUpdatedInput(input map[string]any) PermissionRequestOutput {
	o.updatedInput = input
	return o
}

// WithInterrupt stops the session when true.
func (o permissionRequestOutput) WithInterrupt(v bool) PermissionRequestOutput {
	o.interrupt = v
	return o
}

// WithAdditionalContext injects model context.
func (o permissionRequestOutput) WithAdditionalContext(text string) PermissionRequestOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o permissionRequestOutput) WithContinue(v bool) PermissionRequestOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o permissionRequestOutput) WithStopReason(reason string) PermissionRequestOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o permissionRequestOutput) WithSuppressOutput(v bool) PermissionRequestOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o permissionRequestOutput) WithSystemMessage(msg string) PermissionRequestOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o permissionRequestOutput) WithTerminalSequence(seq string) PermissionRequestOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// PermissionRequestResults is the hook-scoped response builder supplied to On* handlers by registration.
type PermissionRequestResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionRequestOutput
	// Deny returns a deny verdict with a permission message.
	Deny(message string) PermissionRequestOutput
	isPermissionRequestResults()
}

type permissionRequestResults struct{}

func (permissionRequestResults) isPermissionRequestResults() {}

// Allow returns an allow verdict.
func (permissionRequestResults) Allow() PermissionRequestOutput {
	return permissionRequestOutput{behavior: "allow"}
}

// Deny returns a deny verdict with a permission message.
func (permissionRequestResults) Deny(message string) PermissionRequestOutput {
	return permissionRequestOutput{behavior: "deny", message: message}
}

func (permissionRequestOutput) allowedEvents() []string {
	return []string{EventPermissionRequest}
}

func (o permissionRequestOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.behavior != "" {
		dec := map[string]any{"behavior": o.behavior}
		if o.updatedInput != nil {
			dec["updatedInput"] = o.updatedInput
		}
		if o.message != "" {
			dec["message"] = o.message
		}
		if o.interrupt {
			dec["interrupt"] = true
		}
		hso["decision"] = dec
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// OnPermissionRequest registers a PermissionRequest handler.
func OnPermissionRequest(fn func(context.Context, run.Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	return (&chain{}).PermissionRequest(fn)
}

// PermissionRequest registers another PermissionRequest handler on the chain.
func (c *chain) PermissionRequest(fn func(context.Context, run.Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PermissionRequest) (PermissionRequestOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), permissionRequestResults{})
	})
	return c
}
