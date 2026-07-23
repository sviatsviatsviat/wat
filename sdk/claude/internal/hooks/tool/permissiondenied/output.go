package permissiondenied

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for this hook event.
// Construct via Results builders and With* methods.
// A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
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
	retry bool
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && !o.retry
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
	if o.retry {
		hso["retry"] = true
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	return event.MarshalHookOutput(event.PermissionDenied, o.encodeInto)
}

// Merge combines other into this PermissionDenied output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	return output{
		Common: mergedCommon,
		retry:  o.retry || b.retry,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop()
}
