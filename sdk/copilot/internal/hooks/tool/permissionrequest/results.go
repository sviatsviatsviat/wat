package permissionrequest

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Allow returns an allow verdict.
	Allow() Output
	// Deny returns a deny verdict with a permission message.
	Deny(message string) Output
	// Ask returns an ask-style deny that suppresses warn-exit semantics.
	Ask(message string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Allow returns an allow verdict.
func (results) Allow() Output {
	return output{behavior: "allow"}
}

// Deny returns a deny verdict with a permission message.
func (results) Deny(message string) Output {
	return output{behavior: "deny", message: message}
}

// Ask returns an ask-style deny that suppresses warn-exit semantics.
func (results) Ask(message string) Output {
	return output{behavior: "deny", message: message, suppressWarnExit: true}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
