package event

import (
	"encoding/json"
	"fmt"
)

// SuccessExit is exit code 0 for successful encode and no-op responses.
const SuccessExit = 0

// HandlerErrorExit is exit code 1 for mux processing failures (read/decode/handler/encode/write errors).
// Claude Code treats exit 1 as a non-blocking error (fail-open) for most events;
// WorktreeCreate is the documented exception where any non-zero exit fails creation.
const HandlerErrorExit = 1

// BlockExit is exit code 2 for Claude command-hook blocking. Claude ignores stdout
// JSON on this exit and feeds stderr to the model; use only for events that cannot
// express the same deny via JSON decision fields (TeammateIdle, TaskCreated,
// TaskCompleted). Prefer exit 0 with decision:"block" elsewhere.
const BlockExit = 2

// MarshalHookOutput builds Claude stdout JSON from top-level and hookSpecificOutput maps.
func MarshalHookOutput(eventName string, fill func(top, hso map[string]any)) ([]byte, int, error) {
	if eventName == "" {
		return nil, SuccessExit, fmt.Errorf("claude: encode: empty event name")
	}
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
