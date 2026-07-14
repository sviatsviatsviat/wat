package model

import (
	"reflect"
	"testing"
)

func TestTypedResults_Result(t *testing.T) {
	out := "rewritten"
	tests := []struct {
		name string
		got  Result
		want Result
	}{
		{
			name: "pre tool deny",
			got:  PreToolDeny("blocked").Result(),
			want: Result{Decision: DecisionDeny, Reason: "blocked"},
		},
		{
			name: "post tool context",
			got:  PostToolContext("note").Result(),
			want: Result{Context: "note"},
		},
		{
			name: "post tool failure context",
			got:  PostToolFailureContext("retry").Result(),
			want: Result{Context: "retry"},
		},
		{
			name: "stop follow up",
			got:  StopFollowUp("keep going").Result(),
			want: Result{FollowUp: "keep going"},
		},
		{
			name: "session start context",
			got:  SessionStartContext("boot").Result(),
			want: Result{Context: "boot"},
		},
		{
			name: "post tool updated output",
			got:  PostToolResult{UpdatedOutput: &out}.Result(),
			want: Result{UpdatedOutput: &out},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("Result() = %+v, want %+v", tt.got, tt.want)
			}
		})
	}
}

func TestTypedResults_IsZero(t *testing.T) {
	if !(PreToolResult{}).IsZero() {
		t.Fatal("zero PreToolResult should be zero")
	}
	if PreToolDeny("x").IsZero() {
		t.Fatal("deny should not be zero")
	}
	if !(PostToolResult{}).IsZero() {
		t.Fatal("zero PostToolResult should be zero")
	}
}
