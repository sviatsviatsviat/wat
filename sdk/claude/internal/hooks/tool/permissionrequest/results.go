package permissionrequest

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Allow returns an allow verdict.
	Allow() Output
	// Deny returns a deny verdict with a permission message.
	Deny(message string) Output
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
