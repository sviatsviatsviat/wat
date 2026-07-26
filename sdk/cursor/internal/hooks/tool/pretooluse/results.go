package pretooluse

import "github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"

// Results is the hook-scoped response builder supplied to PreToolUse handlers.
type Results interface {
	// Allow returns an allow verdict.
	Allow() event.PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) event.PermissionOutput
	// Ask encodes permission "ask" with an agent-facing message. Cursor's
	// preToolUse schema accepts "ask" but does not enforce it today, so the
	// host will not escalate to the user. Prefer Allow or Deny.
	Ask(agentMessage string) event.PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() event.PermissionOutput
	isResults()
}

type results struct{ event.GateResults }

func (results) isResults() {}

// Ask encodes permission ask; Cursor does not enforce it for preToolUse today.
func (results) Ask(agentMessage string) event.PermissionOutput {
	return event.GateResults{}.Ask(agentMessage)
}
