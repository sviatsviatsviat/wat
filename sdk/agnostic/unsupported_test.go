package agnostic

import (
	"slices"
	"testing"
)

func TestUnsupported_capabilityMatrix(t *testing.T) {
	out := "rewritten"

	tests := []struct {
		name    string
		dialect Dialect
		kind    Kind
		result  Result
		want    []string
	}{
		{
			name:    "copilot block prompt unsupported",
			dialect: Copilot,
			kind:    KindUserPrompt,
			result:  Result{BlockPrompt: true},
			want:    []string{"BlockPrompt"},
		},
		{
			name:    "copilot env unsupported",
			dialect: Copilot,
			kind:    KindSessionStart,
			result:  Result{Env: map[string]string{"K": "V"}},
			want:    []string{"Env"},
		},
		{
			name:    "copilot halt session only on permission request",
			dialect: Copilot,
			kind:    KindStop,
			result:  Result{HaltSession: true},
			want:    []string{"HaltSession"},
		},
		{
			name:    "copilot halt session supported on permission request",
			dialect: Copilot,
			kind:    KindPermissionRequest,
			result:  Result{HaltSession: true, Decision: DecisionDeny},
			want:    nil,
		},
		{
			name:    "copilot decision only on pre tool",
			dialect: Copilot,
			kind:    KindStop,
			result:  Deny("nope"),
			want:    []string{"Decision"},
		},
		{
			name:    "copilot context subset",
			dialect: Copilot,
			kind:    KindPreTool,
			result:  Context("extra"),
			want:    []string{"Context"},
		},
		{
			name:    "copilot context supported on post tool",
			dialect: Copilot,
			kind:    KindPostTool,
			result:  Context("extra"),
			want:    nil,
		},
		{
			name:    "copilot context supported on post tool failure",
			dialect: Copilot,
			kind:    KindPostToolFailure,
			result:  Context("extra"),
			want:    nil,
		},
		{
			name:    "cursor updated input only on pre tool",
			dialect: Cursor,
			kind:    KindStop,
			result:  Result{UpdatedInput: map[string]any{"command": "ls"}},
			want:    []string{"UpdatedInput"},
		},
		{
			name:    "cursor halt session unsupported",
			dialect: Cursor,
			kind:    KindStop,
			result:  Result{HaltSession: true},
			want:    []string{"HaltSession"},
		},
		{
			name:    "cursor ask on subagent start treated as deny",
			dialect: Cursor,
			kind:    KindSubagentStart,
			result:  Ask("confirm"),
			want:    []string{"Ask(treated as deny)"},
		},
		{
			name:    "cursor env only on session start",
			dialect: Cursor,
			kind:    KindPreTool,
			result:  Result{Env: map[string]string{"K": "V"}},
			want:    []string{"Env"},
		},
		{
			name:    "cursor env supported on session start",
			dialect: Cursor,
			kind:    KindSessionStart,
			result:  Result{Env: map[string]string{"K": "V"}},
			want:    nil,
		},
		{
			name:    "cursor updated output only on post tool",
			dialect: Cursor,
			kind:    KindPreTool,
			result:  Result{UpdatedOutput: &out},
			want:    []string{"UpdatedOutput"},
		},
		{
			name:    "cursor block prompt only on user prompt",
			dialect: Cursor,
			kind:    KindPreTool,
			result:  Result{BlockPrompt: true},
			want:    []string{"BlockPrompt"},
		},
		{
			name:    "claude pre tool deny supported",
			dialect: Claude,
			kind:    KindPreTool,
			result:  Deny("blocked"),
			want:    nil,
		},
		{
			name:    "claude updated input only on pre tool",
			dialect: Claude,
			kind:    KindStop,
			result:  Result{UpdatedInput: map[string]any{"command": "ls"}},
			want:    []string{"UpdatedInput"},
		},
		{
			name:    "claude updated output only on post tool",
			dialect: Claude,
			kind:    KindStop,
			result:  Result{UpdatedOutput: &out},
			want:    []string{"UpdatedOutput"},
		},
		{
			name:    "claude follow up only on stop",
			dialect: Claude,
			kind:    KindPreTool,
			result:  Result{FollowUp: "keep going"},
			want:    []string{"FollowUp"},
		},
		{
			name:    "claude set title limited surfaces",
			dialect: Claude,
			kind:    KindStop,
			result:  Result{SetTitle: "title"},
			want:    []string{"SetTitle"},
		},
		{
			name:    "claude env only on session start",
			dialect: Claude,
			kind:    KindPreTool,
			result:  Result{Env: map[string]string{"K": "V"}},
			want:    []string{"Env"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unsupported(tt.dialect, tt.kind, tt.result)
			if len(tt.want) == 0 {
				if len(got) != 0 {
					t.Fatalf("Unsupported() = %v, want empty", got)
				}
				return
			}
			for _, w := range tt.want {
				if !slices.Contains(got, w) {
					t.Fatalf("Unsupported() = %v, missing %q", got, w)
				}
			}
		})
	}
}

func TestDecision_String(t *testing.T) {
	tests := []struct {
		d    Decision
		want string
	}{
		{DecisionUnset, ""},
		{DecisionAllow, "allow"},
		{DecisionAsk, "ask"},
		{DecisionDeny, "deny"},
	}
	for _, tt := range tests {
		if got := tt.d.String(); got != tt.want {
			t.Fatalf("%v.String() = %q, want %q", tt.d, got, tt.want)
		}
	}
}
