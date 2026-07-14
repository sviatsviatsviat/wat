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

// SubagentStartResults is the hook-scoped response builder supplied to Chain handlers by registration.
type SubagentStartResults interface {
	Allow() PermissionOutput
	Deny(agentMessage string) PermissionOutput
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

func init() {
	registerDecoder(EventSubagentStart, decodeAs[SubagentStart])
}

// SubagentStart registers a subagentStart handler.
func (c *Chain) SubagentStart(fn func(context.Context, SubagentStartHook, SubagentStartResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStart) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return &Chain{}
}
