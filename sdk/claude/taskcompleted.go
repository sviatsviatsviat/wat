package claude

import (
	"encoding/json"
)

// TaskCompleted is the TaskCompleted hook event.
type TaskCompleted struct {
	Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (TaskCompleted) EventName() string { return EventTaskCompleted }

func init() {
	registerDecoder(EventTaskCompleted, decodeAs[TaskCompleted])
}
