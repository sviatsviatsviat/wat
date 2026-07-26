package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Envelope holds fields shared by every Cursor hook event payload.
type Envelope = event.Envelope

// ModelParam is a selected model parameter from Cursor's common hook schema.
type ModelParam = event.ModelParam

// Attachment is a file or rule attachment on prompt and read-file events.
type Attachment = event.Attachment

// Edit is a single file edit on afterFileEdit and afterTabFileEdit events.
type Edit = event.Edit

// PermissionDecision is a permission verdict label on permission-gating events.
type PermissionDecision = event.PermissionDecision

const (
	// DecisionAllow permits the action.
	DecisionAllow = event.DecisionAllow
	// DecisionDeny blocks the action.
	DecisionDeny = event.DecisionDeny
	// DecisionAsk escalates to the user.
	DecisionAsk = event.DecisionAsk
)

// PermissionOutput is the response for permission-gating events.
// Construct via PermissionResults builders and With* methods. A nil value is a no-op.
type PermissionOutput = event.PermissionOutput

// PermissionResults is the hook-scoped response builder supplied to permission On* handlers by registration.
type PermissionResults = event.PermissionResults

// PostToolOutput is the response for post-tool events.
// Construct via PostToolResults builders and With* methods. A nil value is a no-op.
type PostToolOutput = event.PostToolOutput

// PostToolResults is the hook-scoped response builder supplied to On* handlers by registration.
type PostToolResults = event.PostToolResults

// Canonical Cursor hook event names for config keys and mux dispatch.
const (
	EventSessionStart         = event.SessionStart
	EventSessionEnd           = event.SessionEnd
	EventBeforeSubmitPrompt   = event.BeforeSubmitPrompt
	EventPreToolUse           = event.PreToolUse
	EventPostToolUse          = event.PostToolUse
	EventPostToolUseFailure   = event.PostToolUseFailure
	EventBeforeShellExecution = event.BeforeShellExecution
	EventAfterShellExecution  = event.AfterShellExecution
	EventBeforeMCPExecution   = event.BeforeMCPExecution
	EventAfterMCPExecution    = event.AfterMCPExecution
	EventBeforeReadFile       = event.BeforeReadFile
	EventAfterFileEdit        = event.AfterFileEdit
	EventSubagentStart        = event.SubagentStart
	EventSubagentStop         = event.SubagentStop
	EventStop                 = event.Stop
	EventPreCompact           = event.PreCompact
	EventAfterAgentResponse   = event.AfterAgentResponse
	EventAfterAgentThought    = event.AfterAgentThought
	EventBeforeTabFileRead    = event.BeforeTabFileRead
	EventAfterTabFileEdit     = event.AfterTabFileEdit
	EventWorkspaceOpen        = event.WorkspaceOpen
)
