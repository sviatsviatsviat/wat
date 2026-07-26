package beforesubmitprompt

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Block blocks prompt submission with a user-facing message. Encoding writes
	// continue:false and user_message with process exit 0. Cursor's
	// beforeSubmitPrompt control channel is the JSON continue field; do not use
	// exit 2 for this event.
	Block(userMessage string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Block blocks prompt submission with continue:false, user_message, and exit 0.
func (results) Block(userMessage string) Output {
	cont := false
	return output{cont: &cont, userMessage: userMessage}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
