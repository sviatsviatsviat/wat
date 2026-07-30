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
	// PermissionManual is a Claude Code alias for default (v2.1.200+).
	PermissionManual = event.PermissionManual
)

// PermissionUpdateType is the discriminator for a permission update entry.
type PermissionUpdateType = event.PermissionUpdateType

const (
	// PermissionUpdateAddRules adds permission rules.
	PermissionUpdateAddRules = event.PermissionUpdateAddRules
	// PermissionUpdateReplaceRules replaces rules of a given behavior at a destination.
	PermissionUpdateReplaceRules = event.PermissionUpdateReplaceRules
	// PermissionUpdateRemoveRules removes matching rules of a given behavior.
	PermissionUpdateRemoveRules = event.PermissionUpdateRemoveRules
	// PermissionUpdateSetMode changes the permission mode.
	PermissionUpdateSetMode = event.PermissionUpdateSetMode
	// PermissionUpdateAddDirectories adds working directories.
	PermissionUpdateAddDirectories = event.PermissionUpdateAddDirectories
	// PermissionUpdateRemoveDirectories removes working directories.
	PermissionUpdateRemoveDirectories = event.PermissionUpdateRemoveDirectories
)

// PermissionDestination controls where a permission update is written.
type PermissionDestination = event.PermissionDestination

const (
	// PermissionDestinationSession keeps the update in memory for the session.
	PermissionDestinationSession = event.PermissionDestinationSession
	// PermissionDestinationLocalSettings writes to .claude/settings.local.json.
	PermissionDestinationLocalSettings = event.PermissionDestinationLocalSettings
	// PermissionDestinationProjectSettings writes to .claude/settings.json.
	PermissionDestinationProjectSettings = event.PermissionDestinationProjectSettings
	// PermissionDestinationUserSettings writes to ~/.claude/settings.json.
	PermissionDestinationUserSettings = event.PermissionDestinationUserSettings
)

// PermissionRule is one tool rule inside an add/replace/removeRules update.
type PermissionRule = event.PermissionRule

// PermissionUpdate is one Claude permission update / suggestion entry.
//
// PermissionRequest input uses these as permission_suggestions. Allow
// responses may echo them via WithUpdatedPermissions.
type PermissionUpdate = event.PermissionUpdate

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

// BackgroundTask describes one in-flight background task on Stop and
// SubagentStop events (Claude Code v2.1.145+).
type BackgroundTask = event.BackgroundTask

// SessionCron describes one session-scoped scheduled wakeup on Stop and
// SubagentStop events (Claude Code v2.1.145+).
type SessionCron = event.SessionCron

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
// plus optional additionalContext (UserPromptExpansion, PreCompact, PostToolBatch,
// ConfigChange, and similar).
type DecisionOutput = event.DecisionOutput

// ExitBlockOutput is a response for TeammateIdle, TaskCreated, and TaskCompleted.
// Block uses Claude exit 2; continue:false stops the teammate entirely.
type ExitBlockOutput = event.ExitBlockOutput

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
