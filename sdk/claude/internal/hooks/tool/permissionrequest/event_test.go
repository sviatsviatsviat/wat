package permissionrequest

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, _, err := results{}.Deny("policy").WithInterrupt(true).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), `"continue":false`) {
		t.Fatalf("interrupt must not set top-level continue: %s", out)
	}
	if !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestUpdatedPermissions(t *testing.T) {
	updates := []event.PermissionUpdate{{
		Type:        event.PermissionUpdateAddRules,
		Behavior:    event.DecisionAllow,
		Destination: event.PermissionDestinationLocalSettings,
		Rules: []event.PermissionRule{{
			ToolName:    "Bash",
			RuleContent: "rm -rf node_modules",
		}},
	}}
	out, code, err := results{}.Allow().WithUpdatedPermissions(updates).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	var top map[string]any
	if err := json.Unmarshal(out, &top); err != nil {
		t.Fatal(err)
	}
	hso, _ := top["hookSpecificOutput"].(map[string]any)
	dec, _ := hso["decision"].(map[string]any)
	got, ok := dec["updatedPermissions"].([]any)
	if !ok || len(got) != 1 {
		t.Fatalf("updatedPermissions = %v", dec["updatedPermissions"])
	}
	entry, _ := got[0].(map[string]any)
	if entry["type"] != "addRules" || entry["behavior"] != "allow" || entry["destination"] != "localSettings" {
		t.Fatalf("entry = %v", entry)
	}
	// Deny must not emit updatedPermissions even if set.
	denied, _, err := results{}.Deny("no").WithUpdatedPermissions(updates).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(denied), "updatedPermissions") {
		t.Fatalf("deny must not emit updatedPermissions: %s", denied)
	}
}

func TestDecode_PermissionRequest(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"PermissionRequest","tool_name":"Write","tool_use_id":"t2"}`, event.PermissionRequest)
}

const permissionRequestPayload = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/Users/...",
  "permission_mode": "default",
  "hook_event_name": "PermissionRequest",
  "tool_name": "Bash",
  "tool_input": {
    "command": "rm -rf node_modules",
    "description": "Remove node_modules directory"
  },
  "permission_suggestions": [
    {
      "type": "addRules",
      "rules": [{ "toolName": "Bash", "ruleContent": "rm -rf node_modules" }],
      "behavior": "allow",
      "destination": "localSettings"
    }
  ]
}`

func TestDecode_PermissionRequestSuggestions(t *testing.T) {
	ev := mustDecode[Event](t, permissionRequestPayload, event.PermissionRequest)
	if ev.ToolName != "Bash" || ev.SessionID != "abc123" {
		t.Fatalf("bad event: %+v", ev)
	}
	if len(ev.PermissionSuggestions) != 1 {
		t.Fatalf("suggestions = %+v", ev.PermissionSuggestions)
	}
	s := ev.PermissionSuggestions[0]
	if s.Type != event.PermissionUpdateAddRules ||
		s.Behavior != event.DecisionAllow ||
		s.Destination != event.PermissionDestinationLocalSettings {
		t.Fatalf("suggestion = %+v", s)
	}
	if len(s.Rules) != 1 || s.Rules[0].ToolName != "Bash" || s.Rules[0].RuleContent != "rm -rf node_modules" {
		t.Fatalf("rules = %+v", s.Rules)
	}
	bash, ok := ev.ToolInput.AsBash()
	if !ok || bash.Command != "rm -rf node_modules" {
		t.Fatalf("tool input = %+v ok=%v", bash, ok)
	}
}

func TestMerge_PermissionRequest_updatedPermissionsOverwriteWarns(t *testing.T) {
	first := []event.PermissionUpdate{{
		Type:        event.PermissionUpdateSetMode,
		Mode:        event.PermissionAcceptEdits,
		Destination: event.PermissionDestinationSession,
	}}
	second := []event.PermissionUpdate{{
		Type:        event.PermissionUpdateAddDirectories,
		Directories: []string{"/tmp/extra"},
		Destination: event.PermissionDestinationSession,
	}}
	origFirstRules := append([]event.PermissionRule(nil), first[0].Rules...)
	origSecondDirs := append([]string(nil), second[0].Directories...)
	a := results{}.Allow().WithUpdatedPermissions(first)
	b := results{}.Allow().WithUpdatedPermissions(second)
	// Mutate caller-owned nested slices after With*; merge must keep clones.
	second[0].Directories[0] = "/mutated"
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] != "updatedPermissions: overwritten by later handler" {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if len(out.updatedPermissions) != 1 ||
		out.updatedPermissions[0].Type != event.PermissionUpdateAddDirectories ||
		!reflect.DeepEqual(out.updatedPermissions[0].Directories, []string{"/tmp/extra"}) {
		t.Fatalf("updatedPermissions = %+v", out.updatedPermissions)
	}
	if !reflect.DeepEqual(first[0].Rules, origFirstRules) {
		t.Fatalf("caller first rules mutated: %+v", first[0].Rules)
	}
	if !reflect.DeepEqual(origSecondDirs, []string{"/tmp/extra"}) {
		t.Fatalf("origSecondDirs = %v", origSecondDirs)
	}
}

func TestEncode_PermissionRequestEchoSuggestion(t *testing.T) {
	ev := mustDecode[Event](t, permissionRequestPayload, event.PermissionRequest)
	out, _, err := results{}.Allow().WithUpdatedPermissions(ev.PermissionSuggestions).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"updatedPermissions"`) ||
		!strings.Contains(string(out), `"localSettings"`) {
		t.Fatalf("bad echo output: %s", out)
	}
}

func init() {
	register(testCodec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != wantName {
		t.Fatalf("EventName() = %q, want %q", ev.EventName(), wantName)
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}
