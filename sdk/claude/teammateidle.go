package claude

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// TeammateIdle is the TeammateIdle hook event.
type TeammateIdle struct {
	Envelope
	hookkit.RawPayload
	// TeammateName is the idle teammate identity.
	TeammateName string `json:"teammate_name"`
	// TeamName is the agent team name when provided.
	TeamName string `json:"team_name"`
}

// EventName returns the hook event name.
func (TeammateIdle) EventName() string { return EventTeammateIdle }

func init() {
	registerDecoder(EventTeammateIdle, decodeAs[TeammateIdle])
}
