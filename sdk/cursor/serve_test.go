package cursor_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

const cursorShell = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "beforeShellExecution",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "user_email": null,
  "transcript_path": null,
  "command": "git push --force",
  "cwd": "/w",
  "sandbox": false
}`

func TestMux_Serve_BeforeShellDeny(t *testing.T) {
	run.Reset()
	cursor.OnBeforeShellExecution(func(ctx context.Context, hook run.Hook[cursor.BeforeShellExecution], r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		return r.Deny("blocked"), nil
	})
	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(cursorShell), &stdout, &bytes.Buffer{})
	if code != cursor.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursor.PermissionDenyExit)
	}
	if !strings.Contains(stdout.String(), `"permission":"deny"`) {
		t.Fatalf("bad stdout: %s", stdout.String())
	}
}
