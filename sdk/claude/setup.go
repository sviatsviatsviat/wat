package claude

// Setup is the Setup hook event.
type Setup struct {
	Envelope
	// Trigger is the setup trigger (init, maintenance).
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (Setup) EventName() string { return EventSetup }

func init() {
	registerDecoder(EventSetup, decodeAs[Setup])
}
