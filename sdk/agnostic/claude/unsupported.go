package claude

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// Unsupported lists Result capabilities Claude cannot express for event kind k.
func Unsupported(k model.Kind, r model.Result) []string {
	var out []string
	if r.SetTitle != "" && k != model.KindSessionStart && k != model.KindUserPrompt {
		out = append(out, "SetTitle")
	}
	if len(r.Env) > 0 && k != model.KindSessionStart {
		out = append(out, "Env")
	}
	if r.Decision != model.DecisionUnset && k != model.KindPreTool && k != model.KindPermissionRequest {
		out = append(out, "Decision")
	}
	if r.UpdatedInput != nil && k != model.KindPreTool && k != model.KindPermissionRequest {
		out = append(out, "UpdatedInput")
	}
	if r.UpdatedOutput != nil && k != model.KindPostTool && k != model.KindPostToolFailure {
		out = append(out, "UpdatedOutput")
	}
	if r.FollowUp != "" && k != model.KindStop && k != model.KindSubagentStop {
		out = append(out, "FollowUp")
	}
	if r.BlockPrompt && k != model.KindUserPrompt {
		out = append(out, "BlockPrompt")
	}
	return out
}
