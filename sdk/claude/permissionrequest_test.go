package claude

import (
	"strings"
	"testing"
)

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, _, err := permissionRequestResults{}.Deny("policy").WithInterrupt(true).Encode()
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

func TestDecode_PermissionRequest(t *testing.T) {
	mustDecode[PermissionRequest](t, `{"session_id":"s","hook_event_name":"PermissionRequest","tool_name":"Write","tool_use_id":"t2"}`, EventPermissionRequest)
}
