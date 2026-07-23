package event

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// CommonOutput is a shared-fields-only response for events that only accept those fields.
// Construct via Results builders and With* methods. A nil value is a no-op.
type CommonOutput interface {
	run.Output
	isCommonOutput()
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) CommonOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) CommonOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) CommonOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) CommonOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) CommonOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) CommonOutput
}

type commonOutput struct {
	Common
	eventName         string
	additionalContext string
}

func (commonOutput) isCommonOutput() {}

// ContextOutput returns a CommonOutput with eventName and optional additional context.
func ContextOutput(eventName, text string) CommonOutput {
	return commonOutput{eventName: eventName, additionalContext: text}
}

// IsZero reports whether this hook response is empty.
func (o commonOutput) IsZero() bool {
	return o.Common.IsZero() && o.additionalContext == ""
}

// WithAdditionalContext injects model context.
func (o commonOutput) WithAdditionalContext(text string) CommonOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o commonOutput) WithContinue(v bool) CommonOutput {
	o.Common = o.Common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o commonOutput) WithStopReason(reason string) CommonOutput {
	o.Common = o.Common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o commonOutput) WithSuppressOutput(v bool) CommonOutput {
	o.Common = o.Common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o commonOutput) WithSystemMessage(msg string) CommonOutput {
	o.Common = o.Common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o commonOutput) WithTerminalSequence(seq string) CommonOutput {
	o.Common = o.Common.WithTerminalSequence(seq)
	return o
}

func (o commonOutput) encodeInto(top, hso map[string]any) {
	ApplyCommon(top, o.Common)
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o commonOutput) Encode() ([]byte, int, error) {
	return MarshalHookOutput(o.eventName, o.encodeInto)
}

// Merge combines other into this CommonOutput.
func (o commonOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(commonOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	eventName := o.eventName
	if eventName == "" {
		eventName = b.eventName
	}
	return commonOutput{
		Common:            mergedCommon,
		eventName:         eventName,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o commonOutput) Stop() bool {
	return o.Common.Stop()
}
