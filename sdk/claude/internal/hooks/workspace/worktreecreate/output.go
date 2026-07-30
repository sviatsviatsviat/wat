package worktreecreate

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Output is the response for this hook event.
// Construct via Results builders and With* methods.
// A nil value is a no-op.
//
// Encode writes a plain worktree path for Claude Code command hooks (the
// transport wat installs). HTTP hooks instead expect JSON
// hookSpecificOutput.worktreePath; wat does not emit that form.
// Shared Common With* fields remain for merge consistency but are not written
// to stdout for this event.
type Output interface {
	hookkit.Output
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
	worktreePath string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && o.worktreePath == ""
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

// Encode renders this output as a plain worktree path on stdout for Claude Code
// command hooks. An empty path is a no-op encode (nil body, exit 0); the host
// treats a missing path as worktree-creation failure.
func (o output) Encode() ([]byte, int, error) {
	if o.worktreePath == "" {
		return nil, event.SuccessExit, nil
	}
	return []byte(o.worktreePath), event.SuccessExit, nil
}

// Merge combines other into this WorktreeCreate output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	worktreePath, w := hookkit.TakeLastString("worktreePath", o.worktreePath, b.worktreePath)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:       mergedCommon,
		worktreePath: worktreePath,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop()
}
