package claude

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/internal/hookkit"
)

// MergeOutputs combines native Claude stdout JSON from multiple handlers.
// Deny outranks ask outranks allow; additionalContext strings concatenate.
func MergeOutputs(outputs [][]byte) ([]byte, error) {
	if len(outputs) == 0 {
		return nil, nil
	}
	if len(outputs) == 1 {
		return outputs[0], nil
	}

	merged := map[string]any{}
	hso := map[string]any{}
	for _, b := range outputs {
		var top map[string]any
		if err := json.Unmarshal(b, &top); err != nil {
			return nil, err
		}
		mergeClaudeTop(merged, top)
		if raw, ok := top["hookSpecificOutput"].(map[string]any); ok {
			mergeClaudeHSO(hso, raw)
		}
	}
	if len(hso) > 0 {
		if existing, ok := merged["hookSpecificOutput"].(map[string]any); ok {
			mergeClaudeHSO(existing, hso)
			hso = existing
		}
		merged["hookSpecificOutput"] = hso
	}
	if len(merged) == 0 {
		return nil, nil
	}
	return json.Marshal(merged)
}

func mergeClaudeTop(dst, src map[string]any) {
	for k, v := range src {
		if k == "hookSpecificOutput" {
			continue
		}
		switch k {
		case "decision":
			hookkit.ApplyRankedDecision(dst, src, "decision", "reason", v, hookkit.BlockDecisionRank)
		case "reason":
			hookkit.ApplyOrphanDetail(dst, "decision", "reason", v)
		case "continue":
			if v == false {
				dst[k] = v
			}
		default:
			if v != nil && v != "" {
				dst[k] = v
			}
		}
	}
}

func mergeClaudeHSO(dst, src map[string]any) {
	for k, v := range src {
		switch k {
		case "permissionDecision":
			hookkit.ApplyRankedDecision(dst, src, "permissionDecision", "permissionDecisionReason", v, hookkit.PermissionRank)
		case "permissionDecisionReason":
			hookkit.ApplyOrphanDetail(dst, "permissionDecision", "permissionDecisionReason", v)
		case "additionalContext":
			dst[k] = hookkit.JoinContext(dst[k], v)
		default:
			if v != nil && v != "" {
				dst[k] = v
			}
		}
	}
}
