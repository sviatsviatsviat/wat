package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolUse is the PreToolUse hook event.
type PreToolUse struct {
	Envelope
	// ToolName is the tool name (matcher field).
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the hook event name.
func (PreToolUse) EventName() string { return EventPreToolUse }

func init() {
	registerDecoder(EventPreToolUse, func(raw []byte) (Event, error) {
		return decodeAsAndThen(raw, func(e *PreToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PreToolUseOutput is the response for PreToolUse events.
type PreToolUseOutput struct {
	Common
	// Decision is the permission verdict (allow, deny, ask, defer).
	Decision PermissionDecision
	// Reason is the agent-facing decision reason.
	Reason string
	// UpdatedInput replaces tool arguments when set.
	UpdatedInput map[string]any
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PreToolUseOutput) isZero() bool {
	return o.Common.isZero() && o.Decision == "" && o.Reason == "" &&
		o.UpdatedInput == nil && o.AdditionalContext == ""
}

// PreToolUseResults is the hook-scoped response builder supplied to Chain.PreToolUse handlers by registration.
type PreToolUseResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolUseOutput
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolUseOutput
	// Ask returns an ask verdict with an agent-facing reason.
	Ask(reason string) PreToolUseOutput
	// Defer returns a defer verdict.
	Defer() PreToolUseOutput
	isPreToolUseResults()
}

type preToolUseResults struct{}

func (preToolUseResults) isPreToolUseResults() {}

// Allow returns an allow verdict.
func (preToolUseResults) Allow() PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (preToolUseResults) Deny(reason string) PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionDeny, Reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason.
func (preToolUseResults) Ask(reason string) PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionAsk, Reason: reason}
}

// Defer returns a defer verdict.
func (preToolUseResults) Defer() PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionDefer}
}

// PreToolUse registers a PreToolUse handler.
func (c *Chain) PreToolUse(fn func(context.Context, PreToolUseHook, PreToolUseResults) (PreToolUseOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PreToolUseOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preToolUseResults{})
	})
	return &Chain{}
}
