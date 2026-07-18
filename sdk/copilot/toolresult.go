package copilot

// ToolResult is a Copilot tool result object.
type ToolResult struct {
	// ResultType is the result type.
	ResultType string `json:"result_type"`
	// TextResultForLLM is the LLM-facing text.
	TextResultForLLM string `json:"text_result_for_llm"`
}

// Text returns the textual result.
func (r ToolResult) Text() string {
	return r.TextResultForLLM
}

// ErrorDetail is a structured Copilot error object.
type ErrorDetail struct {
	// Message is the error message.
	Message string `json:"message"`
	// Name is the error category name.
	Name string `json:"name"`
	// Stack is the error stack trace when provided.
	Stack string `json:"stack"`
}

// RawEvent holds an unknown or future hook event decoded from the shared envelope only.
type RawEvent struct {
	Envelope
}

// EventName returns the canonical or received hook event name.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	if e.receivedName != "" {
		return e.receivedName
	}
	return "unknown"
}
