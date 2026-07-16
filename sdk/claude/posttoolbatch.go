package claude

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
)

// PostToolBatchCall is one tool invocation entry in a PostToolBatch event.
type PostToolBatchCall struct {
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolResponse is the tool response JSON (string or content-block array).
	ToolResponse json.RawMessage `json:"tool_response"`
}

// PostToolBatch is the PostToolBatch hook event.
type PostToolBatch struct {
	Envelope
	// ToolCalls holds per-call metadata for the resolved parallel batch.
	ToolCalls []PostToolBatchCall `json:"tool_calls"`
}

// EventName returns the hook event name.
func (PostToolBatch) EventName() string { return EventPostToolBatch }

func init() {
	registerDecoder(EventPostToolBatch, func(raw []byte) (Event, error) {
		return decodeAsAndThen(raw, bindPostToolBatchToolInputs)
	})
}

func bindPostToolBatchToolInputs(e *PostToolBatch, raw []byte) {
	var wire struct {
		ToolCalls []struct {
			ToolName  string          `json:"tool_name"`
			ToolInput json.RawMessage `json:"tool_input"`
		} `json:"tool_calls"`
	}
	_ = json.Unmarshal(raw, &wire)
	for i := range e.ToolCalls {
		var input json.RawMessage
		if i < len(wire.ToolCalls) {
			input = cloneToolInputRaw(wire.ToolCalls[i].ToolInput)
		}
		e.ToolCalls[i].ToolInput = tools.NewInput(e.ToolCalls[i].ToolName, input)
	}
}

func cloneToolInputRaw(raw json.RawMessage) json.RawMessage {
	if string(raw) == "null" {
		return nil
	}
	return hookkit.CloneBytes(raw)
}
