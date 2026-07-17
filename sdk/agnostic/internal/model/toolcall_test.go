package model

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

// TestNewToolCall_BashShell checks that NewToolCall normalizes Bash to the
// canonical bash name, preserves Native, and extracts the shell command.
func TestNewToolCall_BashShell(t *testing.T) {
	tc := NewToolCall("Bash", []byte(`{"command":"echo hi","description":"test"}`), "id1")
	if tc.Name != tools.ToolBash || tc.Native != "Bash" || tc.Shell != "echo hi" || tc.ID != "id1" {
		t.Fatalf("ToolCall=%+v", tc)
	}
}
