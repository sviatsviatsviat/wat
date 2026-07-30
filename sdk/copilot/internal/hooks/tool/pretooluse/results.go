package pretooluse

import "github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Allow returns an allow verdict.
	Allow() Output
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) Output
	// Ask returns an ask verdict with an agent-facing reason. Under Copilot
	// cloud agent, "ask" is treated as "deny" because no user is available.
	Ask(reason string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Allow returns an allow verdict.
func (results) Allow() Output {
	return output{decision: event.DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (results) Deny(reason string) Output {
	return output{decision: event.DecisionDeny, reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason. Cloud agent treats
// ask as deny.
func (results) Ask(reason string) Output {
	return output{decision: event.DecisionAsk, reason: reason}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
