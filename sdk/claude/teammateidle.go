package claude

// TeammateIdle is the TeammateIdle hook event.
type TeammateIdle struct {
	Envelope
}

// EventName returns the hook event name.
func (TeammateIdle) EventName() string { return EventTeammateIdle }

func init() {
	registerDecoder(EventTeammateIdle, decodeAs[TeammateIdle])
}
