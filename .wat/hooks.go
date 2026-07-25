package hooks

import (
	"context"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hooks contains this project's hook registrations.
var Hooks = []run.Hooks{
	agnostic.UseHooks().OnSessionStart(func(ctx context.Context, hook agnostic.SessionStartEvent, r agnostic.SessionStartResults) (agnostic.SessionStartResult, error) {
		return r.Context("wat hooks are active"), nil
	}),
	cursor.UseHooks().SubagentStart(gateSubagentModel),
}

// gateSubagentModel blocks a subagent spawn that pins a concrete model. Cursor reports
// automatic selection on subagent_model as "", "auto", "default", or "inherit".
// Cursor does not support "ask" on subagentStart (it is treated as "deny"), so this
// denies with a user-facing message.
func gateSubagentModel(_ context.Context, hook cursor.SubagentStart, r cursor.SubagentStartResults) (cursor.PermissionOutput, error) {
	sub := strings.TrimSpace(hook.SubagentModel)
	if isUnpinnedSubagentModel(sub) {
		return nil, nil
	}
	name := strings.TrimSpace(hook.SubagentType)
	if name == "" {
		name = "subagent"
	}
	return r.Deny("wat blocked " + name + " from starting on model " + sub +
		". Re-run it with the default/auto model, or allow this model in .wat/hooks.go."), nil
}

// isUnpinnedSubagentModel reports whether subagent_model is unset or a Cursor
// automatic-selection sentinel.
func isUnpinnedSubagentModel(subagentModel string) bool {
	switch strings.ToLower(strings.TrimSpace(subagentModel)) {
	case "", "auto", "default", "inherit":
		return true
	default:
		return false
	}
}
