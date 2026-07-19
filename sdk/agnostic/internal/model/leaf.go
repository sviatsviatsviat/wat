package model

// Subagent describes subagent lifecycle events.
type Subagent struct {
	// ID is the subagent identifier.
	ID string
	// Type is the subagent type.
	Type string
	// Task is the subagent task description.
	Task string
	// Summary is the subagent result summary.
	Summary string
	// Status is the subagent status (completed, error, aborted, end_turn, ...).
	Status string
	// TranscriptPath is the subagent transcript path when provided.
	TranscriptPath string
	// LoopCount is the subagent loop count when available.
	LoopCount int
}

// TurnEnd describes agent stop events.
type TurnEnd struct {
	// Status is the turn-end status when provided.
	Status string
	// LoopCount is the stop loop count when available.
	LoopCount int
	// StopHookActive is true when a stop hook is already active (Claude).
	StopHookActive bool
	// LastAssistantMessage is the last assistant message when provided.
	LastAssistantMessage string
}

// CompactInfo describes context compaction events.
type CompactInfo struct {
	// Trigger is the compaction trigger reason.
	Trigger string
	// CustomInstructions holds extra compaction instructions when provided.
	CustomInstructions string
}

// Lifecycle describes session start and end events.
type Lifecycle struct {
	// Source is the session start source when provided.
	Source string
	// Reason is the session end reason when provided.
	Reason string
	// Model is the model name when provided.
	Model string
	// InitialPrompt is the session start prompt when provided.
	InitialPrompt string
	// Background is true for background agent sessions when provided.
	Background bool
}
