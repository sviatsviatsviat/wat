package cursor

import (
	"testing"
)

func TestDecode_WorkspaceOpen(t *testing.T) {
	mustDecode[WorkspaceOpen](t, `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`)
}
