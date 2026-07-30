package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/permissionrequest"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/pretooluse"
)

// PreToolUse is the PreToolUse hook event.
type PreToolUse = pretooluse.Event

// PreToolOutput is the response for PreToolUse events.
type PreToolOutput = pretooluse.Output

// PreToolResults is the hook-scoped response builder for PreToolUse.
type PreToolResults = pretooluse.Results

// PostToolUse is the PostToolUse hook event.
type PostToolUse = posttooluse.Event

// PostToolOutput is the response for PostToolUse events.
type PostToolOutput = posttooluse.Output

// PostToolResults is the hook-scoped response builder for PostToolUse.
type PostToolResults = posttooluse.Results

// PostToolUseFailure is the PostToolUseFailure hook event.
type PostToolUseFailure = posttoolusefailure.Event

// PostToolFailureOutput is the response for PostToolUseFailure events.
type PostToolFailureOutput = posttoolusefailure.Output

// PostToolFailureResults is the hook-scoped response builder for PostToolUseFailure.
type PostToolFailureResults = posttoolusefailure.Results

// PermissionRequest is the PermissionRequest hook event.
type PermissionRequest = permissionrequest.Event

// PermissionRequestOutput is the response for PermissionRequest events.
type PermissionRequestOutput = permissionrequest.Output

// PermissionRequestResults is the hook-scoped response builder for
// PermissionRequest. Ask is a soft deny (behavior "deny", exit 0), not a user
// confirmation prompt; prefer Deny to block or Noop to fall through.
type PermissionRequestResults = permissionrequest.Results
