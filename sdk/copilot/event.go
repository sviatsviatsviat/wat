package copilot

// Event is implemented by every decoded GitHub Copilot hook event.
type Event interface {
	EventName() string
	envelope() Envelope
}
