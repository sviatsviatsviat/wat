package agnostic

import "encoding/json"

// Envelope carries shared metadata present on every normalized hook event.
type Envelope struct {
	// Agent is the dialect that emitted this hook event.
	Agent Dialect
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

func envelopeFrom(ev *Event) Envelope {
	if ev == nil {
		return Envelope{}
	}
	return Envelope{
		Agent:          ev.Agent,
		Name:           ev.Name,
		Session:        ev.Session,
		Cwd:            ev.Cwd,
		TranscriptPath: ev.TranscriptPath,
		Raw:            ev.Raw,
	}
}
