package agnostic

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

// TestNewToolCall_BashShell checks that model.NewToolCall normalizes Bash to the
// canonical bash name, preserves Native, and extracts the shell command.
func TestNewToolCall_BashShell(t *testing.T) {
	tc := model.NewToolCall("Bash", []byte(`{"command":"echo hi","description":"test"}`), "id1")
	if tc.Name != tools.ToolBash || tc.Native != "Bash" || tc.Shell != "echo hi" || tc.ID != "id1" {
		t.Fatalf("ToolCall=%+v", tc)
	}
}
