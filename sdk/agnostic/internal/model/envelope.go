package model

import "encoding/json"

// Envelope carries shared metadata present on every normalized hook event.
type Envelope struct {
	// Agent is the dialect that emitted this hook event (e.g. "claude").
	Agent string
	// Name is the native event name as received.
	Name string
	// Session holds session_id, sessionId, or conversation_id from the native payload.
	Session string
	// Cwd is the working directory from the native payload.
	Cwd string
	// TranscriptPath is the conversation transcript path when provided.
	TranscriptPath string
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}
