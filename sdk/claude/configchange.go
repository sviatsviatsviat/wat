package claude

// ConfigChange is the ConfigChange hook event.
type ConfigChange struct {
	Envelope
	// Source is the config source that changed.
	Source string `json:"source"`
}

// EventName returns the hook event name.
func (ConfigChange) EventName() string { return EventConfigChange }

func init() {
	registerDecoder(EventConfigChange, decodeAs[ConfigChange])
}
