package claude

// Event is implemented by every decoded Claude Code hook event.
type Event interface {
	EventName() string
	envelope() Envelope
}
