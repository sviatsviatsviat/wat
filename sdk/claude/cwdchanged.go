package claude

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// CwdChanged is the CwdChanged hook event.
type CwdChanged struct {
	Envelope
	hookkit.RawPayload
	// NewCwd is the new working directory.
	NewCwd string `json:"new_cwd"`
	// OldCwd is the previous working directory.
	OldCwd string `json:"old_cwd,omitempty"`
}

// EventName returns the hook event name.
func (CwdChanged) EventName() string { return EventCwdChanged }

func init() {
	registerDecoder(EventCwdChanged, decodeAs[CwdChanged])
}
