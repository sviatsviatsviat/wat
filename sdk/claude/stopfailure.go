package claude

// StopFailure is the StopFailure hook event.
type StopFailure struct {
	Envelope
	// ErrorType is the error category (rate_limit, overloaded, …).
	ErrorType string `json:"error_type"`
	// Message is the error message when provided.
	Message string `json:"message"`
}

// EventName returns the hook event name.
func (StopFailure) EventName() string { return EventStopFailure }

func init() {
	registerDecoder(EventStopFailure, decodeAs[StopFailure])
}
