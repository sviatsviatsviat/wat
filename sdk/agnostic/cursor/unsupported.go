package cursor

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// Unsupported lists Result capabilities Cursor cannot express for event kind k.
func Unsupported(k model.Kind, r model.Result) []string {
	var out []string
	if r.HaltSession {
		out = append(out, "HaltSession")
	}
	if r.SetTitle != "" {
		out = append(out, "SetTitle")
	}
	if r.Decision == model.DecisionAsk && k == model.KindSubagentStart {
		out = append(out, "Ask(treated as deny)")
	}
	if len(r.Env) > 0 && k != model.KindSessionStart {
		out = append(out, "Env")
	}
	if r.UpdatedOutput != nil && k != model.KindPostTool {
		out = append(out, "UpdatedOutput")
	}
	if r.UpdatedInput != nil && k != model.KindPreTool {
		out = append(out, "UpdatedInput")
	}
	if r.Decision != model.DecisionUnset && k != model.KindPreTool && k != model.KindSubagentStart {
		out = append(out, "Decision")
	}
	if r.BlockPrompt && k != model.KindUserPrompt {
		out = append(out, "BlockPrompt")
	}
	return out
}
