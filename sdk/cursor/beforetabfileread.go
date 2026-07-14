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

// BeforeTabFileReadResults is the hook-scoped response builder supplied to Chain handlers by registration.
type BeforeTabFileReadResults interface {
	Allow() PermissionOutput
	Deny(agentMessage string) PermissionOutput
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

func init() {
	registerDecoder(EventBeforeTabFileRead, decodeAs[BeforeTabFileRead])
}

// BeforeTabFileRead registers a beforeTabFileRead handler.
func (c *Chain) BeforeTabFileRead(fn func(context.Context, BeforeTabFileReadHook, BeforeTabFileReadResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeTabFileRead) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeTabFileReadResults{})
	})
	return &Chain{}
}
