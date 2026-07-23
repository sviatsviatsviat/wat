package teammateidle

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the TeammateIdle hook event.
type Event struct {
	event.Envelope
	// TeammateName is the idle teammate identity.
	TeammateName string `json:"teammate_name"`
	// TeamName is the agent team name when provided.
	TeamName string `json:"team_name"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.TeammateIdle }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.TeammateIdle, hookkit.EventDecoder[Event](c))
}
