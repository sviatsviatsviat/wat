package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// PermissionMode is the Claude Code permission mode on hook events.
type PermissionMode = event.PermissionMode

const (
	// PermissionDefault is the default permission mode.
	PermissionDefault = event.PermissionDefault
	// PermissionPlan is the plan permission mode.
	PermissionPlan = event.PermissionPlan
	// PermissionAcceptEdits auto-accepts edit tools.
	PermissionAcceptEdits = event.PermissionAcceptEdits
	// PermissionAuto is the auto permission mode.
	PermissionAuto = event.PermissionAuto
	// PermissionDontAsk suppresses permission prompts.
	PermissionDontAsk = event.PermissionDontAsk
	// PermissionBypassPermissions bypasses permission checks.
	PermissionBypassPermissions = event.PermissionBypassPermissions
)

// EffortLevel is the effort level on hook events.
type EffortLevel = event.EffortLevel

const (
	// EffortLow is the low effort level.
	EffortLow = event.EffortLow
	// EffortMedium is the medium effort level.
	EffortMedium = event.EffortMedium
	// EffortHigh is the high effort level.
	EffortHigh = event.EffortHigh
	// EffortXHigh is the extra-high effort level.
	EffortXHigh = event.EffortXHigh
	// EffortMax is the maximum effort level.
	EffortMax = event.EffortMax
)

// Effort carries effort metadata on hook events.
type Effort = event.Effort

// Envelope holds fields shared by every Claude Code hook event payload.
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
	// DecisionDefer defers the permission decision.
	DecisionDefer = event.DecisionDefer
)

// CommonOutput is a shared-fields-only response for events that only accept those fields.
type CommonOutput = event.CommonOutput

// DecisionOutput is a response for Claude events that accept decision:"block"
// plus optional additionalContext (PostToolBatch and similar).
type DecisionOutput = event.DecisionOutput

// Claude Code hook_event_name values for config keys and stdin payloads.
const (
	EventSessionStart        = event.SessionStart
	EventSetup               = event.Setup
	EventSessionEnd          = event.SessionEnd
	EventUserPromptSubmit    = event.UserPromptSubmit
	EventUserPromptExpansion = event.UserPromptExpansion
	EventPreToolUse          = event.PreToolUse
	EventPostToolUse         = event.PostToolUse
	EventPostToolUseFailure  = event.PostToolUseFailure
	EventPostToolBatch       = event.PostToolBatch
	EventPermissionRequest   = event.PermissionRequest
	EventPermissionDenied    = event.PermissionDenied
	EventSubagentStart       = event.SubagentStart
	EventSubagentStop        = event.SubagentStop
	EventTaskCreated         = event.TaskCreated
	EventTaskCompleted       = event.TaskCompleted
	EventStop                = event.Stop
	EventStopFailure         = event.StopFailure
	EventTeammateIdle        = event.TeammateIdle
	EventNotification        = event.Notification
	EventMessageDisplay      = event.MessageDisplay
	EventInstructionsLoaded  = event.InstructionsLoaded
	EventConfigChange        = event.ConfigChange
	EventCwdChanged          = event.CwdChanged
	EventFileChanged         = event.FileChanged
	EventWorktreeCreate      = event.WorktreeCreate
	EventWorktreeRemove      = event.WorktreeRemove
	EventPreCompact          = event.PreCompact
	EventPostCompact         = event.PostCompact
	EventElicitation         = event.Elicitation
	EventElicitationResult   = event.ElicitationResult
)
