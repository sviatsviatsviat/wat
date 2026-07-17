package agnostic

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// Kind is a normalized hook event category shared by all supported agents.
type Kind = model.Kind

const (
	// KindSessionStart is the normalized category for session start events.
	KindSessionStart = model.KindSessionStart
	// KindSessionEnd is the normalized category for session end events.
	KindSessionEnd = model.KindSessionEnd
	// KindUserPrompt is the normalized category for user prompt submission events.
	KindUserPrompt = model.KindUserPrompt
	// KindPreTool is the normalized category for pre-tool hook events.
	KindPreTool = model.KindPreTool
	// KindPostTool is the normalized category for successful post-tool events.
	KindPostTool = model.KindPostTool
	// KindPostToolFailure is the normalized category for failed post-tool events.
	KindPostToolFailure = model.KindPostToolFailure
	// KindPermissionRequest is the normalized category for permission request events.
	KindPermissionRequest = model.KindPermissionRequest
	// KindSubagentStart is the normalized category for subagent start events.
	KindSubagentStart = model.KindSubagentStart
	// KindSubagentStop is the normalized category for subagent stop events.
	KindSubagentStop = model.KindSubagentStop
	// KindStop is the normalized category for agent stop events.
	KindStop = model.KindStop
	// KindPreCompact is the normalized category for pre-compaction events.
	KindPreCompact = model.KindPreCompact
	// KindNotification is the normalized category for notification events.
	KindNotification = model.KindNotification
	// KindAgentError is the normalized category for agent runtime error events.
	KindAgentError = model.KindAgentError
	// KindOther is the normalized category for events with no dedicated mapping.
	KindOther = model.KindOther
)

// Event is the unified, agent-independent view of a hook invocation.
// Raw always carries the untouched native payload, so nothing is lost by
// normalization; agent-specific handlers can re-decode it with native types.
type Event = model.Event

// ToolCall describes the tool invocation a pre/post tool event refers to.
type ToolCall = model.ToolCall

// ToolResult describes the outcome of a completed or failed tool call.
type ToolResult = model.ToolResult
