package copilot

import (
	"strings"
	"testing"
)

func TestEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	out, code, err := permissionRequestResults{}.Deny("blocked").WithInterrupt(true).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != WarnExit {
		t.Fatalf("code=%d, want %d", code, WarnExit)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) || !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestAsk(t *testing.T) {
	out, code, err := permissionRequestResults{}.Ask("needs user confirmation").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_PermissionRequest(t *testing.T) {
	e := mustDecode[PermissionRequest](t, `{"hook_event_name":"PermissionRequest","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"create","tool_input":{"path":"a.txt"}}`, EventPermissionRequest)
	if e.NativeToolName() != "create" {
		t.Fatalf("PermissionRequest=%+v", e)
	}
}
