package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolUse is the preToolUse hook event.
type PreToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the canonical hook event name.
func (PreToolUse) EventName() string { return EventPreToolUse }

// ShellCommand extracts the shell command when the tool is Shell.
func (e PreToolUse) ShellCommand() string {
	if e.ToolName != "Shell" && e.ToolName != "shell" {
		return ""
	}
	return hookkit.ExtractShellCommand(e.ToolInput.Raw())
}

func init() {
	registerDecoder(EventPreToolUse, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *PreToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PreToolUse registers a preToolUse handler.
func (c *Chain) PreToolUse(fn func(context.Context, Hook[PreToolUse], PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return c
}
