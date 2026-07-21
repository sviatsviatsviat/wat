package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	codec.Register(EventPreToolUse, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PreToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PreToolUseOutput is the response for PreToolUse events.
// Construct via PreToolUseResults builders and With* methods. A nil value is a no-op.
type PreToolUseOutput interface {
	run.Output
	isPreToolUseOutput()
	// WithUpdatedInput replaces tool arguments when set.
	WithUpdatedInput(input map[string]any) PreToolUseOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) PreToolUseOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) PreToolUseOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) PreToolUseOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) PreToolUseOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) PreToolUseOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) PreToolUseOutput
}

type preToolUseOutput struct {
	common
	decision          PermissionDecision
	reason            string
	updatedInput      map[string]any
	additionalContext string
}

func (preToolUseOutput) isPreToolUseOutput() {}

// IsZero reports whether this hook response is empty.
func (o preToolUseOutput) IsZero() bool {
	return o.common.IsZero() && o.decision == "" && o.reason == "" &&
		o.updatedInput == nil && o.additionalContext == ""
}

// WithUpdatedInput replaces tool arguments when set.
func (o preToolUseOutput) WithUpdatedInput(input map[string]any) PreToolUseOutput {
	o.updatedInput = input
	return o
}

// WithAdditionalContext injects model context.
func (o preToolUseOutput) WithAdditionalContext(text string) PreToolUseOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o preToolUseOutput) WithContinue(v bool) PreToolUseOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o preToolUseOutput) WithStopReason(reason string) PreToolUseOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o preToolUseOutput) WithSuppressOutput(v bool) PreToolUseOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o preToolUseOutput) WithSystemMessage(msg string) PreToolUseOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o preToolUseOutput) WithTerminalSequence(seq string) PreToolUseOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// PreToolUseResults is the hook-scoped response builder supplied to OnPreToolUse handlers by registration.
type PreToolUseResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolUseOutput
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolUseOutput
	// Ask returns an ask verdict with an agent-facing reason.
	Ask(reason string) PreToolUseOutput
	// Defer returns a defer verdict.
	Defer() PreToolUseOutput
	// Noop returns an empty response (silent stdout). Prefer returning nil from handlers when not chaining With*.
	Noop() PreToolUseOutput
	isPreToolUseResults()
}

type preToolUseResults struct{}

func (preToolUseResults) isPreToolUseResults() {}

// Allow returns an allow verdict.
func (preToolUseResults) Allow() PreToolUseOutput {
	return preToolUseOutput{decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (preToolUseResults) Deny(reason string) PreToolUseOutput {
	return preToolUseOutput{decision: DecisionDeny, reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason.
func (preToolUseResults) Ask(reason string) PreToolUseOutput {
	return preToolUseOutput{decision: DecisionAsk, reason: reason}
}

// Defer returns a defer verdict.
func (preToolUseResults) Defer() PreToolUseOutput {
	return preToolUseOutput{decision: DecisionDefer}
}

// Noop returns an empty response (silent stdout).
func (preToolUseResults) Noop() PreToolUseOutput {
	return preToolUseOutput{}
}

func (o preToolUseOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.decision != "" {
		hso["permissionDecision"] = string(o.decision)
		if o.reason != "" {
			hso["permissionDecisionReason"] = o.reason
		}
	} else if o.updatedInput != nil {
		hso["permissionDecision"] = "allow"
	}
	if o.updatedInput != nil {
		hso["updatedInput"] = o.updatedInput
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o preToolUseOutput) Encode() ([]byte, int, error) {
	return marshalHookOutput(EventPreToolUse, o.encodeInto)
}

// Merge combines other into this PreToolUse output.
func (o preToolUseOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(preToolUseOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.common.Merge(b.common)
	decision, reason := hookkit.MergeRankedString(
		string(o.decision), o.reason,
		string(b.decision), b.reason,
		hookkit.PermissionRankString,
	)
	updatedInput, w := hookkit.TakeLastMap("updatedInput", o.updatedInput, b.updatedInput)
	if w != "" {
		warnings = append(warnings, w)
	}
	return preToolUseOutput{
		common:            mergedCommon,
		decision:          PermissionDecision(decision),
		reason:            reason,
		updatedInput:      updatedInput,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o preToolUseOutput) Stop() bool {
	return o.common.Stop() || o.decision == DecisionDeny
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, run.Hook[PreToolUse], PreToolUseResults) (PreToolUseOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[PreToolUse]) (PreToolUseOutput, error) {
		return fn(ctx, hook, preToolUseResults{})
	}))
	return c
}
