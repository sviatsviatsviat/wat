package claude

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is any Claude Code hook response.
type Output = hookkit.Output

type encoder struct{}

func newEncoder() hookkit.Encoder {
	return encoder{}
}

// Encode validates out and renders Claude Code stdout JSON.
func (encoder) Encode(eventName string, out Output) ([]byte, int, error) {
	if eventName == "" {
		return nil, SuccessExit, fmt.Errorf("claude: encode: empty event name")
	}
	if out.IsZero() {
		return nil, SuccessExit, nil
	}
	if err := hookkit.ValidateEncodePair(Dialect, eventName, out, nil); err != nil {
		return nil, SuccessExit, err
	}
	return out.Encode(eventName)
}

// marshalHookOutput builds Claude stdout JSON from top-level and hookSpecificOutput maps.
func marshalHookOutput(eventName string, fill func(top, hso map[string]any)) ([]byte, int, error) {
	top, hso := map[string]any{}, map[string]any{}
	fill(top, hso)
	if len(hso) > 0 {
		hso["hookEventName"] = eventName
		top["hookSpecificOutput"] = hso
	}
	if len(top) == 0 {
		return nil, SuccessExit, nil
	}
	b, err := json.Marshal(top)
	return b, SuccessExit, err
}
