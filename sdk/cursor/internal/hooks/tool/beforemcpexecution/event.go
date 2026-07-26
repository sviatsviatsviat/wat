package beforemcpexecution

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/tools"
)

// Event is the beforeMCPExecution hook event. Cursor calls it before an MCP
// tool runs so authors can allow, deny, or ask. The wire payload includes
// tool_name and tool_input plus either url (remote MCP server) or command
// (stdio MCP server).
//
// Cursor's default hook failure policy is fail-open. For security-critical
// MCP gates, set failClosed: true on the hooks.json handler so crash, timeout,
// or invalid JSON blocks the tool instead of allowing it. This event is
// deferred / not available for cloud agents per Cursor Hooks docs.
type Event struct {
	event.Envelope
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
func (Event) EventName() string { return event.BeforeMCPExecution }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.BeforeMCPExecution, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}
