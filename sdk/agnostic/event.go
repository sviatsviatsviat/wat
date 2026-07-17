package agnostic

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// Event is the unified, agent-independent view of a hook invocation.
// Raw always carries the untouched native payload, so nothing is lost by
// normalization; agent-specific handlers can re-decode it with native types.
type Event = model.Event

// ToolCall describes the tool invocation a pre/post tool event refers to.
type ToolCall = model.ToolCall

// ToolResult describes the outcome of a completed or failed tool call.
type ToolResult = model.ToolResult
