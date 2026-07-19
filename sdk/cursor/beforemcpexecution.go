package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// BeforeMCPExecution is the beforeMCPExecution hook event.
type BeforeMCPExecution struct {
	Envelope
	// ToolName is the native tool name (typically MCP:<tool>).
	ToolName string `json:"tool_name"`
	// ToolInput is the tool arguments from tool_input, bound to ToolName after decode.
	ToolInput tools.Input `json:"-"`
	// URL is the remote MCP server URL when present on the wire.
	URL string `json:"url"`
	// Command is the stdio MCP server command when present on the wire.
	Command string `json:"command"`
}

// EventName returns the canonical hook event name.
func (BeforeMCPExecution) EventName() string { return EventBeforeMCPExecution }

func init() {
	codec.Register(EventBeforeMCPExecution, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *BeforeMCPExecution, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// BeforeMCPExecution registers a BeforeMCPExecution handler on the chain.
func (c *chain) BeforeMCPExecution(fn func(context.Context, run.Hook[BeforeMCPExecution], PermissionResults) (PermissionOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(c.reg, func(ctx context.Context, ev BeforeMCPExecution) (PermissionOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return c
}
