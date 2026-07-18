package copilot

import "encoding/json"

// Event is implemented by every decoded GitHub Copilot hook event.
type Event interface {
	EventName() string
	Raw() json.RawMessage
	envelope() Envelope
}
