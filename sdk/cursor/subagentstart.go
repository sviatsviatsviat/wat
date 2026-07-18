package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStart is the subagentStart hook event.
type SubagentStart struct {
	Envelope
	// SubagentID is the subagent identifier.
	SubagentID string `json:"subagent_id"`
	// SubagentType is the subagent type.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
}

// EventName returns the canonical hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

// SubagentStartResults is the hook-scoped response builder supplied to On* handlers by registration.
type SubagentStartResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) PermissionOutput
	// Ask returns an ask verdict with an agent-facing message.
	Ask(agentMessage string) PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PermissionOutput
	isSubagentStartResults()
}

type subagentStartResults struct{ permissionGateResults }

func (subagentStartResults) isSubagentStartResults() {}

// Allow returns an allow verdict.
func (r subagentStartResults) Allow() PermissionOutput { return r.allow() }

// Deny returns a deny verdict with an agent-facing message.
func (r subagentStartResults) Deny(agentMessage string) PermissionOutput {
	return r.deny(agentMessage)
}

// Ask returns an ask verdict with an agent-facing message.
func (r subagentStartResults) Ask(agentMessage string) PermissionOutput {
	return r.ask(agentMessage)
}

// Noop returns an empty response (silent stdout).
func (r subagentStartResults) Noop() PermissionOutput { return r.noop() }

func init() {
	registerDecoder(EventSubagentStart, decodeAs[SubagentStart])
}

// OnSubagentStart registers a subagentStart handler.
func OnSubagentStart(fn func(context.Context, Hook[SubagentStart], SubagentStartResults) (PermissionOutput, error)) *chain {
	return (&chain{}).SubagentStart(fn)
}

// SubagentStart registers another SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, Hook[SubagentStart], SubagentStartResults) (PermissionOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStart) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return c
}
