package event

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// DecisionOutput is a response for Claude events that accept top-level
// decision:"block" plus optional additionalContext and shared Common fields.
// Construct via ContextDecision / BlockDecision and With* methods. A nil value
// is a no-op. Encode always uses SuccessExit; Claude only processes JSON on
// exit 0.
type DecisionOutput interface {
	hookkit.Output
	isDecisionOutput()
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) DecisionOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) DecisionOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) DecisionOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) DecisionOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) DecisionOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) DecisionOutput
}

type decisionOutput struct {
	Common
	eventName         string
	block             bool
	reason            string
	additionalContext string
}

func (decisionOutput) isDecisionOutput() {}

// ContextDecision returns a DecisionOutput with optional additional context.
func ContextDecision(eventName, text string) DecisionOutput {
	return decisionOutput{eventName: eventName, additionalContext: text}
}

// BlockDecision returns a DecisionOutput with decision:"block" and reason.
func BlockDecision(eventName, reason string) DecisionOutput {
	return decisionOutput{eventName: eventName, block: true, reason: reason}
}

// IsZero reports whether this hook response is empty.
func (o decisionOutput) IsZero() bool {
	return o.Common.IsZero() && !o.block && o.reason == "" && o.additionalContext == ""
}

// WithAdditionalContext injects model context.
func (o decisionOutput) WithAdditionalContext(text string) DecisionOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o decisionOutput) WithContinue(v bool) DecisionOutput {
	o.Common = o.Common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o decisionOutput) WithStopReason(reason string) DecisionOutput {
	o.Common = o.Common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o decisionOutput) WithSuppressOutput(v bool) DecisionOutput {
	o.Common = o.Common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o decisionOutput) WithSystemMessage(msg string) DecisionOutput {
	o.Common = o.Common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o decisionOutput) WithTerminalSequence(seq string) DecisionOutput {
	o.Common = o.Common.WithTerminalSequence(seq)
	return o
}

func (o decisionOutput) encodeInto(top, hso map[string]any) {
	ApplyCommon(top, o.Common)
	if o.block {
		top["decision"] = "block"
		if o.reason != "" {
			top["reason"] = o.reason
		}
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders this output as Claude Code stdout JSON with SuccessExit.
func (o decisionOutput) Encode() ([]byte, int, error) {
	return MarshalHookOutput(o.eventName, o.encodeInto)
}

// Merge combines other into this DecisionOutput.
func (o decisionOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(decisionOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	oDec, bDec := "", ""
	if o.block {
		oDec = "block"
	}
	if b.block {
		bDec = "block"
	}
	dec, reason := hookkit.MergeRankedString(oDec, o.reason, bDec, b.reason, hookkit.BlockDecisionRankString)
	eventName := o.eventName
	if eventName == "" {
		eventName = b.eventName
	}
	return decisionOutput{
		Common:            mergedCommon,
		eventName:         eventName,
		block:             dec == "block",
		reason:            reason,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o decisionOutput) Stop() bool {
	return o.Common.Stop() || o.block
}
