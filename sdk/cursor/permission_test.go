package cursor

import (
	"strings"
	"testing"
)

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := codec.Encode(permissionResults{}.Noop())
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_PermissionUpdatedInput(t *testing.T) {
	out, code, err := codec.Encode(permissionResults{}.Allow().WithUpdatedInput(map[string]any{"command": "ls"}))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"updated_input"`) {
		t.Fatalf("bad output: %s", out)
	}
}
