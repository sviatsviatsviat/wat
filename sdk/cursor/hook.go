package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hook is the handler context for a typed Cursor hook event.
type Hook[E Event] struct {
	// Event is the decoded native hook payload.
	Event E
	inv   run.Invocation
}

// NewHook wraps ev with serve-time invocation settings.
func NewHook[E Event](inv run.Invocation, ev E) Hook[E] {
	return Hook[E]{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h Hook[E]) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload when available.
func (h Hook[E]) Raw() json.RawMessage { return RawBytes(h.Event) }

type (
	// PreToolUseHook is the handler context for preToolUse events.
	PreToolUseHook = Hook[PreToolUse]
	// PostToolUseHook is the handler context for postToolUse events.
	PostToolUseHook = Hook[PostToolUse]
	// PostToolUseFailureHook is the handler context for postToolUseFailure events.
	PostToolUseFailureHook = Hook[PostToolUseFailure]
	// BeforeShellExecutionHook is the handler context for beforeShellExecution events.
	BeforeShellExecutionHook = Hook[BeforeShellExecution]
	// BeforeMCPExecutionHook is the handler context for beforeMCPExecution events.
	BeforeMCPExecutionHook = Hook[BeforeMCPExecution]
	// BeforeReadFileHook is the handler context for beforeReadFile events.
	BeforeReadFileHook = Hook[BeforeReadFile]
	// BeforeTabFileReadHook is the handler context for beforeTabFileRead events.
	BeforeTabFileReadHook = Hook[BeforeTabFileRead]
	// AfterShellExecutionHook is the handler context for afterShellExecution events.
	AfterShellExecutionHook = Hook[AfterShellExecution]
	// AfterMCPExecutionHook is the handler context for afterMCPExecution events.
	AfterMCPExecutionHook = Hook[AfterMCPExecution]
	// AfterFileEditHook is the handler context for afterFileEdit events.
	AfterFileEditHook = Hook[AfterFileEdit]
	// BeforeSubmitPromptHook is the handler context for beforeSubmitPrompt events.
	BeforeSubmitPromptHook = Hook[BeforeSubmitPrompt]
	// StopHook is the handler context for stop events.
	StopHook = Hook[Stop]
	// SubagentStopHook is the handler context for subagentStop events.
	SubagentStopHook = Hook[SubagentStop]
	// SubagentStartHook is the handler context for subagentStart events.
	SubagentStartHook = Hook[SubagentStart]
	// SessionStartHook is the handler context for sessionStart events.
	SessionStartHook = Hook[SessionStart]
	// SessionEndHook is the handler context for sessionEnd events.
	SessionEndHook = Hook[SessionEnd]
	// PreCompactHook is the handler context for preCompact events.
	PreCompactHook = Hook[PreCompact]
	// AfterAgentResponseHook is the handler context for afterAgentResponse events.
	AfterAgentResponseHook = Hook[AfterAgentResponse]
	// AfterAgentThoughtHook is the handler context for afterAgentThought events.
	AfterAgentThoughtHook = Hook[AfterAgentThought]
	// AfterTabFileEditHook is the handler context for afterTabFileEdit events.
	AfterTabFileEditHook = Hook[AfterTabFileEdit]
	// WorkspaceOpenHook is the handler context for workspaceOpen events.
	WorkspaceOpenHook = Hook[WorkspaceOpen]
)

// AnyHook is the handler context for catch-all OnAny handlers.
type AnyHook struct {
	Event
	inv run.Invocation
}

// NewAnyHook wraps ev with serve-time invocation settings.
func NewAnyHook(inv run.Invocation, ev Event) AnyHook {
	return AnyHook{Event: ev, inv: inv}
}

// Invocation returns serve-time settings for this hook invocation.
func (h AnyHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload when available.
func (h AnyHook) Raw() json.RawMessage { return RawBytes(h.Event) }
