package agnostic

import (
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
