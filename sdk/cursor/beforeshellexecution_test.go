package cursor

import (
	"encoding/json"
	"testing"
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

func TestDecodeEncode_BeforeShellDeny(t *testing.T) {
	ev, err := codec.Decode([]byte(cursorShell))
	if err != nil {
		t.Fatal(err)
	}
	shell, ok := ev.(BeforeShellExecution)
	if !ok {
		t.Fatalf("want BeforeShellExecution, got %T", ev)
	}
	if shell.Command != "git push --force" || shell.ConversationID != "c1" {
		t.Fatalf("bad event: %+v", shell)
	}

	out, code, err := codec.Encode(permissionResults{}.Deny("force push blocked"))
	if err != nil {
		t.Fatal(err)
	}
	if code != PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, PermissionDenyExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "deny" || got["agent_message"] != "force push blocked" {
		t.Fatalf("bad output: %s", out)
	}
}
