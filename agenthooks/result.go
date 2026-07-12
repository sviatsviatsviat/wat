package agenthooks

import "maps"

// Decision is the unified gate verdict for pre-events.
type Decision int

const (
	// DecisionUnset means the handler expressed no gate verdict.
	DecisionUnset Decision = iota
	// DecisionAllow permits the gated action to proceed.
	DecisionAllow
	// DecisionAsk escalates the decision to the user.
	DecisionAsk
	// DecisionDeny blocks the gated action.
	DecisionDeny
)

// String returns the agent-facing decision label ("allow", "deny", "ask")
// or an empty string for DecisionUnset.
func (d Decision) String() string {
	switch d {
	case DecisionAllow:
		return "allow"
	case DecisionDeny:
		return "deny"
	case DecisionAsk:
		return "ask"
	default:
		return ""
	}
}

// Result is the unified hook response. The zero value means "no decision":
// no stdout is produced and the agent's default flow applies. Each dialect
// codec encodes only the fields the target agent supports for the given
// event kind; Unsupported reports the remainder.
//
// Result is distinct from Event.Result, which carries incoming post-tool
// payload data on decoded events.
type Result struct {
	// Decision is the gate verdict for PreTool and PermissionRequest events.
	Decision Decision
	// Reason is an agent- or model-facing explanation.
	Reason string
	// UserMessage is a user-facing message.
	UserMessage string
	// UpdatedInput replaces tool arguments on PreTool events.
	UpdatedInput map[string]any
	// UpdatedOutput replaces tool result text on PostTool events.
	UpdatedOutput *string
	// Context is additional context injected for the model.
	Context string
	// FollowUp instructs the agent to keep working on Stop and SubagentStop events.
	FollowUp string
	// BlockPrompt rejects the submitted prompt on UserPrompt events.
	BlockPrompt bool
	// HaltSession stops the whole session where the target agent supports it.
	HaltSession bool
	// Env carries session environment variables, primarily on SessionStart.
	Env map[string]string
	// SetTitle sets the session title where the target agent supports it.
	SetTitle string
}

// Allow returns an allow verdict.
func Allow() Result { return Result{Decision: DecisionAllow} }

// Deny returns a deny verdict with an agent-facing reason.
func Deny(reason string) Result { return Result{Decision: DecisionDeny, Reason: reason} }

// Ask returns an escalate-to-user verdict with an agent-facing reason.
func Ask(reason string) Result { return Result{Decision: DecisionAsk, Reason: reason} }

// Context returns a context-injection-only result.
func Context(text string) Result { return Result{Context: text} }

// IsZero reports whether the result carries no instruction at all.
func (r Result) IsZero() bool {
	return r.Decision == DecisionUnset && r.Reason == "" && r.UserMessage == "" &&
		r.UpdatedInput == nil && r.UpdatedOutput == nil && r.Context == "" &&
		r.FollowUp == "" && !r.BlockPrompt && !r.HaltSession && len(r.Env) == 0 &&
		r.SetTitle == ""
}

// Merge combines b into a. Deny outranks Ask outranks Allow; text fields are
// overridden by non-empty values from b, except Context which accumulates.
func Merge(a, b Result) Result {
	if b.Decision > a.Decision {
		a.Decision = b.Decision
	}
	if b.Reason != "" {
		a.Reason = b.Reason
	}
	if b.UserMessage != "" {
		a.UserMessage = b.UserMessage
	}
	if b.UpdatedInput != nil {
		a.UpdatedInput = b.UpdatedInput
	}
	if b.UpdatedOutput != nil {
		a.UpdatedOutput = b.UpdatedOutput
	}
	if b.Context != "" {
		if a.Context != "" {
			a.Context += "\n\n" + b.Context
		} else {
			a.Context = b.Context
		}
	}
	if b.FollowUp != "" {
		a.FollowUp = b.FollowUp
	}
	a.BlockPrompt = a.BlockPrompt || b.BlockPrompt
	a.HaltSession = a.HaltSession || b.HaltSession
	if a.Env != nil {
		a.Env = maps.Clone(a.Env)
	}
	if len(b.Env) > 0 {
		if a.Env == nil {
			a.Env = make(map[string]string, len(b.Env))
		}
		for k, v := range b.Env {
			a.Env[k] = v
		}
	}
	if b.SetTitle != "" {
		a.SetTitle = b.SetTitle
	}
	return a
}

