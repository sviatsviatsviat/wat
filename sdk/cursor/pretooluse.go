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
	codec.Register(EventPreToolUse, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PreToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PreToolUse registers a PreToolUse handler on the chain.
func (c *chain) PreToolUse(fn func(context.Context, run.Hook[PreToolUse], PermissionResults) (PermissionOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[PreToolUse]) (PermissionOutput, error) {
		return fn(ctx, hook, permissionResults{})
	}))
	return c
}
