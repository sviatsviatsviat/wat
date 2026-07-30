package event

// TaskFields holds Claude task wire fields shared by TaskCreated and
// TaskCompleted events.
type TaskFields struct {
	// TaskID is the identifier of the task.
	TaskID string `json:"task_id"`
	// TaskSubject is the title of the task.
	TaskSubject string `json:"task_subject"`
	// TaskDescription is the detailed task description when provided.
	TaskDescription string `json:"task_description"`
	// TeammateName is the teammate associated with the task when provided.
	TeammateName string `json:"teammate_name"`
	// TeamName is the session-derived team name when provided (deprecated by Claude Code).
	TeamName string `json:"team_name"`
}
