package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/internal/hookkit"
)

// MergeOutputs combines native Copilot stdout JSON from multiple handlers.
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
		mergeCopilotTop(merged, top)
	}
	if len(merged) == 0 {
		return nil, nil
	}
	return json.Marshal(merged)
}

func mergeCopilotTop(dst, src map[string]any) {
	for k, v := range src {
		switch k {
		case "permissionDecision":
			hookkit.ApplyRankedDecision(dst, src, "permissionDecision", "permissionDecisionReason", v, hookkit.PermissionRank)
		case "permissionDecisionReason":
			hookkit.ApplyOrphanDetail(dst, "permissionDecision", "permissionDecisionReason", v)
		case "additionalContext":
			dst[k] = hookkit.JoinContext(dst[k], v)
		case "decision":
			hookkit.ApplyRankedDecision(dst, src, "decision", "reason", v, hookkit.BlockDecisionRank)
		case "reason":
			hookkit.ApplyOrphanDetail(dst, "decision", "reason", v)
		default:
			if v != nil && v != "" {
				dst[k] = v
			}
		}
	}
}
