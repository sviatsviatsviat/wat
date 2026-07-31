package event

import (
	"reflect"
	"testing"
)

func TestClonePermissionUpdates(t *testing.T) {
	if ClonePermissionUpdates(nil) != nil {
		t.Fatal("nil should stay nil")
	}
	src := []PermissionUpdate{{
		Type:        PermissionUpdateAddRules,
		Behavior:    DecisionAllow,
		Destination: PermissionDestinationSession,
		Rules:       []PermissionRule{{ToolName: "Bash", RuleContent: "ls"}},
		Directories: []string{"/a"},
	}}
	got := ClonePermissionUpdates(src)
	if !reflect.DeepEqual(got, src) {
		t.Fatalf("got %+v want %+v", got, src)
	}
	got[0].Rules[0].ToolName = "Write"
	got[0].Directories[0] = "/b"
	if src[0].Rules[0].ToolName != "Bash" || src[0].Directories[0] != "/a" {
		t.Fatal("clone must not alias nested slices")
	}
}
