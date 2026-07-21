package cursor

import (
	"strings"
	"testing"
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
