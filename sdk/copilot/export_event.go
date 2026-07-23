package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Timestamp is an RFC3339 timestamp on Copilot hook payloads.
type Timestamp = event.Timestamp

// Envelope holds fields shared by every GitHub Copilot hook event payload.
type Envelope = event.Envelope

// PermissionDecision is a pre-tool permission verdict label.
type PermissionDecision = event.PermissionDecision

const (
	// DecisionAllow permits the tool call.
	DecisionAllow = event.DecisionAllow
	// DecisionDeny blocks the tool call.
	DecisionDeny = event.DecisionDeny
	// DecisionAsk escalates to the user.
	DecisionAsk = event.DecisionAsk
)

// ToolResult is a Copilot tool result object.
type ToolResult = event.ToolResult

// ErrorDetail is a structured Copilot error object.
type ErrorDetail = event.ErrorDetail

// Canonical PascalCase GitHub Copilot hook event names for config keys and mux dispatch.
const (
	EventSessionStart        = event.SessionStart
	EventSessionEnd          = event.SessionEnd
	EventUserPromptSubmitted = event.UserPromptSubmitted
	EventPreToolUse          = event.PreToolUse
	EventPostToolUse         = event.PostToolUse
	EventPostToolUseFailure  = event.PostToolUseFailure
	EventPermissionRequest   = event.PermissionRequest
	EventSubagentStart       = event.SubagentStart
	EventSubagentStop        = event.SubagentStop
	EventAgentStop           = event.AgentStop
	EventPreCompact          = event.PreCompact
	EventNotification        = event.Notification
	EventErrorOccurred       = event.ErrorOccurred
)
