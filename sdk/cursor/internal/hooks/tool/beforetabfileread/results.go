package beforetabfileread

import "github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"

// Results is the hook-scoped response builder supplied to On* handlers by registration.
// Cursor's beforeTabFileRead schema accepts permission allow|deny only: there is no
// ask verdict and no user_message or agent_message field.
type Results interface {
	// Allow returns an allow verdict encoded as permission only.
	Allow() event.PermissionOutput
	// Deny returns a deny verdict encoded as permission only with process exit 0.
	Deny() event.PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() event.PermissionOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Allow returns a beforeTabFileRead allow with permission only.
func (results) Allow() event.PermissionOutput {
	return event.GateResults{}.PermissionOnlyAllow()
}

// Deny returns a beforeTabFileRead deny with permission only and exit 0.
func (results) Deny() event.PermissionOutput {
	return event.GateResults{}.PermissionOnlyDeny()
}

// Noop returns an empty response (silent stdout).
func (results) Noop() event.PermissionOutput {
	return event.GateResults{}.Noop()
}
