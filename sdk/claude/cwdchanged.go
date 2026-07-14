package claude

// CwdChanged is the CwdChanged hook event.
type CwdChanged struct {
	Envelope
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
