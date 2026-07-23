package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/configchange"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/cwdchanged"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/filechanged"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/instructionsloaded"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreecreate"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreeremove"
)

// CwdChanged is the CwdChanged hook event.
type CwdChanged = cwdchanged.Event

// FileChanged is the FileChanged hook event.
type FileChanged = filechanged.Event

// WorktreeCreate is the WorktreeCreate hook event.
type WorktreeCreate = worktreecreate.Event

// WorktreeCreateOutput is the response for WorktreeCreate events.
type WorktreeCreateOutput = worktreecreate.Output

// WorktreeCreateResults is the hook-scoped response builder for WorktreeCreate.
type WorktreeCreateResults = worktreecreate.Results

// WorktreeRemove is the WorktreeRemove hook event.
type WorktreeRemove = worktreeremove.Event

// InstructionsLoaded is the InstructionsLoaded hook event.
type InstructionsLoaded = instructionsloaded.Event

// ConfigChange is the ConfigChange hook event.
type ConfigChange = configchange.Event
