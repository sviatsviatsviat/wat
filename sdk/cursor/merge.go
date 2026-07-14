package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// MergeOutputs combines native Cursor stdout JSON from multiple handlers.
func MergeOutputs(outputs [][]byte) ([]byte, error) {
	if len(outputs) == 0 {
		return nil, nil
	}
	if len(outputs) == 1 {
		return outputs[0], nil
	}

	merged := map[string]any{}
	for _, b := range outputs {
		var top map[string]any
		if err := json.Unmarshal(b, &top); err != nil {
			return nil, err
		}
		mergeCursorTop(merged, top)
	}
	if len(merged) == 0 {
		return nil, nil
	}
	return json.Marshal(merged)
}

func mergeCursorTop(dst, src map[string]any) {
	for k, v := range src {
		switch k {
		case "permission":
			hookkit.ApplyRankedDecision(dst, src, "permission", "agent_message", v, hookkit.PermissionRank)
		case "agent_message":
			hookkit.ApplyOrphanDetail(dst, "permission", "agent_message", v)
		case "additional_context":
			dst[k] = hookkit.JoinContext(dst[k], v)
		case "followup_message":
			if v != "" {
				dst[k] = v
			}
		default:
			if v != nil && v != "" {
				dst[k] = v
			}
		}
	}
}
