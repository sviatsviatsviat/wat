package event

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// ExitBlockOutput is a response for TeammateIdle, TaskCreated, and TaskCompleted.
// Context and continue:false use JSON on SuccessExit. Block encodes BlockExit with
// a plain-text reason marked BodyOnStderr so run writes it to stderr (Claude
// ignores stdout on exit 2). Prefer Block to roll back or continue the action; use
// WithContinue(false) only when stopping the teammate entirely.
type ExitBlockOutput interface {
	hookkit.Output
	isExitBlockOutput()
	// WithAdditionalContext injects model context (JSON path only).
	WithAdditionalContext(text string) ExitBlockOutput
	// WithContinue sets whether Claude should continue the session.
	// Pass false to stop the teammate entirely (not the same as Block).
	WithContinue(v bool) ExitBlockOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) ExitBlockOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) ExitBlockOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) ExitBlockOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) ExitBlockOutput
}

type exitBlockOutput struct {
	Common
	eventName         string
	blockExit         bool
	reason            string
	additionalContext string
}

func (exitBlockOutput) isExitBlockOutput() {}

// ContextExitBlock returns an ExitBlockOutput with optional additional context.
func ContextExitBlock(eventName, text string) ExitBlockOutput {
	return exitBlockOutput{eventName: eventName, additionalContext: text}
}

// BlockExitBlock returns an ExitBlockOutput that encodes BlockExit with reason.
func BlockExitBlock(eventName, reason string) ExitBlockOutput {
	return exitBlockOutput{eventName: eventName, blockExit: true, reason: reason}
}

// IsZero reports whether this hook response is empty.
func (o exitBlockOutput) IsZero() bool {
	return o.Common.IsZero() && !o.blockExit && o.reason == "" && o.additionalContext == ""
}

// WithAdditionalContext injects model context (JSON path only).
func (o exitBlockOutput) WithAdditionalContext(text string) ExitBlockOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o exitBlockOutput) WithContinue(v bool) ExitBlockOutput {
	o.Common = o.Common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o exitBlockOutput) WithStopReason(reason string) ExitBlockOutput {
	o.Common = o.Common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o exitBlockOutput) WithSuppressOutput(v bool) ExitBlockOutput {
	o.Common = o.Common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o exitBlockOutput) WithSystemMessage(msg string) ExitBlockOutput {
	o.Common = o.Common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o exitBlockOutput) WithTerminalSequence(seq string) ExitBlockOutput {
	o.Common = o.Common.WithTerminalSequence(seq)
	return o
}

func (o exitBlockOutput) encodeInto(top, hso map[string]any) {
	ApplyCommon(top, o.Common)
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders JSON on SuccessExit, or a plain-text reason with BlockExit.
func (o exitBlockOutput) Encode() ([]byte, int, error) {
	if o.blockExit {
		return []byte(o.reason), BlockExit, nil
	}
	return MarshalHookOutput(o.eventName, o.encodeInto)
}

// BodyOnStderr reports whether Encode's body should be written to stderr.
// BlockExit reasons go to stderr; Claude ignores stdout on exit 2.
func (o exitBlockOutput) BodyOnStderr() bool {
	return o.blockExit
}

// Merge combines other into this ExitBlockOutput.
func (o exitBlockOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(exitBlockOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	blockExit := o.blockExit || b.blockExit
	reason, w := hookkit.TakeLastString("reason", o.reason, b.reason)
	if w != "" {
		warnings = append(warnings, w)
	}
	// Exit-2 block wins over JSON context; clear context when blocking.
	additionalContext := hookkit.JoinContextStrings(o.additionalContext, b.additionalContext)
	if blockExit {
		additionalContext = ""
		mergedCommon = Common{}
	}
	eventName := o.eventName
	if eventName == "" {
		eventName = b.eventName
	}
	return exitBlockOutput{
		Common:            mergedCommon,
		eventName:         eventName,
		blockExit:         blockExit,
		reason:            reason,
		additionalContext: additionalContext,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o exitBlockOutput) Stop() bool {
	return o.Common.Stop() || o.blockExit
}
