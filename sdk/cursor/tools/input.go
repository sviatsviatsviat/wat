package tools

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Input is the tool input payload on a Cursor hook event.
type Input struct {
	hookkit.Input
}

// NewInput returns an Input bound to the native tool name and raw JSON.
// It panics if both tool and raw are empty.
func NewInput(tool string, raw json.RawMessage) Input {
	return Input{Input: hookkit.NewInput(tool, raw)}
}

// NewInputFromPayload returns an Input bound to tool and the JSON object field named field.
// It panics if both tool and the extracted field are empty.
func NewInputFromPayload(tool string, payload []byte, field string) Input {
	return NewInput(tool, hookkit.RawObjectField(payload, field))
}
