package copilot

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

var contextKinds = map[model.Kind]bool{
	model.KindSessionStart:    true,
	model.KindSubagentStart:   true,
	model.KindNotification:    true,
	model.KindPostTool:        true,
	model.KindPostToolFailure: true,
}

// Unsupported lists Result capabilities Copilot cannot express for event kind k.
func Unsupported(k model.Kind, r model.Result) []string {
	var out []string
	if r.BlockPrompt {
		out = append(out, "BlockPrompt")
	}
	if len(r.Env) > 0 {
		out = append(out, "Env")
	}
	if r.SetTitle != "" {
		out = append(out, "SetTitle")
	}
	if r.HaltSession && k != model.KindPermissionRequest {
		out = append(out, "HaltSession")
	}
	if r.Decision != model.DecisionUnset && k != model.KindPreTool && k != model.KindPermissionRequest {
		out = append(out, "Decision")
	}
	if r.UpdatedInput != nil && k != model.KindPreTool {
		out = append(out, "UpdatedInput")
	}
	if r.UpdatedOutput != nil && k != model.KindPostTool {
		out = append(out, "UpdatedOutput")
	}
	if r.FollowUp != "" && k != model.KindStop && k != model.KindSubagentStop {
		out = append(out, "FollowUp")
	}
	if r.Context != "" && !contextKinds[k] {
		out = append(out, "Context")
	}
	return out
}
