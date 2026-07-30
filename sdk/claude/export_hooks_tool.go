package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/permissiondenied"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/permissionrequest"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttoolbatch"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/pretooluse"
)

// PreToolUse is the PreToolUse hook event.
type PreToolUse = pretooluse.Event

// PreToolUseOutput is the response for PreToolUse events.
type PreToolUseOutput = pretooluse.Output

// PreToolUseResults is the hook-scoped response builder for PreToolUse.
type PreToolUseResults = pretooluse.Results

// PostToolUse is the PostToolUse hook event.
type PostToolUse = posttooluse.Event

// PostToolUseOutput is the response for PostToolUse and PostToolUseFailure events.
type PostToolUseOutput = posttooluse.Output

// PostToolUseResults is the hook-scoped response builder for PostToolUse.
type PostToolUseResults = posttooluse.Results

// PostToolUseFailure is the PostToolUseFailure hook event.
type PostToolUseFailure = posttoolusefailure.Event

// PostToolUseFailureResults is the hook-scoped response builder for PostToolUseFailure.
type PostToolUseFailureResults = posttoolusefailure.Results

// PostToolBatch is the PostToolBatch hook event.
type PostToolBatch = posttoolbatch.Event

// PostToolBatchCall is one tool invocation entry in a PostToolBatch event.
type PostToolBatchCall = posttoolbatch.Call

// PostToolBatchResults is the hook-scoped response builder for PostToolBatch.
type PostToolBatchResults = posttoolbatch.Results

// PostToolBatchOutput is the response for PostToolBatch events.
type PostToolBatchOutput = DecisionOutput

// PermissionRequest is the PermissionRequest hook event.
type PermissionRequest = permissionrequest.Event

// PermissionRequestOutput is the response for PermissionRequest events.
type PermissionRequestOutput = permissionrequest.Output

// PermissionRequestResults is the hook-scoped response builder for PermissionRequest.
type PermissionRequestResults = permissionrequest.Results

// PermissionDenied is the PermissionDenied hook event.
type PermissionDenied = permissiondenied.Event

// PermissionDeniedOutput is the response for PermissionDenied events.
type PermissionDeniedOutput = permissiondenied.Output

// PermissionDeniedResults is the hook-scoped response builder for PermissionDenied.
type PermissionDeniedResults = permissiondenied.Results
