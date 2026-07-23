package pretooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for this hook event.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
	// WithUpdatedInput replaces tool arguments when set.
	WithUpdatedInput(input map[string]any) Output
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) Output
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) Output
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) Output
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) Output
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) Output
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) Output
}

type output struct {
	event.Common
	decision          event.PermissionDecision
	reason            string
	updatedInput      map[string]any
	additionalContext string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && o.decision == "" && o.reason == "" &&
		o.updatedInput == nil && o.additionalContext == ""
}

// WithUpdatedInput replaces tool arguments when set.
func (o output) WithUpdatedInput(input map[string]any) Output {
	o.updatedInput = input
	return o
}

// WithAdditionalContext injects model context.
func (o output) WithAdditionalContext(text string) Output {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o output) WithContinue(v bool) Output {
	o.Common = o.Common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o output) WithStopReason(reason string) Output {
	o.Common = o.Common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o output) WithSuppressOutput(v bool) Output {
	o.Common = o.Common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o output) WithSystemMessage(msg string) Output {
	o.Common = o.Common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o output) WithTerminalSequence(seq string) Output {
	o.Common = o.Common.WithTerminalSequence(seq)
	return o
}

func (o output) encodeInto(top, hso map[string]any) {
	event.ApplyCommon(top, o.Common)
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
func (o output) Encode() ([]byte, int, error) {
	return event.MarshalHookOutput(event.PreToolUse, o.encodeInto)
}

// Merge combines other into this PreToolUse output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	decision, reason := hookkit.MergeRankedString(
		string(o.decision), o.reason,
		string(b.decision), b.reason,
		hookkit.PermissionRankString,
	)
	updatedInput, w := hookkit.TakeLastMap("updatedInput", o.updatedInput, b.updatedInput)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:            mergedCommon,
		decision:          event.PermissionDecision(decision),
		reason:            reason,
		updatedInput:      updatedInput,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop() || o.decision == event.DecisionDeny
}
