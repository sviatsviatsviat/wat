package beforesubmitprompt

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Block blocks prompt submission with a user-facing message.
	Block(userMessage string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Block blocks prompt submission with a user-facing message.
func (results) Block(userMessage string) Output {
	cont := false
	return output{cont: &cont, userMessage: userMessage}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
