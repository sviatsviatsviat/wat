package event

// Canonical Cursor hook event names for config keys and mux dispatch.
const (
	SessionStart         = "sessionStart"
	SessionEnd           = "sessionEnd"
	BeforeSubmitPrompt   = "beforeSubmitPrompt"
	PreToolUse           = "preToolUse"
	PostToolUse          = "postToolUse"
	PostToolUseFailure   = "postToolUseFailure"
	BeforeShellExecution = "beforeShellExecution"
	AfterShellExecution  = "afterShellExecution"
	BeforeMCPExecution   = "beforeMCPExecution"
	AfterMCPExecution    = "afterMCPExecution"
	BeforeReadFile       = "beforeReadFile"
	AfterFileEdit        = "afterFileEdit"
	SubagentStart        = "subagentStart"
	SubagentStop         = "subagentStop"
	Stop                 = "stop"
	PreCompact           = "preCompact"
	AfterAgentResponse   = "afterAgentResponse"
	AfterAgentThought    = "afterAgentThought"
	BeforeTabFileRead    = "beforeTabFileRead"
	AfterTabFileEdit     = "afterTabFileEdit"
	WorkspaceOpen        = "workspaceOpen"
)
