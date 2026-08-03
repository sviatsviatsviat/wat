package permissionrequest

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Allow returns an allow verdict that short-circuits the host permission flow.
	Allow() Output
	// Deny returns a hard deny with a permission message, WarnExit (2), and
	// Stop so later handlers are skipped. Prefer Deny when the tool must not run.
	Deny(message string) Output
	// SoftDeny returns a soft deny: wire behavior "deny" with exit 0 (no WarnExit)
	// and without Stop. Copilot's permissionRequest schema has no "ask" value;
	// SoftDeny does not open a user confirmation UI. Prefer Deny to block, or
	// Noop to fall through to the host permission service (including user prompting).
	SoftDeny(message string) Output
	// Noop returns empty stdout so the host continues its normal permission
	// flow (rules, session approvals, and user prompting). Prefer nil from
	// handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Allow returns an allow verdict that short-circuits the host permission flow.
func (results) Allow() Output {
	return output{behavior: "allow"}
}

// Deny returns a hard deny with a permission message, WarnExit (2), and Stop.
func (results) Deny(message string) Output {
	return output{behavior: "deny", message: message}
}

// SoftDeny returns a soft deny (behavior "deny", exit 0). It does not escalate
// to the user; prefer Deny or Noop.
func (results) SoftDeny(message string) Output {
	return output{behavior: "deny", message: message, suppressWarnExit: true}
}

// Noop returns an empty response so the host permission flow continues.
func (results) Noop() Output {
	return output{}
}
