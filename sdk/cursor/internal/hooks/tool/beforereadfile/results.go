package beforereadfile

import "github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"

// Results is the hook-scoped response builder supplied to On* handlers by registration.
type Results interface {
	// Allow returns an allow verdict.
	Allow() event.PermissionOutput
	// Deny returns a deny verdict with a user-facing message. Cursor's beforeReadFile
	// schema accepts permission and user_message only; the message is emitted as
	// user_message with process exit 0 so the host applies the JSON permission
	// field instead of treating exit 2 stdout as the message body.
	Deny(userMessage string) event.PermissionOutput
	// SoftDeny returns the same encoding as Deny. Cursor's beforeReadFile schema
	// has no "ask" value; SoftDeny exists so authors do not reach for a prompt API.
	SoftDeny(userMessage string) event.PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() event.PermissionOutput
	isResults()
}

type results struct{ event.GateResults }

func (results) isResults() {}

// Deny returns a beforeReadFile deny with user_message and exit 0.
func (results) Deny(userMessage string) event.PermissionOutput {
	return event.GateResults{}.DenyUserMessage(userMessage)
}

// SoftDeny returns the same encoding as Deny for beforeReadFile.
func (results) SoftDeny(userMessage string) event.PermissionOutput {
	return event.GateResults{}.DenyUserMessage(userMessage)
}
