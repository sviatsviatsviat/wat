package hooks

import (
	"context"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// autoModel is Cursor's automatic model selection value for a subagent.
const autoModel = "auto"

// Hooks contains this project's hook registrations.
var Hooks = []run.Hooks{
	agnostic.UseHooks().OnSessionStart(func(ctx context.Context, hook agnostic.SessionStartEvent, r agnostic.SessionStartResults) (agnostic.SessionStartResult, error) {
		return r.Context("wat hooks are active"), nil
	}),
	cursor.UseHooks().SubagentStart(gateSubagentModel),
}

// gateSubagentModel blocks a subagent spawn that pins a model other than Cursor's
// automatic selection, so a human decides whether to allow it. Cursor does not support
// "ask" on subagentStart (it is treated as "deny"), so this denies and explains the next
// step in a user-facing message rather than silently letting the spawn through.
func gateSubagentModel(_ context.Context, hook cursor.SubagentStart, r cursor.SubagentStartResults) (cursor.PermissionOutput, error) {
	model := strings.TrimSpace(hook.SubagentModel)
	if model == "" || strings.EqualFold(model, autoModel) {
		return nil, nil // no opinion: auto, or the host did not report a model
	}
	name := strings.TrimSpace(hook.SubagentType)
	if name == "" {
		name = "subagent"
	}
	return r.Deny("subagent " + name + " requested model " + model + "; only " + autoModel + " is pre-approved").
		WithUserMessage("wat blocked " + name + " from starting on model " + model +
			". Re-run it with the auto model, or allow this model in .wat/hooks.go."), nil
}
