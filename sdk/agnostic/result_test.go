package agnostic

import (
	"maps"
	"slices"
	"testing"
)

func TestMerge_decisionPriority(t *testing.T) {
	tests := []struct {
		name string
		a, b Result
		want Decision
	}{
		{name: "deny outranks allow", a: Allow(), b: Deny("x"), want: DecisionDeny},
		{name: "deny sticks once set", a: Deny("x"), b: Allow(), want: DecisionDeny},
		{name: "deny outranks ask from allow", a: Deny("x"), b: Ask("y"), want: DecisionDeny},
		{name: "deny outranks ask from ask", a: Ask("y"), b: Deny("x"), want: DecisionDeny},
		{name: "ask outranks allow", a: Allow(), b: Ask("y"), want: DecisionAsk},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.a, tt.b)
			if got.Decision != tt.want {
				t.Fatalf("Merge decision = %v, want %v", got.Decision, tt.want)
			}
		})
	}
}

func TestMerge_contextAccumulation(t *testing.T) {
	got := Merge(Context("a"), Context("b"))
	if got.Context != "a\n\nb" {
		t.Fatalf("context = %q, want %q", got.Context, "a\n\nb")
	}
}

func TestMerge_otherFields(t *testing.T) {
	t.Run("reason override", func(t *testing.T) {
		got := Merge(Deny("first"), Deny("second"))
		if got.Reason != "second" {
			t.Fatalf("reason = %q, want second", got.Reason)
		}
	})

	t.Run("env merge", func(t *testing.T) {
		a := Result{Env: map[string]string{"A": "1", "B": "2"}}
		b := Result{Env: map[string]string{"B": "override", "C": "3"}}
		origA := maps.Clone(a.Env)
		got := Merge(a, b)
		if got.Env["A"] != "1" || got.Env["B"] != "override" || got.Env["C"] != "3" {
			t.Fatalf("env = %v", got.Env)
		}
		if a.Env["B"] != origA["B"] {
			t.Fatalf("caller env mutated: %v, want %v", a.Env, origA)
		}
	})

	t.Run("env clone when b empty", func(t *testing.T) {
		a := Result{Env: map[string]string{"A": "1"}}
		origA := maps.Clone(a.Env)
		got := Merge(a, Result{})
		got.Env["A"] = "mutated"
		if !maps.Equal(a.Env, origA) {
			t.Fatalf("caller env mutated: %v, want %v", a.Env, origA)
		}
		if got.Env["A"] != "mutated" {
			t.Fatalf("returned env = %v", got.Env)
		}
	})

	t.Run("bool OR", func(t *testing.T) {
		got := Merge(Result{BlockPrompt: false}, Result{HaltSession: true})
		if got.BlockPrompt || !got.HaltSession {
			t.Fatalf("got BlockPrompt=%v HaltSession=%v", got.BlockPrompt, got.HaltSession)
		}
		got = Merge(got, Result{BlockPrompt: true})
		if !got.BlockPrompt || !got.HaltSession {
			t.Fatalf("got BlockPrompt=%v HaltSession=%v", got.BlockPrompt, got.HaltSession)
		}
	})

	t.Run("updated output replacement", func(t *testing.T) {
		first := "first"
		second := "second"
		got := Merge(Result{UpdatedOutput: &first}, Result{UpdatedOutput: &second})
		if got.UpdatedOutput == nil || *got.UpdatedOutput != "second" {
			t.Fatalf("updated output = %v", got.UpdatedOutput)
		}
	})
}

func TestResult_IsZero(t *testing.T) {
	tests := []struct {
		name string
		r    Result
		want bool
	}{
		{name: "zero value", r: Result{}, want: true},
		{name: "allow", r: Allow(), want: false},
		{name: "deny", r: Deny("r"), want: false},
		{name: "context", r: Context("c"), want: false},
		{name: "block prompt", r: Result{BlockPrompt: true}, want: false},
		{name: "env", r: Result{Env: map[string]string{"K": "V"}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsZero(); got != tt.want {
				t.Fatalf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstructors(t *testing.T) {
	if Allow().Decision != DecisionAllow {
		t.Fatal("Allow() should set DecisionAllow")
	}
	if got := Deny("blocked"); got.Decision != DecisionDeny || got.Reason != "blocked" {
		t.Fatalf("Deny() = %+v", got)
	}
	if got := Ask("confirm"); got.Decision != DecisionAsk || got.Reason != "confirm" {
		t.Fatalf("Ask() = %+v", got)
	}
	if got := Context("note"); got.Context != "note" || got.Decision != DecisionUnset {
		t.Fatalf("Context() = %+v", got)
	}
}

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
