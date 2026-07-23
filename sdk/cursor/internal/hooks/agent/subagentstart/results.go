package subagentstart

import "github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"

// Results is the hook-scoped response builder supplied to On* handlers by registration.
type Results interface {
	// Allow returns an allow verdict.
	Allow() event.PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) event.PermissionOutput
	// Ask returns an ask verdict with an agent-facing message.
	Ask(agentMessage string) event.PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() event.PermissionOutput
	isResults()
}

type results struct{ event.GateResults }

func (results) isResults() {}
