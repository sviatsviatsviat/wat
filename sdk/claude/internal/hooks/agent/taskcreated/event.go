package taskcreated

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the TaskCreated hook event.
type Event struct {
	event.Envelope
	// TaskID is the identifier of the task being created.
	TaskID string `json:"task_id"`
	// TaskSubject is the title of the task.
	TaskSubject string `json:"task_subject"`
	// TaskDescription is the detailed task description when provided.
	TaskDescription string `json:"task_description"`
	// TeammateName is the teammate creating the task when provided.
	TeammateName string `json:"teammate_name"`
	// TeamName is the session-derived team name when provided (deprecated by Claude Code).
	TeamName string `json:"team_name"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.TaskCreated }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.TaskCreated, hookkit.EventDecoder[Event](c))
}
