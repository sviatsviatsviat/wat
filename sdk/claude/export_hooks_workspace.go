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

// CwdChangedOutput is the response for CwdChanged events.
type CwdChangedOutput = cwdchanged.Output

// CwdChangedResults is the hook-scoped response builder for CwdChanged.
type CwdChangedResults = cwdchanged.Results

// FileChanged is the FileChanged hook event.
type FileChanged = filechanged.Event

// FileChangedOutput is the response for FileChanged events.
type FileChangedOutput = filechanged.Output

// FileChangedResults is the hook-scoped response builder for FileChanged.
type FileChangedResults = filechanged.Results

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

// ConfigChangeOutput is the response for ConfigChange events.
type ConfigChangeOutput = configchange.Output

// ConfigChangeResults is the hook-scoped response builder for ConfigChange.
type ConfigChangeResults = configchange.Results
