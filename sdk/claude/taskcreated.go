package claude

import (
	"encoding/json"
)

// TaskCreated is the TaskCreated hook event.
type TaskCreated struct {
	Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (TaskCreated) EventName() string { return EventTaskCreated }

func init() {
	registerDecoder(EventTaskCreated, decodeAs[TaskCreated])
}
