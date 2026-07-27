package posttoolbatch

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/tools"
)

// Call is one tool invocation entry in a PostToolBatch event.
type Call struct {
	event.ToolFields
	// ToolResponse is the tool response JSON (string or content-block array).
	ToolResponse json.RawMessage `json:"tool_response"`
}

// Event is the PostToolBatch hook event.
type Event struct {
	event.Envelope
	// ToolCalls holds per-call metadata for the resolved parallel batch.
	ToolCalls []Call `json:"tool_calls"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PostToolBatch }

func bindToolInputs(e *Event, raw []byte) {
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
			input = hookkit.NullToNil(wire.ToolCalls[i].ToolInput)
		}
		e.ToolCalls[i].ToolInput = tools.NewInput(e.ToolCalls[i].ToolName, input)
	}
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolBatch, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, bindToolInputs)
	})
}
