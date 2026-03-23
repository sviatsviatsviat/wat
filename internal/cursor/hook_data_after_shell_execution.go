package cursor

import "github.com/sviatsviatsviat/wat/internal/cursor/core"

// hookDataAfterShellExecution is the Cursor payload shape for afterShellExecution.
type hookDataAfterShellExecution = cursorcore.HookDataWithCommon[hookDataAfterShellExecutionFields]

type hookDataAfterShellExecutionFields struct {
	Command  string `json:"command"`
	Output   string `json:"output"`
	Duration float32 `json:"duration"`
	Sandbox  bool   `json:"sandbox"`
}
