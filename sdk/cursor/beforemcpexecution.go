package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// BeforeMCPExecution is the beforeMCPExecution hook event.
type BeforeMCPExecution struct {
	Envelope
	// ToolName is the MCP tool name (MCP: prefix).
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// URL is the remote MCP server URL when present on the wire.
	URL string `json:"url"`
	// Command is the stdio MCP server command when present on the wire.
	Command string `json:"command"`
}

// EventName returns the canonical hook event name.
func (BeforeMCPExecution) EventName() string { return EventBeforeMCPExecution }

func init() {
	registerDecoder(EventBeforeMCPExecution, decodeAs[BeforeMCPExecution])
}

// BeforeMCPExecution registers a beforeMCPExecution handler.
func (c *Chain) BeforeMCPExecution(fn func(context.Context, BeforeMCPExecutionHook, PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeMCPExecution) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return &Chain{}
}
