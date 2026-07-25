package installcfg

import (
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// Install event names per agent. These are the native hook_event_name values
// wat install writes and doctor validates.
var (
	claudeEvents = []string{
		sdkclaude.EventSessionStart,
		sdkclaude.EventSessionEnd,
		sdkclaude.EventUserPromptSubmit,
		sdkclaude.EventPreToolUse,
		sdkclaude.EventPostToolUse,
		sdkclaude.EventPostToolUseFailure,
		sdkclaude.EventPermissionRequest,
		sdkclaude.EventSubagentStart,
		sdkclaude.EventSubagentStop,
		sdkclaude.EventStop,
		sdkclaude.EventPreCompact,
		sdkclaude.EventNotification,
		sdkclaude.EventStopFailure,
	}

	copilotEvents = []string{
		sdkcopilot.EventSessionStart,
		sdkcopilot.EventSessionEnd,
		sdkcopilot.EventUserPromptSubmitted,
		sdkcopilot.EventPreToolUse,
		sdkcopilot.EventPostToolUse,
		sdkcopilot.EventPostToolUseFailure,
		sdkcopilot.EventPermissionRequest,
		sdkcopilot.EventSubagentStart,
		sdkcopilot.EventSubagentStop,
		sdkcopilot.EventAgentStop,
		sdkcopilot.EventPreCompact,
		sdkcopilot.EventNotification,
		sdkcopilot.EventErrorOccurred,
	}

	// cursorEvents includes portable surfaces plus Cursor-only dedicated events
	// (shell/MCP/file read/edit) that install also wires to wat run.
	cursorEvents = []string{
		sdkcursor.EventSessionStart,
		sdkcursor.EventSessionEnd,
		sdkcursor.EventBeforeSubmitPrompt,
		sdkcursor.EventPreToolUse,
		sdkcursor.EventPostToolUse,
		sdkcursor.EventPostToolUseFailure,
		sdkcursor.EventSubagentStart,
		sdkcursor.EventSubagentStop,
		sdkcursor.EventStop,
		sdkcursor.EventPreCompact,
		sdkcursor.EventBeforeShellExecution,
		sdkcursor.EventAfterShellExecution,
		sdkcursor.EventBeforeMCPExecution,
		sdkcursor.EventAfterMCPExecution,
		sdkcursor.EventBeforeReadFile,
		sdkcursor.EventAfterFileEdit,
	}
)
