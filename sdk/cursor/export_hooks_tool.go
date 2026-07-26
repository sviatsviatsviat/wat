package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/afterfileedit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftermcpexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftershellexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftertabfileedit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforemcpexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforereadfile"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforeshellexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforetabfileread"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/pretooluse"
)

// PreToolUse is the preToolUse hook event.
type PreToolUse = pretooluse.Event

// PostToolUse is the postToolUse hook event.
type PostToolUse = posttooluse.Event

// PostToolUseFailure is the postToolUseFailure hook event.
type PostToolUseFailure = posttoolusefailure.Event

// BeforeShellExecution is the beforeShellExecution hook event.
type BeforeShellExecution = beforeshellexecution.Event

// AfterShellExecution is the afterShellExecution hook event.
type AfterShellExecution = aftershellexecution.Event

// BeforeMCPExecution is the beforeMCPExecution hook event.
type BeforeMCPExecution = beforemcpexecution.Event

// AfterMCPExecution is the afterMCPExecution hook event.
type AfterMCPExecution = aftermcpexecution.Event

// BeforeReadFile is the beforeReadFile hook event.
type BeforeReadFile = beforereadfile.Event

// BeforeReadFileResults is the hook-scoped response builder for BeforeReadFile.
type BeforeReadFileResults = beforereadfile.Results

// AfterFileEdit is the afterFileEdit hook event.
type AfterFileEdit = afterfileedit.Event

// BeforeTabFileRead is the beforeTabFileRead hook event (Tab completions only).
type BeforeTabFileRead = beforetabfileread.Event

// BeforeTabFileReadResults is the hook-scoped response builder for BeforeTabFileRead.
// Cursor accepts permission allow|deny only; there is no ask or message field.
type BeforeTabFileReadResults = beforetabfileread.Results

// AfterTabFileEdit is the afterTabFileEdit hook event.
type AfterTabFileEdit = aftertabfileedit.Event
