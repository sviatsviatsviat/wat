package event

import (
	"maps"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

func TestEncode_ZeroOutput(t *testing.T) {
	out := permissionResults{}.Noop()
	if !out.IsZero() {
		t.Fatal("noop should be zero")
	}
}

func TestEncode_PermissionUpdatedInput(t *testing.T) {
	out, code, err := permissionResults{}.Allow().WithUpdatedInput(map[string]any{"command": "ls"}).Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"updated_input"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestMerge_Permission_denyBeatsAllowAndStops(t *testing.T) {
	a := permissionResults{}.Allow().WithUserMessage("ok")
	b := permissionResults{}.Deny("blocked")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(permissionOutput)
	if out.decision != DecisionDeny || out.agentMessage != "blocked" {
		t.Fatalf("got %#v", out)
	}
	if !merged.Stop() {
		t.Fatal("deny should stop")
	}
}

func TestMerge_Permission_updatedInputOverwriteWarns(t *testing.T) {
	first := map[string]any{"command": "a"}
	second := map[string]any{"command": "b"}
	orig := maps.Clone(first)
	a := permissionResults{}.Allow().WithUpdatedInput(first)
	b := permissionResults{}.Allow().WithUpdatedInput(second)
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] != "updatedInput: overwritten by later handler" {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(permissionOutput)
	if out.updatedInput["command"] != "b" {
		t.Fatalf("updatedInput = %v", out.updatedInput)
	}
	if !maps.Equal(first, orig) {
		t.Fatalf("caller map mutated")
	}
}

func TestMerge_Permission_askDoesNotStop(t *testing.T) {
	out := permissionResults{}.Ask("confirm")
	if out.(permissionOutput).Stop() {
		t.Fatal("ask should not stop")
	}
}

func TestEncode_PermissionOnly_stripsMessages(t *testing.T) {
	out, code, err := GateResults{}.PermissionOnlyDeny().
		WithUserMessage("user").
		WithAgentMessage("agent").
		Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if got != `{"permission":"deny"}` {
		t.Fatalf("got %s, want permission-only deny", got)
	}
}

func TestEncode_DenyUserMessage_stripsChainedAgentFields(t *testing.T) {
	out, code, err := GateResults{}.DenyUserMessage("blocked").
		WithAgentMessage("agent").
		WithUpdatedInput(map[string]any{"path": "secret"}).
		Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if !strings.Contains(got, `"permission":"deny"`) {
		t.Fatalf("missing permission deny: %s", got)
	}
	if !strings.Contains(got, `"user_message":"blocked"`) {
		t.Fatalf("missing user_message: %s", got)
	}
	if strings.Contains(got, "agent_message") {
		t.Fatalf("userMessageOnly must omit agent_message: %s", got)
	}
	if strings.Contains(got, "updated_input") {
		t.Fatalf("userMessageOnly must omit updated_input: %s", got)
	}
}
