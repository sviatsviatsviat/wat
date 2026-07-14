package model

// Codec decodes a native payload into the unified Event and encodes a unified
// Result back into the native stdout JSON and process exit code.
type Codec interface {
	Dialect() Dialect
	// Decode parses a native payload. eventHint carries the native event name
	// when it is known from configuration context; it is required for Copilot
	// camelCase payloads, which do not name their event, and ignored when the
	// payload itself carries hook_event_name.
	Decode(raw []byte, eventHint string) (*Event, error)
	// Encode renders res for ev. Empty stdout with exit 0 means "no decision".
	Encode(ev *Event, res Result) (stdout []byte, exitCode int, err error)
}
