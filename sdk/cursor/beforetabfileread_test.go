package cursor

import (
	"strings"
	"testing"
)

func TestEncode_TabFileReadDeny(t *testing.T) {
	out, code, err := codec.Encode(permissionResults{}.Deny("no tab reads"))
	if err != nil {
		t.Fatal(err)
	}
	if code != PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, PermissionDenyExit)
	}
	if !strings.Contains(string(out), `"permission":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_BeforeTabFileRead(t *testing.T) {
	mustDecode[BeforeTabFileRead](t, `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`)
}
