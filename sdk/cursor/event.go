package cursor

// Event is implemented by every decoded Cursor hook event.
type Event interface {
	EventName() string
	envelope() Envelope
}
