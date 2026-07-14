package cursor

// Canonical Cursor hook event names for config keys and mux dispatch.
const (
	EventSessionStart         = "sessionStart"
	EventSessionEnd           = "sessionEnd"
	EventBeforeSubmitPrompt   = "beforeSubmitPrompt"
	EventPreToolUse           = "preToolUse"
	EventPostToolUse          = "postToolUse"
	EventPostToolUseFailure   = "postToolUseFailure"
	EventBeforeShellExecution = "beforeShellExecution"
	EventAfterShellExecution  = "afterShellExecution"
	EventBeforeMCPExecution   = "beforeMCPExecution"
	EventAfterMCPExecution    = "afterMCPExecution"
	EventBeforeReadFile       = "beforeReadFile"
	EventAfterFileEdit        = "afterFileEdit"
	EventSubagentStart        = "subagentStart"
	EventSubagentStop         = "subagentStop"
	EventStop                 = "stop"
	EventPreCompact           = "preCompact"
	EventAfterAgentResponse   = "afterAgentResponse"
	EventAfterAgentThought    = "afterAgentThought"
	EventBeforeTabFileRead    = "beforeTabFileRead"
	EventAfterTabFileEdit     = "afterTabFileEdit"
	EventWorkspaceOpen        = "workspaceOpen"
)
