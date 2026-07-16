package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// WorktreeCreate is the WorktreeCreate hook event.
type WorktreeCreate struct {
	Envelope
}

// EventName returns the hook event name.
func (WorktreeCreate) EventName() string { return EventWorktreeCreate }

func init() {
	registerDecoder(EventWorktreeCreate, decodeAs[WorktreeCreate])
}

// WorktreeCreateOutput is the response for WorktreeCreate events.
// Construct via WorktreeCreateResults builders and With* methods.
// A nil value is a no-op.
type WorktreeCreateOutput interface {
	isWorktreeCreateOutput()
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) WorktreeCreateOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) WorktreeCreateOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) WorktreeCreateOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) WorktreeCreateOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) WorktreeCreateOutput
}

type worktreeCreateOutput struct {
	common
	worktreePath string
}

func (worktreeCreateOutput) isWorktreeCreateOutput() {}
func (o worktreeCreateOutput) isZero() bool {
	return o.common.isZero() && o.worktreePath == ""
}

// WithContinue sets whether Claude should continue the session.
func (o worktreeCreateOutput) WithContinue(v bool) WorktreeCreateOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o worktreeCreateOutput) WithStopReason(reason string) WorktreeCreateOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o worktreeCreateOutput) WithSuppressOutput(v bool) WorktreeCreateOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o worktreeCreateOutput) WithSystemMessage(msg string) WorktreeCreateOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o worktreeCreateOutput) WithTerminalSequence(seq string) WorktreeCreateOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// WorktreeCreateResults is the hook-scoped response builder supplied to Chain handlers by registration.
type WorktreeCreateResults interface {
	// Path returns a worktree-path result.
	Path(path string) WorktreeCreateOutput
	isWorktreeCreateResults()
}

type worktreeCreateResults struct{}

func (worktreeCreateResults) isWorktreeCreateResults() {}

// Path returns a worktree-path result.
func (worktreeCreateResults) Path(path string) WorktreeCreateOutput {
	return worktreeCreateOutput{worktreePath: path}
}

func (worktreeCreateOutput) allowedEvents() []string {
	return []string{EventWorktreeCreate}
}

func (o worktreeCreateOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.worktreePath != "" {
		hso["worktreePath"] = o.worktreePath
	}
}

// WorktreeCreate registers a WorktreeCreate handler.
func (c *Chain) WorktreeCreate(fn func(context.Context, Hook[WorktreeCreate], WorktreeCreateResults) (WorktreeCreateOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev WorktreeCreate) (WorktreeCreateOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), worktreeCreateResults{})
	})
	return &Chain{}
}
