package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// BeforeTabFileRead is the beforeTabFileRead hook event.
type BeforeTabFileRead struct {
	Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
}

// EventName returns the canonical hook event name.
func (BeforeTabFileRead) EventName() string { return EventBeforeTabFileRead }

// BeforeTabFileReadResults is the hook-scoped response builder supplied to On* handlers by registration.
type BeforeTabFileReadResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) PermissionOutput
	// Ask returns an ask verdict with an agent-facing message.
	Ask(agentMessage string) PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PermissionOutput
	isBeforeTabFileReadResults()
}

type beforeTabFileReadResults struct{ permissionGateResults }

func (beforeTabFileReadResults) isBeforeTabFileReadResults() {}

// Allow returns an allow verdict.
func (r beforeTabFileReadResults) Allow() PermissionOutput { return r.allow() }

// Deny returns a deny verdict with an agent-facing message.
func (r beforeTabFileReadResults) Deny(agentMessage string) PermissionOutput {
	return r.deny(agentMessage)
}

// Ask returns an ask verdict with an agent-facing message.
func (r beforeTabFileReadResults) Ask(agentMessage string) PermissionOutput {
	return r.ask(agentMessage)
}

// Noop returns an empty response (silent stdout).
func (r beforeTabFileReadResults) Noop() PermissionOutput { return r.noop() }

func init() {
	registerDecoder(EventBeforeTabFileRead, decodeAs[BeforeTabFileRead])
}

// OnBeforeTabFileRead registers a beforeTabFileRead handler.
func OnBeforeTabFileRead(fn func(context.Context, Hook[BeforeTabFileRead], BeforeTabFileReadResults) (PermissionOutput, error)) *chain {
	return (&chain{}).BeforeTabFileRead(fn)
}

// BeforeTabFileRead registers another BeforeTabFileRead handler on the chain.
func (c *chain) BeforeTabFileRead(fn func(context.Context, Hook[BeforeTabFileRead], BeforeTabFileReadResults) (PermissionOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeTabFileRead) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeTabFileReadResults{})
	})
	return c
}
