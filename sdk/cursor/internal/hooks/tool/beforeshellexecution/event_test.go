package beforeshellexecution

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
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
	ev, err := runtime.Codec.Decode([]byte(cursorShell))
	if err != nil {
		t.Fatal(err)
	}
	shell, ok := ev.(Event)
	if !ok {
		t.Fatalf("want BeforeShellExecution, got %T", ev)
	}
	if shell.Command != "git push --force" || shell.ConversationID != "c1" {
		t.Fatalf("bad event: %+v", shell)
	}

	out, code, err := event.NewPermissionResults().Deny("force push blocked").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, event.PermissionDenyExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "deny" || got["agent_message"] != "force push blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func init() {
	Register(runtime.Codec)
}
