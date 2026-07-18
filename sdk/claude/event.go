package claude

import "encoding/json"

// Event is implemented by every decoded Claude Code hook event.
type Event interface {
	EventName() string
	Raw() json.RawMessage
	envelope() Envelope
}
