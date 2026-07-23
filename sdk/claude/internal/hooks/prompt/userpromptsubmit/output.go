package userpromptsubmit

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
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) Output
	// WithSessionTitle sets the session title.
	WithSessionTitle(title string) Output
	// WithSuppressOriginalPrompt suppresses the original prompt when true.
	WithSuppressOriginalPrompt(v bool) Output
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
	block                  bool
	reason                 string
	additionalContext      string
	sessionTitle           string
	suppressOriginalPrompt bool
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && !o.block && o.reason == "" &&
		o.additionalContext == "" && o.sessionTitle == "" && !o.suppressOriginalPrompt
}

// WithAdditionalContext injects model context.
func (o output) WithAdditionalContext(text string) Output {
	o.additionalContext = text
	return o
}

// WithSessionTitle sets the session title.
func (o output) WithSessionTitle(title string) Output {
	o.sessionTitle = title
	return o
}

// WithSuppressOriginalPrompt suppresses the original prompt when true.
func (o output) WithSuppressOriginalPrompt(v bool) Output {
	o.suppressOriginalPrompt = v
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
	if o.block {
		top["decision"] = "block"
		if o.reason != "" {
			top["reason"] = o.reason
		}
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
	if o.sessionTitle != "" {
		hso["sessionTitle"] = o.sessionTitle
	}
	if o.suppressOriginalPrompt {
		hso["suppressOriginalPrompt"] = true
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	return event.MarshalHookOutput(event.UserPromptSubmit, o.encodeInto)
}

// Merge combines other into this UserPromptSubmit output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
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
	sessionTitle, w := hookkit.TakeLastString("sessionTitle", o.sessionTitle, b.sessionTitle)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:                 mergedCommon,
		block:                  dec == "block",
		reason:                 reason,
		additionalContext:      hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
		sessionTitle:           sessionTitle,
		suppressOriginalPrompt: o.suppressOriginalPrompt || b.suppressOriginalPrompt,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop() || o.block
}