// Unsupported lists Result capabilities that dialect d cannot express for
// event kind k. Codecs drop these fields silently; callers may warn.
func Unsupported(d Dialect, k Kind, r Result) []string {
	var out []string
	switch d {
	case Claude:
		unsupportedClaude(k, r, &out)
	case Copilot:
		unsupportedCopilot(k, r, &out)
	case Cursor:
		unsupportedCursor(k, r, &out)
	}
	return out
}

func unsupportedClaude(k Kind, r Result, out *[]string) {
	if r.SetTitle != "" && k != KindSessionStart && k != KindUserPrompt {
		*out = append(*out, "SetTitle")
	}
	if len(r.Env) > 0 && k != KindSessionStart {
		*out = append(*out, "Env")
	}
	if r.Decision != DecisionUnset && k != KindPreTool && k != KindPermissionRequest {
		*out = append(*out, "Decision")
	}
	if r.UpdatedInput != nil && k != KindPreTool && k != KindPermissionRequest {
		*out = append(*out, "UpdatedInput")
	}
	if r.UpdatedOutput != nil && k != KindPostTool && k != KindPostToolFailure {
		*out = append(*out, "UpdatedOutput")
	}
	if r.FollowUp != "" && k != KindStop && k != KindSubagentStop {
		*out = append(*out, "FollowUp")
	}
	if r.BlockPrompt && k != KindUserPrompt {
		*out = append(*out, "BlockPrompt")
	}
}

func unsupportedCopilot(k Kind, r Result, out *[]string) {
	if r.BlockPrompt {
		*out = append(*out, "BlockPrompt")
	}
	if len(r.Env) > 0 {
		*out = append(*out, "Env")
	}
	if r.SetTitle != "" {
		*out = append(*out, "SetTitle")
	}
	if r.HaltSession && k != KindPermissionRequest {
		*out = append(*out, "HaltSession")
	}
	if r.Decision != DecisionUnset && k != KindPreTool && k != KindPermissionRequest {
		*out = append(*out, "Decision")
	}
	if r.UpdatedInput != nil && k != KindPreTool {
		*out = append(*out, "UpdatedInput")
	}
	if r.UpdatedOutput != nil && k != KindPostTool {
		*out = append(*out, "UpdatedOutput")
	}
	if r.FollowUp != "" && k != KindStop && k != KindSubagentStop {
		*out = append(*out, "FollowUp")
	}
	if r.Context != "" && !copilotContextKinds[k] {
		*out = append(*out, "Context")
	}
}

func unsupportedCursor(k Kind, r Result, out *[]string) {
	if r.HaltSession {
		*out = append(*out, "HaltSession")
	}
	if r.SetTitle != "" {
		*out = append(*out, "SetTitle")
	}
	if r.Decision == DecisionAsk && k == KindSubagentStart {
		*out = append(*out, "Ask(treated as deny)")
	}
	if len(r.Env) > 0 && k != KindSessionStart {
		*out = append(*out, "Env")
	}
	if r.UpdatedOutput != nil && k != KindPostTool {
		*out = append(*out, "UpdatedOutput")
	}
	if r.UpdatedInput != nil && k != KindPreTool {
		*out = append(*out, "UpdatedInput")
	}
	if r.Decision != DecisionUnset && k != KindPreTool && k != KindSubagentStart {
		*out = append(*out, "Decision")
	}
	if r.BlockPrompt && k != KindUserPrompt {
		*out = append(*out, "BlockPrompt")
	}
}

var copilotContextKinds = map[Kind]bool{
	KindSessionStart:    true,
	KindSubagentStart:   true,
	KindNotification:    true,
	KindPostTool:        true,
	KindPostToolFailure: true,
}
