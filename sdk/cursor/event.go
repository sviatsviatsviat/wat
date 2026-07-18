package cursor

import "encoding/json"

// Event is implemented by every decoded Cursor hook event.
type Event interface {
	EventName() string
	Raw() json.RawMessage
	envelope() Envelope
}
