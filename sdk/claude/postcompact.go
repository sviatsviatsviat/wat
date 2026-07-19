package claude

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// PostCompact is the PostCompact hook event.
type PostCompact struct {
	Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (PostCompact) EventName() string { return EventPostCompact }

func init() {
	codec.Register(EventPostCompact, hookkit.EventDecoder[PostCompact](codec))
}
