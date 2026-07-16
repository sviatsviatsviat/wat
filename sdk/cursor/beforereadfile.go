package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// BeforeReadFile is the beforeReadFile hook event.
type BeforeReadFile struct {
	Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
	// Attachments are additional file attachments.
	Attachments []Attachment `json:"attachments"`
}

// EventName returns the canonical hook event name.
func (BeforeReadFile) EventName() string { return EventBeforeReadFile }

// BeforeReadFileResults is the hook-scoped response builder supplied to Chain handlers by registration.
type BeforeReadFileResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) PermissionOutput
	// Ask returns an ask verdict with an agent-facing message.
	Ask(agentMessage string) PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PermissionOutput
	isBeforeReadFileResults()
}

type beforeReadFileResults struct{ permissionGateResults }

func (beforeReadFileResults) isBeforeReadFileResults() {}

// Allow returns an allow verdict.
func (r beforeReadFileResults) Allow() PermissionOutput { return r.allow() }

// Deny returns a deny verdict with an agent-facing message.
func (r beforeReadFileResults) Deny(agentMessage string) PermissionOutput {
	return r.deny(agentMessage)
}

// Ask returns an ask verdict with an agent-facing message.
func (r beforeReadFileResults) Ask(agentMessage string) PermissionOutput {
	return r.ask(agentMessage)
}

// Noop returns an empty response (silent stdout).
func (r beforeReadFileResults) Noop() PermissionOutput { return r.noop() }

func init() {
	registerDecoder(EventBeforeReadFile, decodeAs[BeforeReadFile])
}

// BeforeReadFile registers a beforeReadFile handler.
func (c *Chain) BeforeReadFile(fn func(context.Context, Hook[BeforeReadFile], BeforeReadFileResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeReadFile) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeReadFileResults{})
	})
	return &Chain{}
}
