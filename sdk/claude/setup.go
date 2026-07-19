package claude

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Setup is the Setup hook event.
type Setup struct {
	Envelope
	// Trigger is the setup trigger (init, maintenance).
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (Setup) EventName() string { return EventSetup }

func init() {
	codec.Register(EventSetup, hookkit.EventDecoder[Setup](codec))
}
