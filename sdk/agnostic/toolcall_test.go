package agnostic

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

func TestNewToolCall_BashShell(t *testing.T) {
	tc := newToolCall("Bash", []byte(`{"command":"echo hi","description":"test"}`), "id1")
	if tc.Name != tools.ToolBash || tc.Native != "Bash" || tc.Shell != "echo hi" || tc.ID != "id1" {
		t.Fatalf("ToolCall=%+v", tc)
	}
}
